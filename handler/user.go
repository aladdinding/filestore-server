package handler

import (
	dblayer "filestore-server/db"
	"filestore-server/util"
	"fmt"
	"net/http"
	"time"
)

func SignupHandler() {

}

// UserInfoHandler ： 查询用户信息
func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	//1. 解析请求参数
	r.ParseForm()
	username := r.Form.Get("username")

	user, err := dblayer.GetUserInfo(username)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: user,
	}
	w.Write(resp.JSONBytes())
}

func IsTokenVaild(token string) bool {
	if len(token) != 40 {
		return false
	}
	// TODO:判断token的时效性，是否过期
	// TODO: 从数据库表tbl_user_token查询username对应的token信息
	// TODO: 对比两个token是否一致
	return true
}

func GenToken(username string) string {
	// 40位字符:md5(username=timestamp+token_salt)[:8]
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username + ts + "_tokensalt"))
	return tokenPrefix + ts[:8]

}
