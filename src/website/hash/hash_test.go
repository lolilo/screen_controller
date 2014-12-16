package hash

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"network"
)

var FILE_PATH_TO_USER_AUTHENTICATION_DATA = network.GetRelativeFilePath("./user_authentication_data_for_testing.txt")

func TestInitializeUserAuthenticationMap(t *testing.T) {
	result := InitializeUserAuthenticationMap(FILE_PATH_TO_USER_AUTHENTICATION_DATA)
	expectedResult := map[string][16]byte{"lilo": CreateHashFromString("poke")}
	assert.Equal(t, true, reflect.DeepEqual(result, expectedResult))
}

func TestGetUserNameAndPasswordFromFile(t *testing.T) {
	content := "lilo\npoke"
	username, password := GetUserNameAndPasswordFromFile(content)
	assert.Equal(t, username, "lilo")
	assert.Equal(t, password, "poke")
}

func TestCreateHashFromString(t *testing.T) {
	result := CreateHashFromString("yo")
	assert.NotNil(t, result)
}

func TestIsHashMatchInUserAuthenticationMap(t *testing.T) {
	hashMap := map[string][16]byte{
		"lilo": CreateHashFromString("secret"),
	}
	result := IsHashMatchInUserAuthenticationMap("lilo", "secret", hashMap)
	assert.True(t, result)
}
