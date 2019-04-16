package main

import (
	"crypto/rand"
	b32 "encoding/base32"
	"time"
)

func generateUniqueTokenString() (string, error) {
	t := time.Now()
	ms := uint64(t.Unix())*1000 + uint64(t.Nanosecond()/int(time.Millisecond))

	token := make([]byte, 16)
	token[0] = byte(ms >> 40)
	token[1] = byte(ms >> 32)
	token[2] = byte(ms >> 24)
	token[3] = byte(ms >> 16)
	token[4] = byte(ms >> 8)
	token[5] = byte(ms)

	_, err := rand.Read(token[6:])
	if err != nil {
		return "", err
	}

	tokenStr := b32.StdEncoding.WithPadding(b32.NoPadding).EncodeToString(token)
	return tokenStr, nil
}

func getExpirationTime(days int) int64 {
	return time.Now().AddDate(0, 0, days).Unix()
}

func reachedExpirationTime(t int64) bool {
	return time.Now().After(time.Unix(t, 0))
}
