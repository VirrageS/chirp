FROM logstash

# copy conf
WORKDIR /logstash
COPY logstash.conf .

# download jdbc postgres driver
RUN wget https://jdbc.postgresql.org/download/postgresql-9.4.1212.jar

# run logstash with given conf file
CMD ["-f", "logstash.conf"]