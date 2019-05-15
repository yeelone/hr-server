package user

import (
	h "hrgdrc/handler"
	"hrgdrc/model"
	"hrgdrc/pkg/auth"
	"hrgdrc/pkg/errno"
	"hrgdrc/pkg/token"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
)

// @Summary Login generates the authentication token
// @Produce  json
// @Param username body string true "Username"
// @Param password body string true "Password"
// @Success 200 {string} json "{"code":0,"message":"OK","data":{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1MjgwMTY5MjIsImlkIjowLCJuYmYiOjE1MjgwMTY5MjIsInVzZXJuYW1lIjoiYWRtaW4ifQ.LjxrK9DuAwAzUD8-9v43NzWBN7HXsSLfebw92DKd1JQ"}}"
// @Router /login [post]
func Login(c *gin.Context) {
	// Binding the data with the user struct.
	var r CreateRequest
	if err := c.Bind(&r); err != nil {
		h.SendResponse(c, errno.ErrBind, nil)
		return
	}

	// Get the user information by the login username.
	d, err := model.GetUserByName(r.Username)
	if err != nil {
		log.Warnf("Login function called. ID: %d | username: %s | info: %s  ", 0, r.Username, "user trying to login")
		h.SendResponse(c, errno.ErrUserNotFound, err.Error())
		return
	}

	// Compare the login password with the user password.
	if err := auth.Compare(d.Password, r.Password); err != nil {
		log.Warnf("Login function called. ID: %d | username: %s | info: %s  ", d.ID, d.Username, "password incorrect")
		h.SendResponse(c, errno.ErrPasswordIncorrect, err.Error())
		return
	}

	role := ""
	if len(d.Roles) > 0 {
		role = d.Roles[0].Name
	}
	// Sign the json web token.
	t, err := token.Sign(c, token.Context{ID: d.ID, Username: d.Username, Role: role}, "")
	if err != nil {
		log.Warnf("Login function called. ID: %d | username: %s | info: %s  ", d.ID, d.Username, "token error")
		h.SendResponse(c, errno.ErrToken, err.Error())
		return
	}

	log.Infof("Login function called. ID: %d | username: %s | info: login success", d.ID, d.Username)
	//返回给客户端之前把密码抹除
	d.Password = ""
	c.Set("CurrentUsername", d.Username)
	c.Set("CurrentUserID", d.ID)
	h.SendResponse(c, nil, model.Token{
		Token: t,
		User:  d,
	})
}

func Logout(c *gin.Context) {
	h.SendResponse(c, nil, "Successfully logged out")
}
