FROM golang:1.26.5-alpine AS build
WORKDIR /src
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -o /out/crownfall ./cmd/server

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/crownfall /crownfall
COPY migrations /migrations
ENV CROWNFALL_MIGRATIONS_PATH=/migrations
EXPOSE 8080
ENTRYPOINT ["/crownfall"]
