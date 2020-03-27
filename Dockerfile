FROM golang:alpine as builder

RUN apk add --no-cache curl git

RUN curl -fsSL -o /usr/local/bin/dep https://github.com/golang/dep/releases/download/v0.5.4/dep-linux-amd64 && chmod +x /usr/local/bin/dep

RUN mkdir -p /go/src/github.com/influenzanet/authentication-service
ADD . /go/src/github.com/influenzanet/authentication-service/
WORKDIR /go/src/github.com/influenzanet/authentication-service

COPY Gopkg.toml Gopkg.lock ./
# copies the Gopkg.toml and Gopkg.lock to WORKDIR

RUN dep ensure -vendor-only
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main .
FROM scratch
COPY --from=builder /go/src/github.com/influenzanet/authentication-service/main /app/
VOLUME ["/app/keys"]
WORKDIR /app
EXPOSE 3100:3100
CMD ["./main"]
