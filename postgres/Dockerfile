FROM postgres:9.5

RUN apt-get update
RUN apt-get install -y git build-essential postgresql libpq-dev postgresql-server-dev-all postgresql-common
RUN git clone https://github.com/keithf4/pg_partman.git
RUN cd pg_partman && make NO_BGW=1 install

COPY schema.sql /
COPY fixtures.sql /
COPY init.sh /docker-entrypoint-initdb.d/
