## Asana Extractor

This project is a data extractor and monitor, that fetches information about User and Projects from Asana

### Installation

In order to install it, clone this repository:

```
git clone git@github.com:CristianCurteanu/asana-extractor.git
```

After, you can build it locally:

```
go build -o bin/build main.go && chmod +x bin/build
```

And the run it:

```
./bin/build -asana-access-token=<your-asana-access-token>
```

In order to get the Asana Persona Access Token (PAT), make sure to consult [this documentation page](https://developers.asana.com/docs/personal-access-token)

### Usage

This is how the final build could be used:

```
Usage of ./bin/build:
    -asana-access-token string
        This is the Asana PAT (required)
        Check this page how to set it up https://developers.asana.com/docs/personal-access-token
    -asana-host string 
        This parameter is used in case the Asana API URL will be different that the one provided from official docs (default "https://app.asana.com/api/1.0")
    -extraction-period string
        Period of time between extraction jobs; it's either 30s or 5m (default "30s")
    -output-dir string
        (default "/<your-current-workind-directory>/output")

```

Which means that if you would like to use a different extraction period, every 5 minutes for instance, you would use this command:

```
$ ./bin/build -extraction-period=5m \
              -asana-access-token=<your-asana-access-token>
```

### TODOs
- Replace hardcoded values from Asana API Client, Extractor
- Make the status handlers cleaner for Asana API Client
- Finish End-to-End tests, and add corner case tests, for handling different status responses and errors
- Make the periodic scheduler much more configurable from the cli