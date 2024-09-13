FROM golang:1.23 as builder
ARG CGO_ENABLED=0
WORKDIR /app

COPY . .

RUN go build

FROM scratch
COPY --from=builder /app/backend /backend
ENTRYPOINT ["/backend"]