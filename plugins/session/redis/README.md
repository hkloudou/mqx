## Usage redis session plugin
> add config section to conf/app.ini
``` ini
[session]
; define session plugin provider is redis
plugin = redis
[session.plugin.redis]
; redis server
server = 127.0.0.1:6379
; redis db number
db = 3
; redis connect userName
username = 
; redis connect passWord
password = 
```