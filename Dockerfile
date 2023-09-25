# Build the application from source
FROM golang:1.21 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /example-app

# Deploy the application binary into a lean image
FROM scratch

WORKDIR /

COPY --from=build-stage /example-app /example-app

EXPOSE 8080

ENTRYPOINT ["/example-app"]
