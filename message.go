package main

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/json"
	"time"
)

// SignedMessage ...
type SignedMessage struct {
	ID        []byte `json:"id"`        // sha256 hash of Message
	Message   []byte `json:"message"`   // json encoded message
	Signature []byte `json:"signature"` // signed Message by message.Author
}

// Message ...
type Message struct {
	Author    ed25519.PublicKey `json:"author"`
	Timestamp time.Time         `json:"time"`
	Content   MessageContent    `json:"content"`
}

// MessageType ...
type MessageType string

const (
	// MessageTypePost ...
	MessageTypePost = "post"
)

// MessageContent ...
type MessageContent struct {
	Type MessageType `json:"type"`

	// horrible b/c of go's lack of sum types
	*Post `json:"post,omitempty"`
}

// Post ...
type Post struct {
	Text string `json:"text"`
}

// Sign ...
func (m Message) Sign(privateKey ed25519.PrivateKey) (*SignedMessage, error) {
	serializedMessage, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	h := sha256.New()
	_, err = h.Write(serializedMessage)
	if err != nil {
		return nil, err
	}

	return &SignedMessage{
		ID:        h.Sum(nil),
		Message:   serializedMessage,
		Signature: ed25519.Sign(privateKey, serializedMessage),
	}, nil
}

// Verify ...
func (sm SignedMessage) Verify() (bool, error) {
	var m Message
	err := json.Unmarshal(sm.Message, &m)
	if err != nil {
		return false, err
	}

	return ed25519.Verify(m.Author, sm.Message, sm.Signature), nil
}
