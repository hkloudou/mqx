package face

import (
	"context"
	"errors"
)

type ErrAuthInvalid error

var (
	ErrAuthInvalidUserNamePassword ErrAuthInvalid = errors.New("mqx: auth failed, Invalid username or password")
	ErrAuthServiceUnviable         ErrAuthInvalid = errors.New("mqx: auth faild, Service unavailable")
	ErrAuthInvalidClientIP         ErrAuthInvalid = errors.New("mqx: auth failed, Invalid clientIP")
	ErrAuthInvalidExpired          ErrAuthInvalid = errors.New("mqx: auth failed, Token expired")
	ErrAuthInvalidTooManyTokens    ErrAuthInvalid = errors.New("mqx: auth failed, too many tokens")
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
	// GlobalConfig(options ...authOption) error

	// call the function in your application,ttl -1: loginout
	Update(ctx context.Context, req *AuthRequest, options ...AuthRequestOption) error
	// //
	// Delete(req *AuthRequest) error
	// // when the mqtt broker receive a mqtt.Connect packet
	Check(ctx context.Context, req *AuthRequest, options ...AuthRequestOption) (bool, error)
}
