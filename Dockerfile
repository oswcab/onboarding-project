FROM registry.access.redhat.com/ubi9/go-toolset:1.21.11-7.1724661022 as builder

# Copy the Go Modules manifests
COPY go.mod go.mod
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY main.go main.go

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o main main.go

# Use ubi-minimal as minimal base image to package the manager binary
# Refer to https://catalog.redhat.com/software/containers/ubi8-minimal/5c64772edd19c77a158ea216 for more details
FROM registry.access.redhat.com/ubi9/ubi-minimal:9.4-1194
COPY --from=builder /opt/app-root/src/main /
USER 65532:65532

ENV PORT 8081
EXPOSE 8081

CMD [ "./main" ]
