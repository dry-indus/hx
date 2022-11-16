package util

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

func MD5V(str []byte) string {
	h := md5.New()
	h.Write(str)
	return hex.EncodeToString(h.Sum(nil))
}

func MD5O(o interface{}) string {
	s := fmt.Sprintf("%+v", o)
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
