FROM alpine:3

RUN addgroup -g 1001 -S app && \
    adduser -u 1001 -S app -G app

WORKDIR /app

COPY --chown=app:app zendo-api /app

RUN touch .env # init empty env

USER app

ENTRYPOINT ["/app/zendo-api"]
