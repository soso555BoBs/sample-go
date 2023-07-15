FROM golang:1.20-alpine as dev

WORKDIR /app

COPY . .

CMD [ "go", "run", "main.go" ]

FROM golang:1.20-alpine as build

WORKDIR /build

COPY . .

RUN go build -o test-go-app ./main.go

FROM alpine:latest as prd

WORKDIR /app

COPY --from=build /build/test-go-app .

RUN addgroup go && adduser -D -G go go && chown -R go:go /app/test-go-app

CMD [ "./test-go-app" ]
