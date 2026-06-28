FROM golang:1.26-bookworm AS build 

WORKDIR /src

RUN apt-get update \
	&& apt-get install -y --no-install-recommends gcc libc6-dev \
	&& rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=1
RUN go build -o /forum .

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=build /forum /app/forum
COPY ui ./ui

EXPOSE 8080

CMD ["/app/forum"]
