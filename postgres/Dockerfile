FROM postgres:latest

COPY schema.sql /docker-entrypoint-initdb.d/
COPY fixtures.sql /docker-entrypoint-initdb.d/
