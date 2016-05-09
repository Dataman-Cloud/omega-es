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
    "keyword": "test",
    "source": "echo"
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
    "counter": 3,
    "source": "echo"
}'
```
####POST `/es/index/download`
```shell
curl -X POST http://10.3.11.22:9200/es/index/download \
        -H Authorization:usertoken \
        -H Content-Type:application/json -d '{
    "userid": 1,
    "clusterid": 71,
    "appname": "afgsdfghsdf",
    "start": "2015-12-30T14:16:56.644+08:00",
    "end": "2015-12-30T14:19:56.643+08:00",
    "ipport": "10.3.11.18:[31757]",
    "keyword": "test",
    "source": "echo"
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
    "counter": 3,
    "source": "echo"
}'
```

####POST `/es/alarm/create`
创建报警策略
```shell
curl -X POST http://10.3.20.53:5009/es/alarm/create \
	-H Authorization:usertoken \
	-H Content-Type:application/json -d '{
    "userid": 1,
    "clusterid": 1,
    "appname": "testalarm",
    "interval": 5,
    "gtnum": 10,
    "alarmname": "alarmtest",
    "usertype": "user",
    "keyword": "error",
    "emails": "yqguo@dataman-inc.com"
}'
```

####PUT `/es/alarm/update`
更新报警策略
```shell
curl -X PUT http://10.3.20.53:5009/es/alarm/update \
	-H Authorization:usertoken \
	-H Content-Type:application/json -d '{
    "id":1,
    "clusterid": 1,
    "appname": "testalarm",
    "appalias": "test",
    "interval": 5,
    "gtnum": 10,
    "usertype": "user",
    "keyword": "error",
    "emails": "yqguo@dataman-inc.com"
}'
```
####GET `/es/alarm/:id`
获取策略详细信息
`curl -X -H Authorization:usertoken GET http://10.3.20.53:5009/es/alarm/:id`

####PUT `/es/alarm/stop/:id`
停止策略
`curl -X -H Authorization:usertoken PUT http://10.3.20.53:5009/es/alarm/stop/:id`

####DELETE `/es/alarm/delete/:id`
删除报警策略
`curl -X -H Authorization:usertoken DELETE http://10.3.20.53:5009/es/alarm/delete/:id`

####GET `/es/alarm/list?usertype=usertype&uid=uid&pcount=pcount&pnum=pnum`
查看创建策略列表
`curl -X -H Authorization:usertoken GET http://10.3.20.53:5009/es/alarm/list?usertype=usertype&uid=uidpcount=pcount&pnum=pnum`

####GET `/es/alarm/scheduler/history?jobid=jobid&pcount=pcount&pnum=pnum`
查看策略执行历史记录
`curl  -X -H Authorization:usertoken GET http://10.3.20.53:5009/es/alarm/scheduler/history?id=id&pcount=pcount&pnum=pnum`
