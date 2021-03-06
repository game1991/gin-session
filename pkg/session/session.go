package session

import (
	"bytes"
	"demo/utils"
	"encoding/gob"
	"errors"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/context"
	gs "github.com/gorilla/sessions"
)

const (
	errorFormat = "[sessions] ERROR! %s\n"
)

// Store ...
type Store interface {
	gs.Store
	Options(sessions.Options)
}

// Sessioner Wraps thinly gorilla-session methods.
// Sessioner stores the values and optional configuration for a session.
type Sessioner interface {
	// ID of the session, generated by stores. It should not be used for user data.
	ID() string
	// Get returns the session value associated to the given key.
	Get(key interface{}) interface{}
	// Set sets the session value associated to the given key.
	Set(key interface{}, val interface{})
	// Delete removes the session value associated to the given key.
	Delete(key interface{})
	// Clear deletes all values in the session.
	Clear()
	// AddFlash adds a flash message to the session.
	// A single variadic argument is accepted, and it is optional: it defines the flash key.
	// If not defined "_flash" is used by default.
	AddFlash(value interface{}, vars ...string)
	// Flashes returns a slice of flash messages from the session.
	// A single variadic argument is accepted, and it is optional: it defines the flash key.
	// If not defined "_flash" is used by default.
	Flashes(vars ...string) []interface{}
	// Options sets configuration for a session.
	Options(sessions.Options)
	// Save saves all sessions used during the current request.
	Save() error

	SaveUser(openid string) error //绑定用户id和sessionid到数据库DB

	GetByOpenID(openID string) (*gs.Session, error) //通过用户ID查询session对象

	DeleteBySessionID(session *gs.Session) error //通过sessionID进行删除session对象
}

// Session 重构的 session
type Session struct {
	KeyPrefix string
	//RediStore *redistore.RediStore
	name    string
	request *http.Request
	store   Store
	session *gs.Session
	written bool
	writer  http.ResponseWriter
}

// SaveUser 绑定用户id和sessionid到DB
func (s *Session) SaveUser(openID string) error {
	defer utils.Def()
	store, ok := s.store.(*RedisStore)
	if !ok {
		return errors.New("SaveUser store类型错误，应为RedisStore")
	}

	// 用户id信息绑定到session并且存到redis
	conn := store.RedisStore.Pool.Get()
	defer conn.Close()
	if err := conn.Err(); err != nil {
		return err
	}

	age := s.Session().Options.MaxAge
	if age == 0 {
		age = store.RedisStore.DefaultMaxAge
	}

	if _, err := conn.Do("SETEX", openID, age, s.Session().ID); err != nil {
		return err
	}
	return nil
}

// GetByOpenID 通过openID获取session对象
func (s *Session) GetByOpenID(openID string) (*gs.Session, error) {
	defer utils.Def()
	store, ok := s.store.(*RedisStore)
	if !ok {
		return nil, errors.New("GetByOpenID store类型错误，应为RedisStore")
	}

	conn := store.RedisStore.Pool.Get()
	defer conn.Close()
	resp, err := conn.Do("GET", openID)
	if err != nil {
		return nil, err
	}
	if resp == nil || string(resp.([]byte)) == "" {
		return nil, nil // no data was associated with this key
	}

	sessionID := string(resp.([]byte))

	// 通过sessionID查询redis中存储的session对象的Values字段map内容
	data, err := conn.Do("GET", s.KeyPrefix+sessionID)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil // no data was associated with this key
	}
	b, err := redis.Bytes(data, err)
	if err != nil {
		return nil, err
	}
	ss := gs.NewSession(s.store, s.name)
	ss.ID = sessionID
	dec := gob.NewDecoder(bytes.NewBuffer(b))
	if err := dec.Decode(&ss.Values); err != nil {
		return nil, err
	}
	return ss, nil

}

// DeleteBySessionID 通过sessionID进行删除session对象
func (s *Session) DeleteBySessionID(session *gs.Session) error {
	defer utils.Def()
	if session == nil || session.ID == "" {
		return errors.New("Delete(session *gs.Session)方法调用 session对象 入参为空")
	}
	store, ok := s.store.(*RedisStore)
	if !ok {
		return errors.New("GetByOpenID store类型错误，应为RedisStore")
	}

	conn := store.RedisStore.Pool.Get()
	defer conn.Close()

	if _, err := conn.Do("DEL", s.KeyPrefix+session.ID); err != nil {
		return err
	}
	return nil
}

// Sessions ...重构
func Sessions(name string, store Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		s := &Session{"session_", name, c.Request, store, nil, false, c.Writer}
		c.Set(sessions.DefaultKey, s)
		defer context.Clear(c.Request)
		c.Next()
	}
}

// SessionsMany 重构
func SessionsMany(names []string, store sessions.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := make(map[string]*Session, len(names))
		for _, name := range names {
			session[name] = &Session{"session_", name, c.Request, store, nil, false, c.Writer}
		}
		c.Set(sessions.DefaultKey, session)
		defer context.Clear(c.Request)
		c.Next()
	}
}

// ID 获取gs.Session的ID
func (s *Session) ID() string {
	return s.Session().ID
}

// Get 获取gs.Session中的Values字段的map中key对应的数据
func (s *Session) Get(key interface{}) interface{} {
	return s.Session().Values[key]
}

// Set  设置gs.Session中的Values字段的map中key对应的数据
func (s *Session) Set(key interface{}, val interface{}) {
	s.Session().Values[key] = val
	s.written = true
}

// Delete 删除gs.Session中的Values字段的map中key对应的数据
func (s *Session) Delete(key interface{}) {
	delete(s.Session().Values, key)
	s.written = true
}

// Clear 清空gs.Session对象中的字段Values的map
func (s *Session) Clear() {
	for key := range s.Session().Values {
		s.Delete(key)
	}
}

// AddFlash 添加一个flash message to session
func (s *Session) AddFlash(value interface{}, vars ...string) {
	s.Session().AddFlash(value, vars...)
	s.written = true
}

// Flashes 返回全部的flash message in session
func (s *Session) Flashes(vars ...string) []interface{} {
	s.written = true
	return s.Session().Flashes(vars...)
}

// Options 将store传入的options传给底层session对象
func (s *Session) Options(options sessions.Options) {
	s.Session().Options = options.ToGorillaOptions()
}

// Save 更新gs.Session结构体对象的Values这个字段的map内容并且会更新session在redis数据库中的存活时间，maxage传入的options，当sessionID不存在时，会生成一个新的sessionID
func (s *Session) Save() error {
	if s.Written() {
		e := s.Session().Save(s.request, s.writer)
		if e == nil {
			s.written = false
		}
		return e
	}
	return nil
}

// Session ...
func (s *Session) Session() *gs.Session {
	if s.session == nil {
		var err error
		s.session, err = s.store.Get(s.request, s.name) //通过http.Request的Header中Cookie信息获取sessionID，如果没有sessionID就新生成一个底层的gs.Session对象
		if err != nil {
			log.Printf(errorFormat, err)
		}
	}
	return s.session
}

// Written ...
func (s *Session) Written() bool {
	return s.written
}

// Default shortcut to get session
func Default(c *gin.Context) Sessioner {
	return c.MustGet(sessions.DefaultKey).(Sessioner)
}

// DefaultMany shortcut to get session with given name
func DefaultMany(c *gin.Context, name string) Sessioner {
	return c.MustGet(sessions.DefaultKey).(map[string]Sessioner)[name]
}
