FROM golang:latest as builder

WORKDIR /app
COPY . .
RUN go mod download
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o bin/api cmd/api/main.go

FROM scratch

WORKDIR /app
COPY --from=builder /app/bin/api .
COPY --from=builder /app/.env .

ENTRYPOINT [ "./api" ]