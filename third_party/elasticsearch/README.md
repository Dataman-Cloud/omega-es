# elasticsearch server with ik plugin

Docker run command

node1
```
docker run --rm --net=host \
   -e NODES="127.0.0.1:9300,127.0.0.1:9301,127.0.0.1:9302" \
   -e HOSTNAME=node1 -e HTTPPORT=9200 -e TRANSPORT=9300 \
   -it demoregistry.dataman-inc.com/srypoc/elastciser-2.1.0:20160608
```

node2
```
docker run --rm --net=host \
   -e NODES="127.0.0.1:9300,127.0.0.1:9301,127.0.0.1:9302" \
   -e HOSTNAME=node2 -e HTTPPORT=9201 -e TRANSPORT=9301 \
   -it demoregistry.dataman-inc.com/srypoc/elastciser-2.1.0:20160608
```

node3
```
docker run --rm --net=host \
   -e NODES="127.0.0.1:9300,127.0.0.1:9301,127.0.0.1:9302" \
   -e HOSTNAME=node3 -e HTTPPORT=9202 -e TRANSPORT=9302 \
   -it demoregistry.dataman-inc.com/srypoc/elastciser-2.1.0:20160608
```
