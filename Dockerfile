FROM golang:1.22.1-alpine3.19

WORKDIR /app

COPY ./ ./

RUN go build -o /bin/app cmd/main.go

ENTRYPOINT ["app"]