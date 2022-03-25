# mqx (developing, not stable)


# session plugin

# auth pugin
> Authentication is an important part of most applications. MQTT protocol supports username/password authentication. Enabling authentication can effectively prevent illegal client connections.

- redis `done`
https://github.com/hkloudou/mqx/tree/main/plugins/auth/redis
- http `todo`
https://github.com/hkloudou/mqx/tree/main/plugins/auth/http
- jwt `todo`
https://github.com/hkloudou/mqx/tree/main/plugins/auth/jwt
- ldap `todo`
https://github.com/hkloudou/mqx/tree/main/plugins/auth/ldap
- mongo `todo`
https://github.com/hkloudou/mqx/tree/main/plugins/auth/mongo
- mysql `todo`
https://github.com/hkloudou/mqx/tree/main/plugins/auth/mysql
- postgresql `todo`
https://github.com/hkloudou/mqx/tree/main/plugins/auth/postgresql

# retain plugin
> When the server receives a PUBLISH packet with a Retain flag of 1, it will treat the message as a retained message. In addition to being normally forwarded, the retained message will be stored on the server. There can only be one retained message under each topic. Therefore, if there is already a retained message of the same topic, the retained message is replaced.

- momory `done`
https://github.com/hkloudou/mqx/tree/main/plugins/retain/redis
- redis `done`
https://github.com/hkloudou/mqx/tree/main/plugins/retain/redis

# TODO LIST
- [x] conf interface
- [x] auth interface
- [x] retain interface
- [ ] acl interface
- [ ] store interface let qos>0 infight
- [ ] store interface to store psub session
- [ ] $sys and $usr message support