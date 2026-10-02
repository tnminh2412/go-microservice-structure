#FROM busybox AS health-builder

#ARG BUSYBOX_VERSION=1.31.0-i686-uclibc
#ADD https://busybox.net/downloads/binaries/$BUSYBOX_VERSION/busybox_WGET /wget
#RUN chmod a+x /wget


#FROM overclockedllama/docker-chromium
ARG REGISTRY
FROM $REGISTRY/chromium
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

EXPOSE $PORT
CMD [ "/app", "serve" ]

