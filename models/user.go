package models

import (
	"fmt"

	"demo/utils"

	"pkg.deepin.com/golang/lib/uuid"
)

// User ...
type User struct {
	ID    string
	Name  string
	Age   int
	Phone string
	Email string
}

// NewUser 模拟新增用户
func NewUser() func() *User {
	var i int
	return func() *User {
		i++
		return &User{
			ID:    uuid.UUID32(),
			Name:  fmt.Sprintf("用户%d", i),
			Age:   utils.RandRangeNum(18, 65),
			Phone: "1886666" + utils.IncrStr(i, 4),
			Email: fmt.Sprintf("%d@%d.com", i, i),
		}
	}
}
