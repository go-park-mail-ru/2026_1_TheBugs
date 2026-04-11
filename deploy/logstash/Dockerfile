FROM docker.elastic.co/logstash/logstash:8.12.0

RUN mkdir -p /usr/share/logstash/drivers && \
    curl -L -o /usr/share/logstash/drivers/postgresql.jar \
    https://jdbc.postgresql.org/download/postgresql-42.7.3.jar && \
    chown logstash:logstash /usr/share/logstash/drivers/postgresql.jar && \
    chmod 644 /usr/share/logstash/drivers/postgresql.jar