## Usage memory acl plugin 
> add config acl to conf/app.ini
``` ini
[acl]
; define session plugin provider is memory
plugin = memory



; allow all user subscribe public topic
[acl.plugin.memory.sub.allow.0]
patterns = "#"

; only all user subscribe self private topic
[acl.plugin.memory.sub.allow.1]
patterns = "$$usr/<username>/#,$$cid/<username>/#"

; admin can subscribe all private topic
[acl.plugin.memory.sub.allow.2]
patterns = "$$usr/#,$$cid/#"
user = admin


; mqtt user is test user,deny to subscribe private topic
; [acl.plugin.memory.sub.deny.0]
; patterns = "$$usr/#,$$cid/#"
; user = mqtt

; local ip user mqtt can publish public topic
[acl.plugin.memory.pub.allow.0]
patterns = "#"
cidr = "127.0.0.1/16"
user = mqtt

; local ip user mqtt can publish private topic
[acl.plugin.memory.pub.allow.1]
patterns = "$$usr/#"
cidr = "127.0.0.1/16"
user = mqtt
```