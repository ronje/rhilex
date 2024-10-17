package apis

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"
	"unicode/utf8"

	common "github.com/hootrhino/rhilex/component/apiserver/common"
	"github.com/hootrhino/rhilex/component/apiserver/model"
	"github.com/hootrhino/rhilex/component/apiserver/server"
	"github.com/hootrhino/rhilex/component/apiserver/service"
	"github.com/hootrhino/rhilex/component/internotify"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func InitUserRoute() {
	userApi := server.RouteGroup(server.ContextUrl("/users"))
	{
		userApi.GET(("/"), server.AddRoute(Users))
		userApi.POST(("/"), server.AddRoute(CreateUser))
		userApi.PUT(("/update"), server.AddRoute(UpdateUser))
		userApi.GET(("/detail"), server.AddRoute(UserDetail))
		userApi.POST(("/logout"), server.AddRoute(LogOut))
	}
}

const (
	__SECRET_KEY = "you-can-not-get-this-secret"
)

// All Users
type user struct {
	Role        string `json:"role"`
	Username    string `json:"username"`
	Description string `json:"description"`
}

func Users(c *gin.Context, ruleEngine typex.Rhilex) {
	users := []user{}
	for _, u := range service.AllMUser() {
		users = append(users, user{
			Role:        u.Role,
			Username:    u.Username,
			Description: u.Description,
		})
	}
	c.JSON(common.HTTP_OK, common.OkWithData(users))
}
func isLengthBetween6And16(str string) bool {
	length := utf8.RuneCountInString(str)
	return length >= 6 && length <= 16
}

// CreateUser
func CreateUser(c *gin.Context, ruleEngine typex.Rhilex) {
	type Form struct {
		Role        string `json:"role" binding:"required"`
		Username    string `json:"username" binding:"required"`
		Password    string `json:"password" binding:"required"`
		Description string `json:"description"`
	}
	form := Form{}
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	if !isLengthBetween6And16(form.Username) {
		c.JSON(common.HTTP_OK, common.Error("Username Length must Between 6 ~ 16"))
		return
	}
	if !isLengthBetween6And16(form.Password) {
		c.JSON(common.HTTP_OK, common.Error("Password Length must Between 6 ~ 16"))
		return
	}
	if _, err := service.GetMUser(form.Username); err != nil {
		service.InsertMUser(&model.MUser{
			Role:        form.Role,
			Username:    form.Username,
			Password:    md5Hash(form.Password),
			Description: form.Description,
		})
		c.JSON(common.HTTP_OK, common.Ok())
		return
	}
	c.JSON(common.HTTP_OK, common.Error("user already exists:"+form.Username))
}

// UpdateUser
func UpdateUser(c *gin.Context, ruleEngine typex.Rhilex) {
	type Form struct {
		Username    string `json:"username" binding:"required"`
		Password1   string `json:"password1" binding:"required"`
		Password2   string `json:"password2" binding:"required"`
		Description string `json:"description"`
	}
	form := Form{}
	if err1 := c.ShouldBindJSON(&form); err1 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err1))
		return
	}
	if !isLengthBetween6And16(form.Username) {
		c.JSON(common.HTTP_OK, common.Error("Username Length must Between 6 ~ 16"))
		return
	}
	if !isLengthBetween6And16(form.Password1) {
		c.JSON(common.HTTP_OK, common.Error("Password1 Length must Between 6 ~ 16"))
		return
	}
	if !isLengthBetween6And16(form.Password2) {
		c.JSON(common.HTTP_OK, common.Error("Password2 Length must Between 6 ~ 16"))
		return
	}
	_, errLogin := service.Login(form.Username, md5Hash(form.Password1))
	if errLogin != nil {
		c.JSON(common.HTTP_OK, common.Error("Invalid password"))
		return
	}

	token := c.GetHeader("Authorization")
	claims, errParse := parseToken(token)
	if errParse != nil {
		c.JSON(common.HTTP_OK, common.Error400(errParse))
		return
	}
	if err2 := service.UpdateMUser(claims.Username, &model.MUser{
		Username:    form.Username,
		Password:    md5Hash(form.Password2),
		Description: form.Description,
	}); err2 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err2))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())
}

/*
*
* Md5 计算
*
 */
