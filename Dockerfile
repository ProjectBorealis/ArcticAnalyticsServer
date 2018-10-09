FROM golang:1.11-alpine as builder

WORKDIR /go/src/projectborealisgitlab.site/project-borealis/programming/dev-ops/aa-server
COPY . .

ENV GO111MODULE=on
RUN apk add --no-cache git
RUN CGO_ENABLED=0 GOOS=linux go test ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo .

FROM alpine:latest
RUN \
    apk --no-cache add ca-certificates && \
    addgroup app && adduser -S -G app app && \
    mkdir -p /home/app/data && \
    chown -R app:app /home/app/data

WORKDIR /home/app
USER app

COPY --from=builder /go/src/projectborealisgitlab.site/project-borealis/programming/dev-ops/aa-server/aa-server .

CMD ["./aa-server"]