FROM alpine:3.2
MAINTAINER Eagle Liut <eagle@dantin.me>

ENV VERSION v0.0.2
ENV DOWNLOAD_URL https://github.com/liut/staffio/releases/download/$VERSION/staffio-linux-amd64-$VERSION.tar.gz

ENV STAFFIO_HTTP_LISTEN=":80"
ENV STAFFIO_ROOT /app

RUN mkdir /app
WORKDIR /app

RUN apk add --virtual build-dependencies --update \
  curl \
  ca-certificates \
  && curl -L $DOWNLOAD_URL | tar xvz -C /usr/local/bin \
  && apk del build-dependencies \
  && rm -rf /var/cache/apk/*

ADD templates /app/templates
# ADD htdocs /app/htdocs

EXPOSE 80

# ENTRYPOINT ["/usr/local/bin/staffio"]