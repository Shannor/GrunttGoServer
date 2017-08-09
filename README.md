# Gruntt-Comics Backend

A webcrawler service written in Go to get comic information from hosting websites.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### Prerequisites

1. Download [Google Cloud SDK](https://cloud.google.com/sdk/docs/)
2. Install go gcloud component

    ```
    gcloud components install app-engine-go
    ```

3. Need to have [Python 2.7](https://www.python.org/download/releases/2.7/) installed

    Google Cloud SDK may come with Python bundled but safe to installed manually.

4. Install [Glide](https://glide.sh/)

### Installing

After cloning down the repo, use glide to get the dependencies.

```
glide install
```

### Local Deployment
To deploy the server locally run the command
```
dev_appserver.py main/app.yaml
```
Then the server can be reached by using "localhost:8080"

## Running the tests

Tests are done using Ginkgo and Gomega.
To run the tests run
```
go test
```

or

```
ginkgo
```
### Break down into end to end tests

Explain what these tests test and why

```
Give an example
```

### And coding style tests

Explain what these tests test and why

```
Give an example
```

## Deployment

Add additional notes about how to deploy this on a live system

## Contributing

Please read [CONTRIBUTING.md](https://gist.github.com/PurpleBooth/b24679402957c63ec426) for details on our code of conduct, and the process for submitting pull requests to us.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/your/project/tags).


## Acknowledgments

* Hat tip to anyone who's code was used
* Inspiration
* etc
