package secret

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/AbePlays/gophercises/17-secrets/encrypt"
)

type FileVault struct {
	encodingKey string
	filePath    string
	mutex       sync.Mutex
	keyValues   map[string]string
}

func (v *FileVault) Get(key string) (string, error) {
	v.mutex.Lock()
	defer v.mutex.Unlock()

	if err := v.loadKeyValues(); err != nil {
		return "", err
	}

	res, ok := v.keyValues[key]
	if !ok {
		return "", errors.New("key not found")
	}

	return res, nil
}

func (v *FileVault) Set(key, value string) error {
	v.mutex.Lock()
	defer v.mutex.Unlock()

	err := v.loadKeyValues()
	if err != nil {
		return err
	}

	v.keyValues[key] = value
	err = v.saveKeyValues()
	return err
}

func File(encodingKey, filePath string) FileVault {
	return FileVault{
		encodingKey: encodingKey,
		filePath:    filePath,
		keyValues:   make(map[string]string),
	}
}

func (v *FileVault) loadKeyValues() error {
	file, err := os.Open(v.filePath)
	if err != nil {
		v.keyValues = make(map[string]string)
		return nil
	}
	defer file.Close()

	var sb strings.Builder
	_, err = io.Copy(&sb, file)
	if err != nil {
		return err
	}

	decryptedJson, err := encrypt.Decrypt(v.encodingKey, sb.String())
	if err != nil {
		return err
	}

	reader := strings.NewReader(decryptedJson)
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&v.keyValues); err != nil {
		return err
	}

	return nil
}

func (v *FileVault) saveKeyValues() error {
	var sb strings.Builder
	encoder := json.NewEncoder(&sb)
	if err := encoder.Encode(v.keyValues); err != nil {
		return err
	}

	encryptedJson, err := encrypt.Encrypt(v.encodingKey, sb.String())
	if err != nil {
		return err
	}

	file, err := os.OpenFile(v.filePath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprint(file, encryptedJson)
	return err
}
