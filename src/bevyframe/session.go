package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
)

func (app Frame) getSessionToken(email string, token string) (string, error) {
	data, err := json.Marshal(map[string]string{"email": email, "token": token})
	if err != nil {
		return "", err
	}

	iv := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	block, err := aes.NewCipher(app.secret)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nil, iv, data, nil)
	tag := ciphertext[len(ciphertext)-aesGCM.Overhead():]
	ciphertext = ciphertext[:len(ciphertext)-aesGCM.Overhead()]

	return fmt.Sprintf("%s:%s:%s", hex.EncodeToString(iv), hex.EncodeToString(ciphertext), hex.EncodeToString(tag)), nil
}

func (app Frame) getSession(token string) (map[string]string, error) {
	parts := bytes.Split([]byte(token), []byte(":"))
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format")
	}

	iv, err := hex.DecodeString(string(parts[0]))
	if err != nil {
		return nil, err
	}

	ciphertext, err := hex.DecodeString(string(parts[1]))
	if err != nil {
		return nil, err
	}

	tag, err := hex.DecodeString(string(parts[2]))
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(app.secret)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(tag) != aesGCM.Overhead() {
		return nil, fmt.Errorf("invalid tag length")
	}

	plaintext, err := aesGCM.Open(nil, iv, append(ciphertext, tag...), nil)
	if err != nil {
		return nil, err
	}

	var result map[string]string
	if err := json.Unmarshal(plaintext, &result); err != nil {
		return nil, err
	}

	return result, nil
}
