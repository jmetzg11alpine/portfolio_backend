FROM rust:latest

ARG DATABASE_URL
ENV DATABASE_URL=${DATABASE_URL}

ARG RUST_BACKTRACE
ENV RUST_BACKTRACE=${RUST_BACKTRACE}

# Install cargo-watch and sqlx-cli as root
RUN cargo install cargo-watch \
    && cargo install sqlx-cli --no-default-features --features rustls,mysql

# Create the user
RUN useradd -ms /bin/bash rustacean

# Set permissions for the rustacean user
RUN mkdir -p /home/rustacean/.cargo \
    && chown -R rustacean:rustacean /home/rustacean \
    && chown -R rustacean:rustacean /usr/local/cargo/registry \
    && chown -R rustacean:rustacean /usr/local/cargo/bin

# Switch to rustacean user
USER rustacean

WORKDIR /home/rustacean

# Copy the application code
COPY --chown=rustacean:rustacean . .

EXPOSE 8080
