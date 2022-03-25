## Usage memory auth plugin 
> add config section to conf/app.ini
``` ini
[auth]
; define auth plugin provider is redis
plugin = memory
ttl = 0
max_tokens = 0

; 0.AuthDiscardOld(recommended)
; 1.AuthDiscardNew(deny new connect)
; don't change this value,otherwise you really know what's this parame
discard = 0
[auth.public]
; public account is useful
enable = true
username = mqtt
password = public
; only public account supprted in memory provider
```