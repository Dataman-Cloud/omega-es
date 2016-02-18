#omega-es REST API
* search
 * [POST /search/index](#searchindex): 进入日志首页根据选择条件查询日志
 * [POST /search/jump](#searchjump): 根据一条日志查询日志上下文

##Search
####POST `/es/index`
日志首页根据选择条件查询
```shell
curl -X POST http://123.59.58.58:8080/es/index \
        -H Authorization:usertoken \
        -H Content-Type:application/json -d '{
    "userid": 1,
    "clusterid": 71,
    "appname": "afgsdfghsdf",
    "start": "2015-12-30T14:16:56.644+08:00",
    "end": "2015-12-30T14:19:56.643+08:00",
    "from": 0,
    "size": 20,
    "ipport": "10.3.11.18:[31757]",
    "keyword": "test"
}'
```
####POST `/es/content`
```shell
curl -X POST http://123.59.58.58:8080/es/content \
        -H Authorization:usertoken \
        -H Content-Type:application/json -d '{
    "userid": 1,
    "clusterid": 71,
    "appname": "htmltest",
    "timestamp": "2015-12-30T16:08:07.272+08:00",
    "ipport": "10.3.11.18:[31092]",
    "counter": 3
}'
```
####POST `/es/index/download`
```shell
curl -X POST http://10.3.11.22:9200/es/index/download/log.json \
        -H Authorization:usertoken \
        -H Content-Type:application/json -d '{
    "userid": 1,
    "clusterid": 76,
    "appname": "haha",
    "start": "2016-02-17T10:16:56.644+08:00",
    "end": "2016-02-17T14:19:56.643+08:00",
    "from": 0,
    "size": 20
}'
```
####POST `/es/content/download`
```shell
curl -X POST http://123.59.58.58:8080/es/content/download \
        -H Authorization:usertoken \
        -H Content-Type:application/json -d '{
    "userid": 1,
    "clusterid": 71,
    "appname": "htmltest",
    "timestamp": "2015-12-30T16:08:07.272+08:00",
    "ipport": "10.3.11.18:[31092]",
    "counter": 3
}'
```
