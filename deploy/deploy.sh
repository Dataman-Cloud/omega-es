#!/bin/bash

curl -v -X POST $MARATHON_API_URL/v2/apps -H Content-Type:application/json -d \
'{
      "id": "shurenyun-'$TASKENV'-'$SERVICE'",
      "cpus": '$CPUS',
      "mem": '$MEM',
      "instances": '$INSTANCES',
      "constraints": [["hostname", "LIKE", "'$DEPLOYIP'"], ["hostname", "UNIQUE"]],
      "container": {
                     "type": "DOCKER",
                     "docker": {
                                     "image": "'$SERVICE_IMAGE'",
                                     "network": "BRIDGE",
				     "privileged": '$PRIVILEGED',
				     "forcePullImage": '$FORCEPULLIMAGE',
				     "portMappings": [
                                             { "containerPort": '$BAMBOO_TCP_PORT', "hostPort": 0, "protocol": "tcp"}
                                     ]
                                }
                   },
      "env": {
		    "BAMBOO_TCP_PORT": "'$BAMBOO_TCP_PORT'",
            "BAMBOO_PRIVATE": "'$BAMBOO_PRIVATE'",
            "BAMBOO_PROXY":"'$BAMBOO_PROXY'",
		    "BAMBOO_BRIDGE": "'$BAMBOO_BRIDGE'",
		    "BAMBOO_HTTP": "'$BAMBOO_HTTP'"
             },
      "uris": [
               "'$CONFIGSERVER'/config/demo/config/registry/docker.tar.gz"
       ]
}'
