FROM golang:alpine AS dependencies
WORKDIR /app
COPY go.mod go.sum .
RUN go mod graph | awk '{if ($1 !~ "@") print $2}' | xargs go get

FROM dependencies AS builder
COPY . .
RUN go build -o web

FROM alpine AS final
COPY --from=builder /app/web .
ENTRYPOINT ["/web"]
