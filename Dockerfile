# build stage
FROM golang:1.14 as builder

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build .

RUN mkdir publish && cp hx publish && \
    cp -r docs publish && cp -r config publish

# final stage
FROM scratch

WORKDIR /app

COPY --from=builder /app/publish .

ENV GIN_MODE=release \
    PORT=8080

EXPOSE 8080

ENTRYPOINT ["./hx"]