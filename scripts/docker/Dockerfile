# Build the application from source
FROM golang:1.20 AS build-stage

WORKDIR /service

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o class-booking-service ./app/services/booking


# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /service/class-booking-service /class-booking-service
COPY --from=build-stage /service/scripts /scripts

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/class-booking-service"]
