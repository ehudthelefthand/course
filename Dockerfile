FROM golang:1.19

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o api

EXPOSE 8080

ENTRYPOINT [ "./api" ]
