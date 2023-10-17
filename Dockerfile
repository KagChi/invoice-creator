FROM golang:1.21-alpine as build-stage

WORKDIR /tmp/build

COPY . .

# Build the project
RUN go build .

FROM alpine:3

LABEL name "Invoice Creator"
LABEL maintainer "KagChi"

WORKDIR /app

# Install needed deps
RUN apk add --no-cache tini

COPY --from=build-stage /tmp/build/invoice-creator main

ENTRYPOINT ["tini", "--"]
CMD ["/app/main"]