FROM golang:1.14 as builder

ENV BUILDPATH $GOPATH/src/github.com/ricmalta/urlshortner/
ENV GO111MODULE=on

WORKDIR $BUILDPATH

COPY . .

RUN CGO_ENABLED=0 go build -mod=vendor -a -tags netgo -ldflags '-w' -o /app .

FROM scratch as release

COPY --from=builder /app /app
COPY --from=builder /go/src/github.com/ricmalta/urlshortner/internal/config/config.yaml /config.yaml
COPY --from=builder /etc/passwd /etc/passwd

EXPOSE 3000

USER nobody

ENTRYPOINT ["/app", "-config", "./config.yaml"]
