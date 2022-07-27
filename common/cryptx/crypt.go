package cryptx

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/scrypt"
)

const (
	SECRET = "licat2333333"
)

func aesKey(plaintext string) (ciphertext []byte, err error) {
	ciphertext, err = scrypt.Key([]byte(plaintext), []byte(SECRET), 1<<15, 8, 1, 32)
	return
}

// 加密密码
func HashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		return ""
	}
	return string(hash)
}

func AesGcmEncrypt(unencrypted []byte, password string) ([]byte, error) {
	key, err := aesKey(password)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// generate a random nonce (makes encryption stronger)
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	encrypted := gcm.Seal(nil, nonce, unencrypted, nil)
	// we need nonce for decryption so we put it at the beginning
	// of encrypted text
	return append(nonce, encrypted...), nil
}

func AesGcmDecrypt(encrypted []byte, password string) ([]byte, error) {
	key, err := aesKey(password)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(encrypted) < gcm.NonceSize() {
		return nil, errors.New("Invalid data")
	}

	// extract random nonce we added to the beginning of the file
	nonce := encrypted[:gcm.NonceSize()]
	encrypted = encrypted[gcm.NonceSize():]

	return gcm.Open(nil, nonce, encrypted, nil)
}

/** MD5V
 * description: md5加密，不可逆
 * param: str []byte, b ...byte
 * return: string
 */
func MD5V(str []byte, b ...byte) string {
	h := md5.New()
	h.Write(str)
	return hex.EncodeToString(h.Sum(b))
}

func ToMD5V(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
