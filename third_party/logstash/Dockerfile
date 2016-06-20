FROM logstash:2.3.2

RUN && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
env ELASTICSEARCH elasticsearch:9200
env PORT 4999
ADD dataman-2.3.2.conf /etc/logstash/conf.d/dataman-2.3.2.conf
ADD logstash-2.1.0.json /etc/logstash/conf.d/logstash-2.1.0.json
ADD run.sh /run.sh

CMD /run.sh
