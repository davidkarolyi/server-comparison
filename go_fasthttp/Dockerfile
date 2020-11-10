FROM golang:1.15

WORKDIR /usr/src/app
COPY ./test_data.json .

WORKDIR /usr/src/app/go_fasthttp
COPY ./go_fasthttp/ .

RUN go build

EXPOSE 3000

CMD [ "./go_fasthttp"]