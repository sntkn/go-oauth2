package main

import (
	"encoding/base64"
	"fmt"
	"time"

	"crypto/ed25519"
	"crypto/rand"

	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/oauth2/internal/common/accesstoken"
)

func GenerateEd25519KeyPair() (string, string, error) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return "", "", err
	}
	publicKeyBase64 := base64.StdEncoding.EncodeToString(publicKey)
	privateKeyBase64 := base64.StdEncoding.EncodeToString(privateKey)

	// 結果を出力
	fmt.Println("Public Key (Base64):", publicKeyBase64)
	fmt.Println("Private Key (Base64):", privateKeyBase64)
	return publicKeyBase64, privateKeyBase64, nil
}

func main() {
	pub, pri, err := GenerateEd25519KeyPair()
	if err != nil {
		fmt.Println(err)
		return
	}

	// privateKeyBase64 := "PFsL4Ebcf4hW3EyNjxrnv0++NET1ZBLJMOxbez/xSXCDNO7cqWnldg2X2rPz7yyHTDFYdacUmC7wtr7nwg4/qQ=="

	t := accesstoken.TokenParams{
		UserID:    uuid.MustParse("00000000-0000-0000-0000-000000000000"),
		ClientID:  uuid.MustParse("00000000-0000-0000-0000-000000000000"),
		Scope:     "",
		ExpiresAt: time.Now(),
	}
	tt := &accesstoken.Token{}
	token, err := tt.Generate(&t, pri)
	if err != nil {
		panic(err)
	}

	// fmt.Println(token, err)

	// publicKey := "gzTu3Klp5XYNl9qz8+8sh0wxWHWnFJgu8La+58IOP6k="

	tt.Parse(token, pub)
	// fmt.Println(, err)
}
