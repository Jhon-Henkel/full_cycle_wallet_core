FROM golang:1.22

WORKDIR /app/

RUN apt-get update && apt-get install -y librdkafka-dev

CMD ["go", "run", "/app/cmd/wallet_core/main.go"]