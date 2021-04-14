FROM golang:alpine3.12 AS stage
RUN apk add --no-cache make git
RUN go get github.com/anz-bank/go-bindata/...

WORKDIR /usr/arrai
COPY . .
RUN make build

FROM golang:alpine3.12
COPY --from=stage /usr/arrai/arrai /bin/arrai

ENTRYPOINT ["/bin/arrai"]
