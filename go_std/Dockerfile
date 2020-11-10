FROM golang:1.15

WORKDIR /usr/src/app
COPY ./test_data.json .

WORKDIR /usr/src/app/go_std
COPY ./go_std/ .

RUN go build

EXPOSE 3000

CMD [ "./go_std"]