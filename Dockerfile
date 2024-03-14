FROM golang:latest as build

WORKDIR /app
COPY . .

ARG SERVER_PORT=3001
ENV SERVER_PORT=${SERVER_PORT}
RUN go mod download
# ignore error from test
RUN go test -p 1 -v -cover -short ./... || true

RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm  go build -o /go/bin/app cmd/server/main.go

FROM gcr.io/distroless/static-debian12

WORKDIR /usr/src/app
COPY --from=build /go/bin/app ./server
CMD ["/usr/src/app/server"]
EXPOSE ${SERVER_PORT}
