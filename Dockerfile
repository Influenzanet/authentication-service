FROM golang:alpine as builder
RUN mkdir /build 
ADD . /build/
WORKDIR /build
RUN apk add --no-cache git \
  && go get github.com/gin-gonic/gin \
  && go get github.com/dgrijalva/jwt-go \
  && apk del git
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main .
FROM scratch
COPY --from=builder /build/main /app/
WORKDIR /app
EXPOSE 3100:3100
CMD ["./main"]
