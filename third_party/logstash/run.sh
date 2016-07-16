#!/bin/bash

CONF=/etc/logstash/conf.d/dataman-2.3.2.conf

sed -i s/ELASTICSEARCH/$ELASTICSEARCH/g $CONF 
sed -i s/PORT/$PORT/g $CONF 

until curl -X PUT http://$ELASTICSEARCH/_template/dataman -d '
{
    "order": 0,
    "template": "dataman-*",
    "settings": {
        "index": {
            "query": {
                "default_field": "@message"
            },
            "store": {
                "compress": {
                    "stored": true,
                    "tv": true
                }
            }
        }
    },
    "mappings": {
        "_default_": {
            "dynamic_templates": [
                {
                    "string_template": {
                        "match": "*",
                        "mapping": {
                            "type": "string",
                            "index": "not_analyzed"
                        },
                        "match_mapping_type": "string"
                    }
                }
            ],
            "_ttl": {
                "enabled": true,
                "default": "7d"
            },
            "properties": {
                "clusterid": {
                    "type": "string",
                    "index": "not_analyzed"
                },
                "ip": {
                    "type": "string",
                    "index": "not_analyzed"
                },
                "hostname": {
                    "type": "string",
                    "index": "not_analyzed"
                },
                "type": {
                    "type": "string",
                    "index": "not_analyzed"
                },
                "userid": {
                    "type": "string",
                    "index": "not_analyzed"
                },
                "counter": {
                    "type": "long",
                    "index": "not_analyzed"
                },
                "timestamp": {
                    "type": "multi_field",
                    "fields": {
                        "timestamp": {
                            "type": "date",
                            "index": "not_analyzed"
                        },
                        "sort": {
                            "type": "string",
                            "index": "not_analyzed"
                        }
                    }
                },
                "message": {
                    "type": "string",
                    "index": "not_analyzed"
                },
                "uuid": {
                    "type": "string",
                    "index": "not_analyzed"
                },
                "ipport": {
                    "type": "string",
                    "index": "not_analyzed"
                },
                "taskid": {
                    "type": "string",
                    "index": "not_analyzed"
                },
                "platform": {
                    "type": "string",
                    "index": "not_analyzed"
                },
                "source": {
                    "type": "string",
                    "index": "not_analyzed"
                },
                "time": {
                    "type": "string",
                    "index": "not_analyzed"
                },
                "appname": {
                    "type": "string",
                    "index": "not_analyzed"
                },
                "msg": {
                    "type": "string",
                    "analyzer": "ik",
                    "store": "yes"
                }
            },
            "_all": {
                "enabled": true
            },
            "_source": {
                "compress": true
            }
        }
    },
    "aliases": {}
}'
do
	sleep 1
done

logstash -f $CONF
