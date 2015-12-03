#omega-es REST API
* search
 * [POST /search/index](#searchindex): 进入日志首页根据选择条件查询日志
 * [POST /search/jump](#searchjump): 根据一条日志查询日志上下文

##Search
####POST `/search/index`
日志首页根据选择条件查询
```shell
curl -X POST http://123.59.58.58:8080/search/index \
        -H Authorization:usertoken \
        -H Content-Type:application/json -d '{
        	"clusterid": 19, 
            "appname": "chronos", 
            "hosts":["10.3.11.2"], 
            "start": "2015-11-01T15:17:24", 
            "end": "2015-11-01T15:17:24", 
            "from":0, 
            "size":20, 
            "userid":1,
            "keyword": "test"
       }'
```
####POST `/search/jump`
```shell
curl -X POST http://123.59.58.58:8080/search/jump \
        -H Authorization:usertoken \
        -H Content-Type:application/json -d '{
            "userid":1,
        	"clusterid": 19, 
            "taskid":"adfasd",
            "start": "2015-11-01T15:17:24", 
            "end": "2015-11-01T15:17:24", 
            "appname": "chronos", 
            "from":0, 
            "size":20
       }'
```

