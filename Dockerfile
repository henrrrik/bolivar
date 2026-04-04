FROM golang:1.26.1 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /bolivar .

FROM alpine:3
WORKDIR /app
COPY --from=build /bolivar /app/bolivar
EXPOSE 8080
USER 1000
ENTRYPOINT ["/app/bolivar"]
