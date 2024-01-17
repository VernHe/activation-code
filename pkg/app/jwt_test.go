package app

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestToken(t *testing.T) {
	token, err := CreateToken(UserInfo{
		UserId:   "testId",
		Username: "testName",
		Roles:    []string{"testRole"},
	})
	assert.NoError(t, err)
	err = VerifyToken(token)
	assert.NoError(t, err)

	time.Sleep(6 * time.Second)
	err = VerifyToken(token)
	assert.Error(t, err)
}
