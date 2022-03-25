## Usage memory auth plugin 
> add config section to conf/app.ini
``` ini
[auth]
; define auth plugin provider is redis
plugin = memory
[auth.public]
; public account is useful
enable = true
username = mqtt
password = public
; only public account supprted in memory provider
```