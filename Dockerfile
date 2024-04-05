###############################################################################
# BEGIN build-stage
# Compile the binary
FROM --platform=$BUILDPLATFORM docker.io/library/golang:1.22.2@sha256:c4fb952e712efd8f787bcd8e53fd66d1d83b7dc26adabc218e9eac1dbf776bdf AS build-stage

ARG BUILDPLATFORM
ARG TARGETARCH

WORKDIR /app

COPY vendor ./vendor
COPY go.mod go.sum ./
COPY cmd ./cmd
COPY pkg ./pkg

RUN CGO_ENABLED=0 GOOS=linux GOARCH="${TARGETARCH}" go build -ldflags="-w -s" -o /simple-fileserver ./cmd/

#
# END build-stage
###############################################################################

###############################################################################
# BEGIN test-stage
# Run the tests in the container
FROM docker.io/library/golang:1.22.2@sha256:c4fb952e712efd8f787bcd8e53fd66d1d83b7dc26adabc218e9eac1dbf776bdf AS test-stage

WORKDIR /app

COPY --from=build-stage /app /app
# Not needed for testing, but needed for later stage
COPY --from=build-stage /simple-fileserver /

RUN go test -v ./...

#
# END test-stage
###############################################################################

###############################################################################
# BEGIN final-stage
# Create final docker image
FROM scratch AS final-stage

WORKDIR /

COPY --from=test-stage /simple-fileserver /

EXPOSE 8080

USER 1001

ENTRYPOINT ["/simple-fileserver", "-webroot", "/webroot"]

#
# END final-stage
###############################################################################
