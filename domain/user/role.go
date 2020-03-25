package user

import ()

type Role int8

const (
  Guest Role = 0
  Member Role = 1
  Admin Role = 2
)
