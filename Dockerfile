FROM golang:stretch AS build
WORKDIR /build
RUN apt-get update && \
    apt-get install -y xz-utils
ADD https://github.com/upx/upx/releases/download/v3.95/upx-3.95-amd64_linux.tar.xz /usr/local
RUN xz -d -c /usr/local/upx-3.95-amd64_linux.tar.xz | \
    tar -xOf - upx-3.95-amd64_linux/upx > /bin/upx && \
    chmod a+x /bin/upx
COPY . .
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux \
    go build -a -tags netgo -ldflags '-w -s'  main.go && \
    upx main


# Last stage
FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /build/.env .env
COPY --from=build /build/public /public
COPY --from=build /build/main /main
WORKDIR /
EXPOSE 8000
CMD ["./main"]