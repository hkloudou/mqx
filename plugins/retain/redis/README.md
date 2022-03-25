## Usage redis retain plugin
> add config section to conf/app.ini
``` ini
[retain]
; set retain plugin provider to redis
plugin = redis
[retain.plugin.redis]
; redis server
server = 127.0.0.1:6379
; redis db number
db = 3
; redis connect userName
username = 
; redis connect passWord
password = 
```