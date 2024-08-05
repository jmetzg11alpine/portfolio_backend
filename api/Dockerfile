FROM rust:latest

RUN cargo install cargo-watch

RUN useradd -ms /bin/bash rustacean

WORKDIR /home/rustacean

COPY --chown=rustacean:rustacean . .

RUN chown -R rustacean:rustacean /home/rustacean
RUN chown -R rustacean:rustacean /usr/local/cargo/registry

USER rustacean

EXPOSE 8080
