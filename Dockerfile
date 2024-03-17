# syntax=docker/dockerfile:1

FROM golang:1.22-bookworm as build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /simplesso ./cmd/simplesso


FROM gcr.io/distroless/base-debian12 AS build-release-stage

WORKDIR /

COPY --from=build-stage /simplesso /simplesso

EXPOSE 5000

CMD ["/simplesso"]


