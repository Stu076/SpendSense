package models

import (
	"time"

	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	bun.BaseModel `bun:"table:users" json:"-"`

	ID        int       `bun:",pk,autoincrement"`
	Username  string    `bun:",unique,notnull"`
	Email     string    `bun:",unique,notnull"`
	Password  string    `bun:",notnull"`
	Role      string    `bun:",notnull,default:'user'"`
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
}

func (u *User) HashPassword() error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashed)
	return nil
}

func (u *User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil
}
