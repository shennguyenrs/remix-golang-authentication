package models

import (
	"time"

	"github.com/uptrace/bun"
)

// Using pointer and "omitempty" in json tag for optinal tag
// Because if the value is missing in the unmarshal step
// the value will record as empty
// by using * and "omitempty" there no field in the json object
type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`
	ID            int       `bun:",pk,autoincrement"   json:"id"`
	Name          string    `bun:",unique,notnull"     json:"name"       validate:"required,alphanumunicode"`
	Email         string    `bun:",unique,notnull"     json:"email"      validate:"required,email"`
	Password      string    `bun:",notnull"            json:"password"   validate:"required,alphanumunicode"`
	LastLogin     time.Time `bun:""                    json:"last_login"                                     vaidate:"reuired,datetime"`
}

type LoginForm struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,alphanumunicode"`
}

type RegisterForm struct {
	Name     string `json:"name"     validate:"required,alphanumunicode"`
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,alphanumunicode"`
}
