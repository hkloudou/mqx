## redis auth plugin
``` ini
[auth]
; define auth plugin provider is redis
plugin = redis
ttl = 0
max_tokens = 0

[auth.public]
; public account is useful
enable = true
username = mqtt
password = public

[auth.plugin.redis]
; redis server
server = 127.0.0.1:6379
; redis db number
db = 3
; redis connect userName
username = 
; redis connect passWord
password = 
```