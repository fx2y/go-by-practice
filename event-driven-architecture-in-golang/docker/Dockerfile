FROM golang:1.19-alpine AS builder

WORKDIR /mallbots

COPY go.mod go.sum ./
RUN go mod download

COPY .. ./

RUN go build -v -o monolith ./cmd/mallbots

FROM alpine:latest AS runtime

COPY --from=builder /mallbots/docker/wait-for .
RUN chmod +x /wait-for

COPY --from=builder /mallbots/monolith /mallbots/monolith

CMD ["/mallbots/monolith"]