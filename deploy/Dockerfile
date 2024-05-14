FROM --platform=linux/amd64 debian:bookworm-20230904-slim

COPY asm-live_timing-parser-server /app/server

EXPOSE 4242

CMD ["/app/server", "-c", "/config/config.toml"]
