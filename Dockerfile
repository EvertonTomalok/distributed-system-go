FROM golang:1.19 AS build
WORKDIR /go/src/app
COPY . .
RUN go install -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -o exe main.go

FROM alpine
RUN apk add --no-cache ca-certificates
COPY --from=build /go/src/app .
