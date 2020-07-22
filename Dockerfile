FROM golang:alpine3.12
RUN apk add --no-cache make git
WORKDIR /usr/arrai
COPY . .
RUN make install
ENTRYPOINT ["arrai"]
