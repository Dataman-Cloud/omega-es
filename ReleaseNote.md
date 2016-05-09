#Release Note

更新版本v0.1.050900
##需要运维更新
 * 修改nginx配置文件
```
在原有的nginx配置文件基础上增减update stop restart 等
location ~ /es/alarm/(create|list|delete|update|stop|restart|scheduler/history|\d) {
  if ($request_method = OPTIONS ) {
      add_header Access-Control-Allow-Origin "*" ;
      add_header Access-Control-Allow-Methods "GET,PUT,POST,DELETE,OPTIONS,PATCH";
      add_header Access-Control-Allow-Headers "Content-Type, Depth, User-Agent, X-File-Size, X-Requested-With, X-Requested-By, If-Modified-Since, X-File-Name, Cache-Control, X-XSRFToken, Authorization" ;
      add_header Access-Control-Allow-Credentials "true" ;
      add_header Content-Length 0 ;
      add_header Content-Type application/json ;
      return 204;
  }
  if ($request_method != 'OPTIONS') {
      add_header 'Access-Control-Allow-Origin' '*' always;
      add_header 'Access-Control-Allow-Credentials' 'true' always;
      add_header 'Access-Control-Allow-Methods' 'GET,PUT,POST,DELETE,OPTIONS,PATCH' always;
      add_header 'Access-Control-Allow-Headers' 'DNT,X-CustomHeader,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type' always;
  }
  auth_request    /_auth;
  proxy_pass      http://10.3.20.53:5009;
}
```
