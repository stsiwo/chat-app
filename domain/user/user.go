package user

import (
	"github.com/google/uuid"
)

type User struct {

  id string

  name string

  role Role
}

func NewGuestUser(name string) *User {
  return &User{
    id: uuid.New().String(),
    name: name,
    role: Guest,
  }
}

func NewMemberUser(name string) *User {
  return &User{
    id: uuid.New().String(),
    name: name,
    role: Member,
  }
}

func NewAdminUser(name string) *User {
  return &User{
    id: uuid.New().String(),
    name: name,
    role: Admin,
  }
}

func (u *User) Id() string {
  return u.id
}

func (u *User) Name() string {
  return u.name
}

func (u *User) Role() Role {
  return u.role
}

