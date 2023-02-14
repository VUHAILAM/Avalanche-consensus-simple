FROM golang:1.18-alpine as build
WORKDIR /
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o snow ./cmd/snow/main.go

# This is for documentation purposes only.
# To actually open the port, runtime parameters
# must be supplied to the docker command.
EXPOSE 8080

FROM alpine:latest
WORKDIR /
COPY --from=build /snow ./snow
ENTRYPOINT ["./snow"]


