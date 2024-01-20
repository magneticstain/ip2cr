# syntax=docker/dockerfile:1

FROM golang:alpine3.19

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o ./ip2cr

ENTRYPOINT [ "./ip2cr" ]
CMD [ "--help" ]