package keyring

import (
	"fmt"

	"github.com/zalando/go-keyring"
)

const (
	serviceName = "spin"
)

type Credential struct {
	Username   string
	SessionKey string
}

func SetCredential(username string, sessionKey string) error {
	return keyring.Set(serviceName, username, sessionKey)
}

func GetCredential(username string) (Credential, error) {
	item, err := keyring.Get(serviceName, username)
	if err != nil {
		return Credential{}, fmt.Errorf("credential not found for user %s: %w", username, err)
	}

	return Credential{
		Username:   username,
		SessionKey: item,
	}, nil
}

func DeleteCredential(username string) error {
	return keyring.Delete(serviceName, username)
}