func md5Hash(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

type LoginResultVo struct {
	Token      string       `json:"token"`
	UserDetail UserDetailVo `json:"userDetail"`
}

// Login
// TODO: 下个版本实现用户基础管理
func Login(c *gin.Context, ruleEngine typex.Rhilex) {
	type _user struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	clientIP := c.ClientIP()
	var u _user
	if err := c.BindJSON(&u); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	Ts := uint64(time.Now().UnixMilli())
	_, errLogin := service.Login(u.Username, md5Hash(u.Password))
	if errLogin != nil {
		glogger.GLogger.Warn("User Login Failed:", clientIP)
		internotify.Push(internotify.BaseEvent{
			Type:    `WARNING`,
			Event:   `event.system.user.login.failed`,
			Ts:      Ts,
			Summary: "User Login Failed",
			Info: fmt.Sprintf(`User Login Failed, Username: %s, RemoteAddr: %s`,
				u.Username, clientIP),
		})
		c.JSON(common.HTTP_OK, common.Error400(errLogin))
		return
	}
	token, err1 := generateToken(u.Username)
	if err1 != nil {
		glogger.GLogger.Warn("User Login Failed:", clientIP)
		internotify.Push(internotify.BaseEvent{
			Type:    `WARNING`, // INFO | ERROR | WARNING
			Event:   `event.system.user.login.failed`,
			Ts:      Ts,
			Summary: "User Login Failed",
			Info: fmt.Sprintf(`User Login Failed, Username: %s, RemoteAddr: %s`,
				u.Username, clientIP),
		})
		c.JSON(common.HTTP_OK, common.Error400(err1))
		return
	}
	glogger.GLogger.Info("User Login Success:", clientIP)
	internotify.Push(internotify.BaseEvent{
		Type:    `INFO`, // INFO | ERROR | WARNING
		Event:   `event.system.user.login.success`,
		Ts:      Ts,
		Summary: "User Login Success",
		Info: fmt.Sprintf(`User Login Success, Username: %s, RemoteAddr: %s`,
			u.Username, clientIP),
	})
	claims, err := parseToken(token)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	UserDetail := UserDetailVo{
		Username:       claims.Username,
		Role:           "ADMIN",
		Description:    "RHILEX ADMIN",
		BeginAuthorize: typex.License.BeginAuthorize,
		EndAuthorize:   typex.License.EndAuthorize,
		Type:           typex.License.Type,
		Token:          token,
	}
	c.JSON(common.HTTP_OK, common.OkWithData(UserDetail))

}

/*
*
* 退出
*
 */
func LogOut(c *gin.Context, ruleEngine typex.Rhilex) {
	c.JSON(common.HTTP_OK, common.Ok())
}

type UserDetailVo struct {
	BeginAuthorize int64  `json:"beginAuthorize"`
	Description    string `json:"description"`
	EndAuthorize   int64  `json:"endAuthorize"`
	Role           string `json:"role"`
	Token          string `json:"token"`
	Type           string `json:"type"`
	Username       string `json:"username"`
}

/*
*
* 用户信息
*
 */
func UserDetail(c *gin.Context, ruleEngine typex.Rhilex) {
	token := c.GetHeader("Authorization")
	claims, err := parseToken(token)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.OkWithData(UserDetailVo{
		Username:       claims.Username,
		Role:           "ADMIN",
		Description:    "RHILEX ADMIN",
		BeginAuthorize: typex.License.BeginAuthorize,
		EndAuthorize:   typex.License.EndAuthorize,
		Type:           typex.License.Type,
		Token:          token,
	}))
}

type JwtClaims struct {
	Username string
	jwt.StandardClaims
}

/*
*
* 生成Token
*
 */
func generateToken(username string) (string, error) {
	claims := &JwtClaims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(60*60*24*7) * time.Second).Unix(),
			Issuer:    username,
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(__SECRET_KEY))
	return token, err
}

/*
*
* 解析Token
*
 */
func parseToken(tokenString string) (*JwtClaims, error) {
	if tokenString == "" {
		return nil, fmt.Errorf("expected token string on headers")
	}

	token, err := jwt.ParseWithClaims(tokenString, &JwtClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(__SECRET_KEY), nil
		})
	if err != nil {
		return nil, err
	}
	if token.Claims == nil {
		return nil, fmt.Errorf("invalid Claims: %v", token.Raw)
	}
	if claims, ok := token.Claims.(*JwtClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
