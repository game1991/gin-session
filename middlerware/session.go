package middlerware

import (
	"demo/models"
	"demo/resp"
	"log"
	"net/http"

	mySession "demo/pkg/session"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
)

func Session() gin.HandlerFunc {
	store, err := redis.NewStoreWithDB(
		50,               //最大连接空闲数
		"tcp",            //网络方式
		"localhost:6379", //连接地址
		"",               //密码
		"",               //DB库号
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
}
