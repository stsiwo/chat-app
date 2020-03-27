package user

import (
	"github.com/google/uuid"
  "testing"
  "encoding/json"
  "log"
	"github.com/stretchr/testify/assert"
  "github.com/stsiwo/chat-app/domain/user"
)

func TestUserJsonEncode(t *testing.T) {

  dummyUser := user.NewAdminUser(uuid.New().String(), "sample-admin")

  jsoned, err := json.Marshal(dummyUser)
  if err != nil {
    log.Fatalf("err during encoding instance to json byte: %v", err)
  }

  var newDummyUser *user.User
  json.Unmarshal(jsoned, &newDummyUser)

  assert.Equal(t, dummyUser.Id(), newDummyUser.Id())
  assert.Equal(t, dummyUser.Name(), newDummyUser.Name())
  assert.Equal(t, dummyUser.Role(), newDummyUser.Role())
}
