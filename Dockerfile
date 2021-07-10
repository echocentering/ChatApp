FROM postgres:latest
ENV POSTGRES_PASSWORD mysecretpassword
ENV POSTGRES_DB some-postgres
COPY world.sql /docker-entrypoint-initdb.d/