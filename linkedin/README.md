# LinkedIn Profile Grabber

This project is meant to be a demonstration of some of my abilities in the software development space. It consists of a fully functional _Golang_ application that is able to retrieve my own personal
LinkedIn profile information in real-time and provide the data as a _JSON_ object.

## Build

Run `go build` to build the project.

## Running unit tests

Run `go test` to test the project. This project has pretty complete unit tests using the _ginkgo_ testing framework.

## Fully automated build and deploy

This project is integrated with _Google Cloud Build_ to have commits on the `main` branch automatically trigger _CI/CD_ automation that builds and deploys the Golang application.
The application is deployed into _Google Cloud Run_ as a REST API that is then used by my actual [web site project](https://github.com/invinity/mattpitts-site) to retrieve my own LinkedIn profile information to generate a live resume.
