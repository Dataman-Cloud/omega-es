#!/bin/bash

CONF=/etc/logstash/conf.d/dataman-2.3.2.conf

sed -i s/ELASTICSEARCH/$ELASTICSEARCH/g $CONF 
sed -i s/PORT/$PORT/g $CONF 

logstash -f $CONF
