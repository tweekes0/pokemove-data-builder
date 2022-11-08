# syntax=docker/dockerfile:1
FROM golang:1.18 as builder
WORKDIR /usr/local/go/src/pokemoves-backend
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . ./
RUN CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-w -s" -o /usr/local/bin/server ./src/backend/cmd/

FROM alpine:latest
RUN apk --no-cache add ca-certificates \
  && update-ca-certificates
WORKDIR /pokemoves-server/
COPY --from=builder /usr/local/bin/server ./
# COPY --from=builder /usr/local/go/src/pokemoves-backend/src/internal/conf ./conf
ENTRYPOINT [ "./server" ] 