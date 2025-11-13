FROM golang:1.25.3 AS builder

WORKDIR /build

# build app
COPY go.* ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal

RUN CGO_ENABLED=0 go build -o /build/app ./cmd/pr-review/main.go

FROM alpine:3.22.2

RUN apk update && apk add --no-cache curl

# run as non-root user
RUN adduser --disabled-password -u 1000 user
USER user

WORKDIR /app
COPY --chown=user:user --from=builder /build/app app

RUN mkdir ./docs
COPY --chown=user:user docs/swagger.* ./docs

CMD [ "/app/app" ]