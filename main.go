package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"time"
)

func main() {
	// generate temporary identity
	pubKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}

	// build and serialize test message
	m := Message{
		Author:    pubKey,
		Timestamp: time.Now(),
		Content: MessageContent{
			Type: MessageTypePost,
			Post: &Post{
				Text: "Hello, world!",
			},
		},
	}
	signed, err := m.Sign(privateKey)
	if err != nil {
		panic(err)
	}

	serialized, err := json.Marshal(signed)
	if err != nil {
		panic(err)
	}

	// validate message on receiver end
	var deserialized SignedMessage
	err = json.Unmarshal(serialized, &deserialized)
	if err != nil {
		panic(err)
	}

	valid, err := deserialized.Verify()
	if err != nil {
		panic(err)
	}

	fmt.Println(valid)
}
