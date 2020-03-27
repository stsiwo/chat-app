package user

import (
  "encoding/json"
  "strconv"
  "log"
  "strings"
)

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

func (u *User) MarshalJSON() ([]byte, error) {
  return json.Marshal(map[string]interface{}{
    "id": u.id,
    "name": u.name,
    "role": u.role,
  })
}

func (u *User) UnmarshalJSON(rawData []byte) error {
  log.Println("unmarshal json is called inside user struct")
  var objMap map[string]json.RawMessage
  err := json.Unmarshal(rawData, &objMap)
  if err != nil {
    log.Fatalf("err during decoding user json data to map struct: %v", err)
    return err
  }
  // need to remove double quote from rawMessage
  u.id = strings.Trim(string(objMap["id"]), "\"")
  // need to remove double quote from rawMessage
  u.name = strings.Trim(string(objMap["name"]), "\"")
  roleNum, err := strconv.Atoi(string(objMap["role"]))
  if err != nil {
    log.Fatalf("err during converting user role number string to int: %v", err)
    return err
  }
  u.role = Role(roleNum)
  return nil
}


