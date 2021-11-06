package controller

import (
	"demo/models"
	"log"
	"net/http"

	"demo/resp"

	"demo/pkg/session"

	"github.com/gin-gonic/gin"
)

// Login 登录
func Login(ctx *gin.Context) {
	user := models.NewUser()()
	// 此处模拟登录后已经拿到用户信息user
	log.Println("login user", user)

	sess := session.Default(ctx)

	/*set方法不会生成sessionID，
		但是会生成一个底层gorilla/session的一个session结构体对象并且存储到了store对象中，
		这里的store是"github.com/gin-contrib/sessions/redis"包的定义的结构体，
		它包含了"github.com/boj/redistore"包里面的对象RediStore。
		当然最终都是被包裹在了"github.com/gin-contrib/sessions"包里面的session结构体对象中
		type session struct {
		name    string
		request *http.Request
		store   Store
		session *sessions.Session
		written bool
		writer  http.ResponseWriter
	}

		于是乎，gin通过ctx传递这个session结构体对象。
		在sess.Get()调用时，其实是通过http.Request的header中的Cookie解密后取到之前存储的底层session结构体对象，并且赋值给最外层的store对象。
	*/

	sess.Set("userinfo", user)
	if err := sess.Save(); err != nil {
		resp.Fail(ctx, http.StatusInternalServerError, err)
		return
	}

	if err := sess.SaveUser(user.ID); err != nil {
		resp.Fail(ctx, http.StatusInternalServerError, err)
		return
	}
}
