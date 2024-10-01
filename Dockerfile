###############################################################################
# BEGIN build-stage
# Compile the binary
FROM --platform=$BUILDPLATFORM docker.io/library/golang:1.23.1 AS build-stage

ARG BUILDPLATFORM
ARG TARGETARCH
ARG RELEASE_VERSION

WORKDIR /app

COPY go.mod go.sum ./
COPY cmd ./cmd
COPY pkg ./pkg
COPY hack ./hack

RUN --mount=type=bind,target=/app/.git,source=.git GOOS=linux GOARCH="${TARGETARCH}" hack/build.sh

#
# END build-stage
###############################################################################

###############################################################################
# BEGIN final-stage
# Create final docker image
FROM scratch AS final-stage

WORKDIR /

COPY --from=build-stage /app/bin/simple-fileserver /

EXPOSE 8080

USER 1001

ENTRYPOINT ["/simple-fileserver", "-webroot", "/webroot"]

#
# END final-stage
###############################################################################
