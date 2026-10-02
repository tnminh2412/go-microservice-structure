#FROM busybox AS health-builder

#ARG BUSYBOX_VERSION=1.31.0-i686-uclibc
#ADD https://busybox.net/downloads/binaries/$BUSYBOX_VERSION/busybox_WGET /wget
#RUN chmod a+x /wget

ARG REGISTRY
FROM $REGISTRY/runner-alpine
WORKDIR /
#COPY --from=health-builder /wget /usr/bin/wget
ARG SERVICE
ARG PORT
ARG ENV
ENV env=$ENV
ENV app_port=$PORT
COPY $SERVICE/config.$ENV.ini /config.$ENV.ini
COPY $SERVICE/bin /app
COPY $SERVICE/docs.html /docs.html
COPY $SERVICE/swagger.json /swagger.json
COPY health.sh /health.sh
RUN chmod +x /health.sh

HEALTHCHECK --interval=5s \
            --timeout=1s \
            CMD ["/health.sh"]
EXPOSE $PORT
CMD [ "/app", "serve" ]
