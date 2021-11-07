package middlerware

import (
	"demo/models"
	"demo/resp"
	"encoding/gob"
	"log"
	"net/http"

	mySession "demo/pkg/session"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
)

func init() {
	gob.Register(&models.User{})
}

// Session 中间件
func Session() gin.HandlerFunc {
	store, err := redis.NewStoreWithDB(
		50,               //最大连接空闲数
		"tcp",            //网络方式
		"localhost:6379", //连接地址
		"",               //密码
		"0",              //DB库号
		[]byte("secret"), // 密匙
	)
	if err != nil {
		log.Fatal("redis NewStoreWithDB err:", err)
	}

	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   20 * 60,
		HttpOnly: true,
	})

	err, rediStore := redis.GetRedisStore(store)
	if err != nil {
		log.Fatal("Session中间件 GetRedisStore err:", err)
	}

	sessionSecureCodec := securecookie.CodecsFromPairs([]byte("secret"))
	myStore := &mySession.RedisStore{Pairs: sessionSecureCodec, Store: store, RedisStore: rediStore}
	return mySession.Sessions("mysession", myStore)
}

// SessionCheck 登录检查
func SessionCheck(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	u := sess.Get("user")
	if u == nil {
		resp.Fail(ctx, http.StatusUnauthorized, nil)
		return
	}
	user, ok := u.(*models.User)
	if !ok {
		resp.Fail(ctx, http.StatusUnauthorized, nil)
		return
	}
	ctx.Set("user", user)
	// 如果需要将session时长续期，可以调用 sess.Flashs()和sess.Save()，或者调用sess.Options(&sessions.Options{MaxAge = ?})和sess.Save()
	// 需要注意的是如果需要续期并且用户id与sessionid有关联，请调用sess.SaveUser()对用户id关联绑定续期。
}
