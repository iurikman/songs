FROM golang:1.22 as builder

WORKDIR /usr/src/app
RUN go mod download
COPY . .
RUN go build -v -o songs ./...

FROM debian:stable-slim
WORKDIR /bin
COPY --from=builder /usr/src/app/songs ./
ENTRYPOINT ["./songs"]