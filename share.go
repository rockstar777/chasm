package main

import (
	"crypto/rand"
	"encoding/base64"
)

// ShareID is a uniqiue id to represent uploaded shares
type ShareID string

// Share represents a secret share of a file
type Share struct {
	SID  ShareID
	Data []byte
}

// CreateShares creates two shares from secret
// for now just 2-out-of-n
func CreateShares(secret []byte, sid ShareID, n int) []Share {
	if n > 2 {
		panic("n > 2. unsupported for now.")
	}

	secretLength := len(secret)

	firstShareBytes := make([]byte, secretLength)
	_, err := rand.Read(firstShareBytes)
	if err != nil {
		panic(err)
	}

	secondShareBytes := xor(secret, firstShareBytes)

	return []Share{
		Share{SID: sid, Data: firstShareBytes},
		Share{SID: sid, Data: secondShareBytes},
	}
}

// CombineShares restores the secret by adding
func CombineShares(shares []Share) []byte {
	if len(shares) > 2 {
		panic("n > 2. unsupported for now.")
	}

	return xor(shares[0].Data, shares[1].Data)
}

/// Helper Functions ///

func xor(a, b []byte) []byte {
	if len(a) != len(b) {
		panic("xor must take equal arrays")
	}
	n := len(a)

	res := make([]byte, n)

	for i := 0; i < len(a); i++ {
		res[i] = a[i] ^ b[i]
	}

	return res
}

// RandomShareID randomly generates a 16 byte base64URL encoded string
func RandomShareID() ShareID {
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}

	// base64 URL Encode the bytes:
	return ShareID(base64.URLEncoding.EncodeToString(randomBytes))
}
