package handler

import (
	"testing"

	"github.com/kr/pretty"

	"time"

	proto "github.com/ddosakura/starmap/srv/auth/proto"
)

// todo this test, use `CONFIG_KEY_JWT=123456 go test -v`
func TestJWT(t *testing.T) {
	u := &proto.UserInfo{
		Nickname: "Hello World!",
		Homepage: "www.baidu.com",
	}
	TOKEN := buildUserJWT(u)
	pretty.Println(TOKEN)

	jwt, e := TOKEN.sign(time.Second * 3)
	if e != nil {
		t.Fatal(e)
		return
	}
	t.Log(jwt)

	token, e := ValidUserJWT(jwt)
	if e != nil {
		t.Fatal(e)
		return
	}
	pretty.Println(token)

	time.Sleep(time.Second * 5)

	token, e = ValidUserJWT(jwt)
	if e == nil {
		t.Fatal("JWT 未过期")
	} else {
		t.Log("过期成功")
	}
	pretty.Println(token)
}
