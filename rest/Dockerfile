FROM golang:1.25-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/echelon-rest \
    ./rest/cmd/echelon-rest

FROM alpine:3.23.5

RUN addgroup -S echelon && adduser -S -G echelon echelon

WORKDIR /app

COPY --from=builder /out/echelon-rest /app/echelon-rest
COPY cli/cmd/echelon-utility/config.yaml /app/config.yaml

USER echelon

EXPOSE 8080

ENTRYPOINT ["/app/echelon-rest"]
CMD ["-address", ":8080"]
