FROM alpine:3.16.2

WORKDIR app

# TODO: Make sure you build the app binary before building the docker container
# IMPORTANT!!! `strings` is a pre-existing binary in the alpine image
# so we need to rename our binary to something else
COPY build/strings-linux-amd64 /bin/strings-linux-amd64

RUN mkdir sql

COPY infrastructure/postgres/sql/ sql/

ENV DB_HOST=host.docker.internal
ENV SQL_MIGRATION_SCRIPTS=/app/sql
ENV GIN_MODE=release

CMD strings-linux-amd64