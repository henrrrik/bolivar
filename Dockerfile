FROM golang:1.26.1 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /bolivar .

FROM gcr.io/distroless/static-debian12
WORKDIR /app
COPY --from=build /bolivar /app/bolivar
COPY static/ /app/static/
ENTRYPOINT ["/app/bolivar"]
