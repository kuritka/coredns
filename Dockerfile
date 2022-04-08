# Dockerfile
# Build the manager binary
FROM golang:1.18 as builder
ARG EXPOSE_PORT=5053
WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY coredns.go coredns.go
COPY directives_generate.go directives_generate.go
COPY owners_generate.go owners_generate.go
COPY core/ core/
COPY coremain/ coremain/
COPY pb/ pb/
COPY request/ request/
COPY plugin/ plugin/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./coredns

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
# used for tetsing purposes
# FROM alpine:latest
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/coredns .
USER 65532:65532
EXPOSE $EXPOSE_PORT $EXPOSE_PORT/udp
ENTRYPOINT ["/coredns"]
