FROM golang:1.20.2-bullseye

RUN mkdir /app
ADD . /app/
WORKDIR /app

ENV PATH="${PATH}:/app"

RUN go mod download && go mod verify
RUN go build -o ssid-jungle main.go

CMD ["./ssid-jungle"]