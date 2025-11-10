package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/CristianCurteanu/reco/pkg/asana"
	"github.com/CristianCurteanu/reco/pkg/storage"
	"github.com/CristianCurteanu/reco/pkg/ticker"
)

var (
	asanaAccessToken = flag.String("asana-access-token", "", "This is the Asana PAT (required)\nCheck this page how to set it up https://developers.asana.com/docs/personal-access-token")
	asanaAPIHost     = flag.String("asana-host", "https://app.asana.com/api/1.0", "This parameter is used in case the Asana API URL will be different that the one provided from official docs")
	extractionPeriod = flag.String("extraction-period", "30s", "Period of time between extraction jobs; it's either 30s or 5m")
)

func main() {

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	outputDir := flag.String("output-dir", filepath.Join(wd, "output"), "")

	flag.Parse()
	// Step 1: Initialize the Asana API Client

	if asanaAccessToken == nil {
		log.Fatalf("unable to initialize without Asana API Token; check this page how to set it up https://developers.asana.com/docs/personal-access-token")
	}

	apiClient := asana.NewAPIClient(*asanaAPIHost, *asanaAccessToken)

	log.Printf("Asana API Extractor running (pid: %d)", os.Getpid())

	// Step 3: Setup the Asana Extractor, and inject it to Periodic Extractor
	asanaExtractor := asana.NewExtractor(apiClient)

	// Step 4: Run the Periodic Extractor
	period, found := ticker.GetExtractionPeriod(*extractionPeriod)
	if !found {
		log.Fatal("please specify proper period config value, either `30s` or `5m`")
	}
	scheduler := ticker.NewScheduler()
	defer scheduler.Stop()

	fileStorage := storage.NewFile(*outputDir)
	scheduler.Run("get all users", period, func() error {
		users, err := asanaExtractor.GetAllUsers()
		if err != nil {
			return err
		}

		usersData, err := json.MarshalIndent(users, "", "  ")
		if err != nil {
			log.Printf("failed to marshal users to JSON, err=%q", err)

			return err
		}

		tn := time.Now().UTC()

		return fileStorage.Store(fmt.Sprintf("%d_users.json", tn.Unix()), usersData)
	})

	scheduler.Run("get all projects", period, func() error {
		projects, err := asanaExtractor.GetAllProjects()
		if err != nil {
			return err
		}

		projectsData, err := json.MarshalIndent(projects, "", "  ")
		if err != nil {
			log.Printf("failed to marshal users to JSON, err=%q", err)

			return err
		}

		tn := time.Now().UTC()

		return fileStorage.Store(fmt.Sprintf("%d_projects.json", tn.Unix()), projectsData)
	})

	scheduler.Wait()
}
