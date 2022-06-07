package crypto

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"hash"
)

// Sha1Stream Sha1Stream
type Sha1Stream struct {
	_sha1 hash.Hash
}

// Update Update Sha1Stream
func (obj *Sha1Stream) Update(data []byte) {
	if obj._sha1 == nil {
		obj._sha1 = sha1.New()
	}
	obj._sha1.Write(data)
}

// Sum Sum
func (obj *Sha1Stream) Sum() string {
	return hex.EncodeToString(obj._sha1.Sum([]byte("")))
}

// Sha1 Sha1
func Sha1(data string) string {
	h := sha1.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum([]byte("")))
}

// MD5 MD5
func MD5(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum([]byte("")))
}
