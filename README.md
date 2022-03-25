# mqx (developing, not stable)


# session plugin
- memory `default`
https://github.com/hkloudou/mqx/tree/main/plugins/session/memory
# auth pugin
> Authentication is an important part of most applications. MQTT protocol supports username/password authentication. Enabling authentication can effectively prevent illegal client connections.

- redis `default`
https://github.com/hkloudou/mqx/tree/main/plugins/auth/redis
<!-- - http `todo`
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
https://github.com/hkloudou/mqx/tree/main/plugins/auth/postgresql -->

# retain plugin
> When the server receives a PUBLISH packet with a Retain flag of 1, it will treat the message as a retained message. In addition to being normally forwarded, the retained message will be stored on the server. There can only be one retained message under each topic. Therefore, if there is already a retained message of the same topic, the retained message is replaced.

- memory`default`
https://github.com/hkloudou/mqx/tree/main/plugins/retain/memory
- redis
https://github.com/hkloudou/mqx/tree/main/plugins/retain/redis

# TODO LIST
## interface
- [x] conf interface
- [x] auth interface
- [x] retain interface
- [x] session interface
- [ ] acl interface
## provider
- [x] redis auth provider
- [x] memory session provider
- [ ] redis session provider
- [x] memory retain provider
- [x] redis retain provider
## other
- [ ] $sys and $usr message support