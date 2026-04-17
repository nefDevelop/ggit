FROM golang:1.26-alpine3.23 AS binarybuilder
RUN apk --no-cache --no-progress add --virtual \
  build-deps \
  build-base \
  git \
  linux-pam-dev

WORKDIR /app/ggit
COPY . .

RUN ./docker/build/install-task.sh
RUN TAGS="cert pam" task build

FROM alpine:3.23
RUN apk --no-cache --no-progress add \
  bash \
  ca-certificates \
  curl \
  git \
  linux-pam \
  openssh \
  s6 \
  shadow \
  socat \
  tzdata \
  rsync \
  "zlib>1.3.2"

ENV GGIT_CUSTOM=/data/ggit

# Configure LibC Name Service
COPY docker/nsswitch.conf /etc/nsswitch.conf

WORKDIR /app/ggit
COPY docker ./docker
COPY --from=binarybuilder /app/ggit/.bin/ggit .

RUN ./docker/build/finalize.sh

# Configure Docker Container
VOLUME ["/data", "/backup"]
EXPOSE 22 3000
HEALTHCHECK CMD (curl --noproxy localhost -o /dev/null -sS http://localhost:3000/healthcheck) || exit 1
ENTRYPOINT ["/app/ggit/docker/start.sh"]
CMD ["/usr/bin/s6-svscan", "/app/ggit/docker/s6/"]
