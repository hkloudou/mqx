## redis session plugin
``` ini
[session]
; define retain plugin provider is redis
plugin = memory
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