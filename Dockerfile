# ---------- BUILD STAGE ----------
FROM golang:1.26.1-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

# cache deps
COPY go.mod go.sum ./
RUN go mod download

# copy source
COPY . .

# build static binary (smaller + portable)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./server

# ---------- RUN STAGE ----------
FROM alpine:latest

WORKDIR /app

# add ca certs (important for HTTP clients later)
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/app .

EXPOSE 9090

CMD ["./app"]
