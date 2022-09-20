FROM golang:1.17

RUN mkdir /go-img-resizer
WORKDIR /go-img-resizer
COPY . .

RUN make build

ENTRYPOINT ["./go-img-resizer"]