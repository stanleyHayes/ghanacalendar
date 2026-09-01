FROM golang:1.23-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /ghanacalendar ./cmd/api
FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=build /ghanacalendar /app/ghanacalendar
COPY data /app/data
EXPOSE 8080
ENTRYPOINT ["/app/ghanacalendar"]
