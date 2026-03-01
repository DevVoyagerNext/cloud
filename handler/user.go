package handler

import (
	dblayer "cloud/db"
	"cloud/util"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	pwd_salt = "*#890"
)

// SignupHandler 处理用户注册请求
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		data, err := ioutil.ReadFile("./static/view/signup.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
		return
	}
	r.ParseForm()
	username := r.Form.Get("username")
	passwd := r.Form.Get("password")

	if len(username) < 3 || len(passwd) < 5 {
		w.Write([]byte("Invalid username or password"))
		return
	}

	enc_password := util.Sha1([]byte(passwd + pwd_salt))
	suc := dblayer.UserSignup(username, enc_password)
	if suc {
		w.Write([]byte("User signed up"))
	} else {
		w.Write([]byte("User not signed up"))
	}
}

// SignInHandler 登录接口
func SignInHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")
	passwd := r.Form.Get("password")
	encPassword := util.Sha1([]byte(passwd + pwd_salt))
	// 1.校验用户名及密码
	pwdChecked := dblayer.UserSignin(username, encPassword)
	if !pwdChecked {
		w.Write([]byte("Invalid password"))
		return
	}
	// 2.生成访问凭证（token）
	token := GenToken(username)
	upRes := dblayer.UpdateToken(username, token)
	if !upRes {
		w.Write([]byte("Failed"))
		return
	}
	// 3. 登录成功后重定向到首页
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: struct {
			Location string
			Username string
			Token    string
		}{
			Location: "http://" + r.Host + "/static/view/home.html",
			Username: username,
			Token:    token,
		},
	}
	w.Write(resp.JSONBytes())
}

func GenToken(username string) string {
	// 40位字符的token:md5(username+timestamp+token_salt)+timestamp[:8]
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username + ts + "_tokensalt"))
	return tokenPrefix + ts[:8]
}

func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	//1.解析请求参数
	r.ParseForm()
	username := r.Form.Get("username")
	//token := r.Form.Get("token")

	////2.验证token是否有效
	//IsTokenValid := IsTokenValid(token)
	//if !IsTokenValid {
	//	w.WriteHeader(http.StatusForbidden)
	//	return
	//}
	//3.查询用户信息
	user, err := dblayer.GetUserInfo(username)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	//4.组装并且响应用户数据
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: user,
	}
	w.Write(resp.JSONBytes())
}

// IsTokenValid token是否有效
func IsTokenValid(token string) bool {
	if len(token) != 40 {
		return false
	}
	//TODO:判断token时效性,是否过期
	//TODO:从数据库里查询是否有这个token
	//TODO:比较token是否一致
	return true
}
