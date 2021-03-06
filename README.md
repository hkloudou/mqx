# mqx (developing, not stable)

# TODO LIST
## interface
- [x] conf interface
- [x] auth interface
- [x] retain interface
- [x] session interface
- [x] acl interface
## acl privider
> **Publish/Subscribe** ACL refers to **permission control** for **PUBLISH/SUBSCRIBE** operations.
- [x] memory `default` https://github.com/hkloudou/mqx/tree/main/plugins/acl/memory
## auth provider
> Authentication is an important part of most applications. MQTT protocol supports username/password authentication. Enabling authentication can effectively prevent illegal client connections.
- [x] memory `default` https://github.com/hkloudou/mqx/tree/main/plugins/auth/memory
- [x] redis https://github.com/hkloudou/mqx/tree/main/plugins/auth/redis
- [ ] http
- [ ] jwt
- [ ] ldap
- [ ] mysql
- [ ] mongo
- [ ] postgresql
- [ ] nrpc(https://github.com/hkloudou/nrpc)
## session provider
- [x] memory `default`
https://github.com/hkloudou/mqx/tree/main/plugins/session/memory
- [x] redis https://github.com/hkloudou/mqx/tree/main/plugins/session/redis
- [ ] disk
- [ ] nrpc(https://github.com/hkloudou/nrpc)
## retain provider
> When the server receives a PUBLISH packet with a Retain flag of 1, it will treat the message as a retained message. In addition to being normally forwarded, the retained message will be stored on the server. There can only be one retained message under each topic. Therefore, if there is already a retained message of the same topic, the retained message is replaced.
- [x] memory `default`
https://github.com/hkloudou/mqx/tree/main/plugins/retain/memory
- [x] redis https://github.com/hkloudou/mqx/tree/main/plugins/retain/redis
- [ ] disk
- [ ] etcd
- [ ] nrpc(https://github.com/hkloudou/nrpc)
- [ ] s3
## other
- [ ] $SYS message
> system message
- [ ] $SHARE message
> shared message `$shared/group/topic`
- [ ] $USR message 
> user message `$USR/username/topic`
- [ ] $CID message
> cid message `$CID/clientid/topic`