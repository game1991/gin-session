package controller

import (
	"demo/resp"
	"fmt"
	"net/http"
	"strings"

	mySession "demo/pkg/session"

	"github.com/gin-contrib/sessions"

	"github.com/gin-gonic/gin"
)

// Logout ...
func Logout(ctx *gin.Context) {
	sess := mySession.Default(ctx)
	openid := strings.TrimSpace(ctx.Query("openid"))
	if openid != "" {
		// 如果存在传参openid，则需要通过对应的sessionid进行退出
		/*
			由源码中可以知道sess.Get()方式是通过http.Request中Header携带的Cookie信息进行查询底层redis的
			所以不能通过此方法去解决sessionid查询以及对象删除。需要直接去底层数据库redis删除对应sessionid即可
		*/
		ses, err := sess.GetByOpenID(openid)
		if err != nil {
			resp.Fail(ctx, http.StatusInternalServerError, err)
			return
		}

		if ses == nil {
			resp.Fail(ctx, http.StatusBadRequest, fmt.Errorf("openID：%v传参得到:%v", openid, ses))
			return
		}

		if err := sess.DeleteBySessionID(ses); err != nil {
			resp.Fail(ctx, http.StatusInternalServerError, err)
			return
		}

		ctx.Redirect(http.StatusTemporaryRedirect, "http://www.baidu.com")
		return
	}

	sess.Options(sessions.Options{MaxAge: -1})
	sess.Clear()
	if err := sess.Save(); err != nil {
		resp.Fail(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Redirect(http.StatusTemporaryRedirect, "http://www.baidu.com")
	return
}
