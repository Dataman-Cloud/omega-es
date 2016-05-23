FROM logstash:2.3.2

ADD dataman-2.3.2.conf /etc/logstash/conf.d/dataman-2.3.2.conf
ADD logstash-2.1.0.json /etc/logstash/conf.d/logstash-2.1.0.json

CMD ["logstash", "-f", "/etc/logstash/conf.d/dataman-2.3.2.conf"]
