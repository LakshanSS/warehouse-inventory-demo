FROM golang:1.22-alpine as build-env

RUN mkdir /app
WORKDIR /app

COPY go.mod ./
RUN go mod download

# Create non-root user
RUN addgroup -g 10014 choreo && \
    adduser --disabled-password --no-create-home --uid 10014 --ingroup choreo choreouser

COPY . .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Generate go.sum and build the binary
RUN go mod tidy && go build -a -installsuffix cgo -o /go/bin/inventory-api -buildvcs=false

FROM alpine:3.19

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

COPY --from=build-env /go/bin/inventory-api /go/bin/inventory-api
COPY --from=build-env /etc/passwd /etc/passwd

USER 10014

EXPOSE 9090

ENTRYPOINT ["/go/bin/inventory-api"]
