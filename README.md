# mqx

# auth pugin
> Authentication is an important part of most applications. MQTT protocol supports username/password authentication. Enabling authentication can effectively prevent illegal client connections.

- redis `release`
https://github.com/hkloudou/mqx/tree/main/plugins/auth/redis
- http `future`
https://github.com/hkloudou/mqx/tree/main/plugins/auth/http
- jwt `future`
https://github.com/hkloudou/mqx/tree/main/plugins/auth/jwt
- ldap `future`
https://github.com/hkloudou/mqx/tree/main/plugins/auth/ldap
- mongo `future`
https://github.com/hkloudou/mqx/tree/main/plugins/auth/mongo
- mysql `future`
https://github.com/hkloudou/mqx/tree/main/plugins/auth/mysql
- postgresql `future`
https://github.com/hkloudou/mqx/tree/main/plugins/auth/postgresql

# retain plugin
> When the server receives a PUBLISH packet with a Retain flag of 1, it will treat the message as a retained message. In addition to being normally forwarded, the retained message will be stored on the server. There can only be one retained message under each topic. Therefore, if there is already a retained message of the same topic, the retained message is replaced.

- redis `release`
https://github.com/hkloudou/mqx/tree/main/plugins/retain/redis