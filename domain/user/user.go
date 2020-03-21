package user

import ()

type User struct {

  id string

  name string

  role Role
}

func NewGuestUser(id string, name string) *User {
  return &User{
    id: id,
    name: name,
    role: Guest,
  }
}

func NewMemberUser(id string, name string) *User {
  return &User{
    id: id,
    name: name,
    role: Member,
  }
}

func NewAdminUser(id string, name string) *User {
  return &User{
    id: id,
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

