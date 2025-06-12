FROM rust:1.87.0 AS builder
RUN rustup target add x86_64-unknown-linux-musl
RUN apt update && apt install -y musl-tools musl-dev
RUN update-ca-certificates
WORKDIR /
COPY ./ .
RUN cargo build --target x86_64-unknown-linux-musl --release

FROM scratch
COPY --from=builder /target/x86_64-unknown-linux-musl/release/ifconfig ./
ENTRYPOINT ["/ifconfig"]
