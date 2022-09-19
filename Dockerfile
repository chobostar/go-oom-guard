FROM postgres:14

ENV POSTGRES_HOST_AUTH_METHOD=trust
ENV POSTGRES_PASSWORD=password

RUN apt-get update && apt-get install -y procps lsof && rm -rf /var/lib/apt/lists/*

COPY ./bin/go-oom-guard .