FROM alpine:latest

RUN apk add --no-cache curl busybox-suid

COPY crontab.txt /etc/crontabs/root

ENTRYPOINT ["crond", "-f", "-L", "/proc/1/fd/1"]
