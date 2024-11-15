FROM golang:1.22 as builder

WORKDIR /opt/ac/serverless-cron

COPY main.go main.go
COPY go.mod go.mod
COPY go.sum go.sum

RUN go install .

# ==============================================================================

FROM gcr.io/distroless/base

COPY --from=builder /go/bin/serverless-cron .

ENV VERSION=production
CMD [ "./serverless-cron" ]
