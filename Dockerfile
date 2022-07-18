FROM golang:1.17.7-alpine3.15 as build
ARG OS
ARG ARCH
COPY . /build/
WORKDIR /build
RUN go mod download && GOOS=${OS} GOARCH=${ARCH} go build -o systemstat

FROM alpine:3.15
ARG VERSION
ARG user=navilg
ARG group=navilg
ARG uid=1000
ARG gid=1000
USER root
WORKDIR /app
COPY --from=build /build/systemstat /app/systemstat
COPY container-entrypoint.sh /app/container-entrypoint.sh
RUN apk update && apk --no-cache add bash && addgroup -g ${gid} ${group} && adduser -h /app -u ${uid} -G ${group} -s /bin/bash -D ${user}
RUN chown navilg:navilg /app/systemstat && chmod +x /app/systemstat && \
    chown navilg:navilg /app/container-entrypoint.sh && chmod +x /app/container-entrypoint.sh
USER navilg
EXPOSE 30000
ENTRYPOINT [ "/app/container-entrypoint.sh"]