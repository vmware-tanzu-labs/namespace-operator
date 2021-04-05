# Build the manager binary
FROM golang:1.14 as builder

# take the username/password to the private repo as an argument
# NOTE: only do this in the builder stage as not to leave this behind in the final image
ARG GIT_HTTPS_USERNAME
ARG GIT_HTTPS_PASSWORD
RUN echo "machine gitlab.eng.vmware.com login ${GIT_HTTPS_USERNAME} password ${GIT_HTTPS_PASSWORD}" > /root/.netrc

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY main.go main.go
COPY api/ api/
COPY controllers/ controllers/

# Build
# TODO: the following command must be run manually prior to building
# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o manager main.go
COPY manager manager

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/manager .
USER nonroot:nonroot

ENTRYPOINT ["/manager"]
