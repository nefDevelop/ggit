FROM golang:1.26-alpine3.23 AS binarybuilder
RUN apk --no-cache add build-base git linux-pam-dev

WORKDIR /app/ggit
COPY . .

RUN ./docker/build/install-task.sh
RUN TAGS="cert pam bindata" task build

FROM alpine:3.23
RUN apk --no-cache add \
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
WORKDIR /app/ggit

COPY docker ./docker
COPY --from=binarybuilder /app/ggit/.bin/ggit .
COPY frontend/build ./frontend/build

RUN ./docker/build/finalize.sh

VOLUME ["/data", "/backup"]
EXPOSE 22 3005
HEALTHCHECK CMD (curl --noproxy localhost -o /dev/null -sS http://localhost:3005/healthcheck) || exit 1
ENTRYPOINT ["/app/ggit/docker/start.sh"]
CMD ["/usr/bin/s6-svscan", "/app/ggit/docker/s6/"]
