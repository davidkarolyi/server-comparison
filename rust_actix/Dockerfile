FROM rust:latest

WORKDIR /usr/src/app
COPY ./test_data.json .

WORKDIR /usr/src/app/rust_actix
COPY ./rust_actix/ .

RUN cargo build --release

EXPOSE 3000

CMD [ "./target/release/rust_actix"]