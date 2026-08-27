FROM golang:1.25-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/echelon-grpc \
    ./grpc/cmd/echelon-grpc

FROM alpine:3.23.5

RUN addgroup -S echelon && adduser -S -G echelon echelon

WORKDIR /app

COPY --from=builder /out/echelon-grpc /app/echelon-grpc
COPY cli/cmd/echelon-utility/config.yaml /app/config.yaml

USER echelon

EXPOSE 9090

ENTRYPOINT ["/app/echelon-grpc"]
CMD ["-address", ":9090"]
