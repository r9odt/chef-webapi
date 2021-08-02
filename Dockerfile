FROM golang:1.16.6-alpine as binary_build

ENV GO_SRC=/usr/local/go/src
WORKDIR ${GO_SRC}/github.com/JIexa24/chef-webapi
COPY . .
RUN apk add --no-cache git make
RUN go get -d ./... && go mod vendor

RUN make build && ${GO_SRC}/github.com/JIexa24/chef-webapi/bin/web -version

#-------------------------------------------------------------------------------
FROM node:15.14.0-alpine as web_build

RUN apk add --no-cache git python make g++
WORKDIR /app
COPY content/webjs .
ENV NODE_ENV=production
RUN yarn install
RUN PUBLIC_URL="" yarn build

#-------------------------------------------------------------------------------
FROM alpine:3.13.5

ARG IMAGE_USER=webapp
ENV GO_SRC=/usr/local/go/src \
    IMAGE_USER=${IMAGE_USER} \ 
    GOSU_VERSION=1.13

RUN addgroup -g 4253 -S ${IMAGE_USER} \
    && adduser -u 4253 -G ${IMAGE_USER} -s /usr/sbin/nologin -H -D ${IMAGE_USER}
RUN set -eux; \
    apk add --no-cache --virtual .gosu-deps \
      ca-certificates \
      dpkg \
      gnupg \
      wget \
    ; \
    dpkgArch="$(dpkg --print-architecture | awk -F- '{ print $NF }')"; \
    wget --quiet -O /usr/local/bin/gosu "https://github.com/tianon/gosu/releases/download/$GOSU_VERSION/gosu-$dpkgArch"; \
    wget --quiet -O /usr/local/bin/gosu.asc "https://github.com/tianon/gosu/releases/download/$GOSU_VERSION/gosu-$dpkgArch.asc"; \
    export GNUPGHOME="$(mktemp -d)"; \
    gpg --batch --keyserver hkps://keys.openpgp.org --recv-keys B42F6819007F00F88E364FD4036A9C25BF357DD4; \
    gpg --batch --verify /usr/local/bin/gosu.asc /usr/local/bin/gosu; \
    command -v gpgconf && gpgconf --kill all || :; \
    rm -rf "$GNUPGHOME" /usr/local/bin/gosu.asc; \
    apk del --no-network .gosu-deps; \
    chmod +x /usr/local/bin/gosu; \
    gosu --version; \
    gosu nobody true

WORKDIR /app
COPY --from=binary_build --chown=4253:4253 ${GO_SRC}/github.com/JIexa24/chef-webapi/bin /app/
COPY --from=binary_build --chown=4253:4253 ${GO_SRC}/github.com/JIexa24/chef-webapi/keys.key /app/
COPY --from=web_build --chown=4253:4253 /app/build /app/content/webjs/build

EXPOSE 3000 8082

ENTRYPOINT ["/app/docker-entrypoint.sh"]
CMD ["/app/web"]
