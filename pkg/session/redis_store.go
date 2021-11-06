package session

import (
	"net/http"

	"github.com/boj/redistore"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gorilla/securecookie"
	gs "github.com/gorilla/sessions"
)

// RedisStore Store
type RedisStore struct {
	Pairs []securecookie.Codec
	redis.Store
	RedisStore *redistore.RediStore
}

// Save Save
func (s *RedisStore) Save(r *http.Request, w http.ResponseWriter, sess *gs.Session) error {
	err := s.Store.Save(r, w, sess)
	if err != nil {
		return err
	}

	oldMaxAge := sess.Options.MaxAge
	if oldMaxAge != -1 {
		//为了浏览器退出cookie时效 加的logic
		sess.Options.MaxAge = 0
		encoded, _ := securecookie.EncodeMulti(sess.Name(), sess.ID, s.Pairs...)
		w.Header().Set("Set-Cookie", gs.NewCookie(sess.Name(), encoded, sess.Options).String())
		//还原现场
		sess.Options.MaxAge = oldMaxAge
	}
	return nil
}

// Get Get
func (s *RedisStore) Get(r *http.Request, name string) (*gs.Session, error) {

	return gs.GetRegistry(r).Get(s, name)
}
