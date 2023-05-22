ARG go_ver=1.19
ARG alpine_ver=3.16

ARG DOCKER_BASE=golang:${go_ver}-alpine${alpine_ver}
FROM ${DOCKER_BASE} AS stage

RUN apk add --no-cache make git

WORKDIR /usr/arrai
COPY . .
RUN make build

FROM ${DOCKER_BASE}
COPY --from=stage /usr/arrai/arrai /bin/arrai

ENTRYPOINT ["/bin/arrai"]
