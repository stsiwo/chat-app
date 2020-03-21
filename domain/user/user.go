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

