ARG go_ver=1.19
ARG alpine_ver=3.16

FROM golang:${go_ver}-alpine${alpine_ver} AS stage
RUN apk add --no-cache make git

WORKDIR /usr/arrai
COPY . .
RUN make build

FROM golang:${go_ver}-alpine${alpine_ver}
COPY --from=stage /usr/arrai/arrai /bin/arrai

ENTRYPOINT ["/bin/arrai"]
