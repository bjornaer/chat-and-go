FROM golang:1.9

WORKDIR /go/src/app
COPY . .
WORKDIR /go/src/app/api
RUN go-wrapper download
RUN go-wrapper install
CMD ["go-wrapper", "run"]
