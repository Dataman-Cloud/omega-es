##omega-es环境变量说明

注释括号里面代表原来相对应配置文件的字段

```
OMEGAES_NET_HOST=0.0.0.0             	#服务监听地址(host)
OMEGAES_NET_PORT=5009                	#服务监听端口(port)
OMEGAES_NET_APPURL=localhost:8000    	#app服务api地址(appurl)
OMEGAES_LOG_CONSOLE=true				#日志是否标准输出(log.console)
OMEGAES_LOG_APPENDFILE=false			#日志是否输出到文件件(log.appendfile)
OMEGAES_LOG_FILE=./log/omega-es.log	#日志输出文件路径(log.file)
OMEGAES_LOG_LEVEL=debug					#日志输出级别(log.level)
OMEGAES_LOG_FORMATTER=text				#日志输出格式(log.formatter)
OMEGAES_LOG_MAXSIZE=1024000				#日志输出最大字节(log.maxSize)
OMEGAES_ES_HOSTS=localhost				#es地址多个用逗号分隔(es.hosts)
OMEGAES_ES_PORT=9200					#es端口(es.port)
OMEGAES_REDIS_HOST=localhost			#redis地址(redis.host)
OMEGAES_REDIS_PORT=6379					#redis端口(redis.port)
OMEGAES_MYSQL_HOST=localhost			#mysql地址(mysql.host)
OMEGAES_MYSQL_PORT=3306					#mysql端口(mysql.port)
OMEGAES_MYSQL_MAXIDLECONNS=5			#mysql最大闲置链接数(mysql.maxIdleConns)
OMEGAES_MYSQL_MAXOPENCONNS=50			#mysql最大打开链接数(mysq.maxOpenConns)
OMEGAES_MYSQL_DATABASE=alarm			#mysql数据库名(mysql.database)
OMEGAES_MYSQL_USERNAME=root				#mysql用户名(mysql.username)
OMEGAES_MYSQL_PASSOWRD=111111			#mysql用户密码(mysql.password)
```
