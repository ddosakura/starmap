package handler

import (
	"time"

	proto "github.com/ddosakura/starmap/srv/auth/proto"
	"github.com/ddosakura/starmap/srv/auth/raw"
	"github.com/ddosakura/starmap/srv/common"
	jwt "github.com/dgrijalva/jwt-go"
)

//jwt.StandardClaims{
//	// Id: "",
//	Issuer: "starmap", // 签发方
//	// Subject: "", // 面向用户
//	// Audience: "", // 接收方
//	IssuedAt: iat.Unix(), // 签发时间
//	// NotBefore: 0, // 生效时间
//	ExpiresAt: expireTime, // 过期时间
//}

// UserJWT with UserInfo
type UserJWT struct {
	jwt.StandardClaims
	UserInfo *proto.UserInfo
}

func buildUserJWT(u *proto.UserInfo) *UserJWT {
	return &UserJWT{
		StandardClaims: jwt.StandardClaims{},
		UserInfo:       u,
	}
}

// sign JWT
func (c *UserJWT) sign(d time.Duration) (string, error) {
	iat := time.Now()
	expireTime := iat.Add(d).Unix()
	c.Issuer = "starmap"
	c.IssuedAt = iat.Unix()
	c.ExpiresAt = expireTime
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return jwtToken.SignedString(common.ENV.KeyJWT())
}

// ValidUserJWT for UserInfo
func ValidUserJWT(tokenStr string) (t *UserJWT, e error) {
	token, err := jwt.ParseWithClaims(tokenStr, &UserJWT{}, func(token *jwt.Token) (interface{}, error) {
		return common.ENV.KeyJWT(), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, raw.ErrInvalidJWT
	}

	if t, ok := token.Claims.(*UserJWT); ok {
		return t, nil
	}
	return nil, raw.ErrInvalidJWT
}
