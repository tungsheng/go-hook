FROM alpine:3.7

LABEL maintainer="Tony Lee <tungsheng@gmail.com>"

RUN apk add --no-cache ca-certificates tzdata && \
  rm -rf /var/cache/apk/*

EXPOSE 3003

ADD bin/gohook /

ENTRYPOINT ["/gohook"]
CMD ["server"]
