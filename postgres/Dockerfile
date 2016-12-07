FROM postgres:latest

COPY schema.sql /
COPY fixtures.sql /
COPY init.sh /docker-entrypoint-initdb.d/
