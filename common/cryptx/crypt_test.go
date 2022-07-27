package cryptx

import (
	"encoding/base64"
	"fmt"
	"testing"
)

func Test_crypt(t *testing.T) {
	v, _ := AesGcmEncrypt([]byte("licat"), "233")
	fmt.Println(string(v))
	s := base64.StdEncoding.EncodeToString(v)
	fmt.Println(s)
	// v1, _ := AesGcmDecrypt(v, "233")
	// fmt.Println(string(v1))
}
