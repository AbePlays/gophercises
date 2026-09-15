package secret

import (
	"errors"

	"github.com/AbePlays/gophercises/17-secrets/encrypt"
)

type MemoryVault struct {
	encodingKey string
	keyValues   map[string]string
}

func (v *MemoryVault) Get(key string) (string, error) {
	res, ok := v.keyValues[key]
	if !ok {
		return "", errors.New("key not found")
	}

	decryptedValue, err := encrypt.Decrypt(v.encodingKey, res)
	if err != nil {
		return "", err
	}

	return decryptedValue, nil
}

func (v *MemoryVault) Set(key, value string) error {
	encryptedValue, err := encrypt.Encrypt(v.encodingKey, value)
	if err != nil {
		return err
	}

	v.keyValues[key] = encryptedValue
	return nil
}

func Memory(encodingKey string) MemoryVault {
	return MemoryVault{
		encodingKey: encodingKey,
		keyValues:   make(map[string]string),
	}
}
