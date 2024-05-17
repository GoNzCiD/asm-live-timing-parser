FROM --platform=linux/amd64 debian:bookworm-20230904-slim

# Copy first the server binary, the templates folder and the assets folder into the current Dockerfile 
# directory to put these data into the context and use `docker build . -t image_name:image_tag`
COPY asm-live_timing-parser-server /app/server
COPY templates /app/templates
COPY assets /app/assets

WORKDIR /app

EXPOSE 4242

CMD ["/app/server", "-c", "/config/config.toml"]
