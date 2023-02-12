FROM golang:1.19.3-alpine AS builder
WORKDIR /app

RUN --mount=type=secret,id=git-credentials cp /run/secrets/git-credentials ~/.git-credentials

# todo: remove this after repositories are public (+ git-credentials)
RUN apk add git
RUN apk add gcc g++
RUN git config --global credential.helper store

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY cmd/ ./cmd/
COPY internal/ ./internal/

RUN GOOS=linux CGO_ENABLED=1 go build -ldflags="-extldflags=-static" -o subscribe ./cmd/subscribe/subscribe.go

FROM scratch
WORKDIR /app

COPY --from=builder /app /app/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

CMD [ "/app/subscribe" ]
