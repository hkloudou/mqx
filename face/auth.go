package face

import (
	"context"
	"errors"
	"time"
)

type ErrAuthInvalid error

var (
	ErrAuthInvalidUserName ErrAuthInvalid = errors.New("mqx: auth failed, invalid userName")
	ErrAuthInvalidPassword ErrAuthInvalid = errors.New("mqx: auth failed, invalid password")
	ErrAuthInvalidClientIP ErrAuthInvalid = errors.New("mqx: auth failed, invalid clientIP")
	ErrAuthInvalidExpired  ErrAuthInvalid = errors.New("mqx: auth failed, token expired")
)

type AuthRequest struct {
	ClientId       string `json:"clientId"`       // reCommended to use deviceID
	UserName       string `json:"userName"`       // reCommended to use userID `uint64`
	PassWord       string `json:"passWord"`       // reCommended to use temporary password or http sessionID token
	TlsServerName  string `json:"tlsServerName"`  // tls ShakeHande ServerName
	TlsSubjectName string `json:"tlsSubjectName"` // tls PeerCertificates subject.CommonName
	ClientIp       string `json:"clientIp"`       // real clientip
}

type AuthReply struct {
	UserName  string `json:"userName"`
	ExpiredAt int64  `json:"expiredAt"` // 0 if never expire
}

// Auth interface
// default public mqtt account is mqtt:public
type Auth interface {
	// Init() error
	GlobalConfig(options ...authOption) error
	// call the function in your application
	Update(ctx context.Context, req *AuthRequest, ttl time.Duration) error
	// //
	// Delete(req *AuthRequest) error
	// // when the mqtt broker receive a mqtt.Connect packet
	Check(ctx context.Context, req *AuthRequest, ttl time.Duration) (bool, error)
}
