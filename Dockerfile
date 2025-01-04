# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM golang:alpine AS build-stage

# Create and change to the app directory.
WORKDIR /app

# Retrieve application dependencies using go modules.
# Allows container builds to reuse downloaded dependencies.
COPY go.* ./
RUN go mod download

# Copy local code to the container image.
COPY . ./

# Install chromium so go-rod can use it directly without downloading inside the container
RUN apk add chromium

# Build the binary.
# -mod=readonly ensures immutable go.mod and go.sum in container builds.
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -v -o server

FROM build-stage AS run-test-stage
RUN go get github.com/Thatooine/go-test-html-report
RUN go get github.com/jstemmer/go-junit-report
RUN go install github.com/Thatooine/go-test-html-report
RUN go install github.com/jstemmer/go-junit-report
RUN CGO_ENABLED=0 GOOS=linux go test -v -cover -coverprofile coverage.out -json ./... | go-junit-report

# Build the runtime container image from scratch, copying what is needed from the previous stage.  
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
# FROM scratch

# Copy the binary to the production image from the builder stage.
# COPY --from=builder /app/server /server

# Run the web service on container startup.
EXPOSE 8080/tcp
ENTRYPOINT ["/app/server"]