package msg

import (
  "github.com/stsiwo/chat-app/domain/user"
  "time"
)

type Message struct {

  id string

  // must be Member or Guest
  sender user.User

  // must be Admin or nil (not assign yet; message from new Member/Guest)
  reciever user.User

  // any type => 'interface{}' and use type switch for handling each type
  content interface{}

  date time.Time

}

