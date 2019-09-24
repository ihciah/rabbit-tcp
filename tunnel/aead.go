package tunnel

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"golang.org/x/crypto/chacha20poly1305"
	"strings"
)

var ErrCipherNotSupported = errors.New("cipher not supported")

const (
	aeadAes128Gcm        = "AEAD_AES_128_GCM"
	aeadAes192Gcm        = "AEAD_AES_192_GCM"
	aeadAes256Gcm        = "AEAD_AES_256_GCM"
	aeadChacha20Poly1305 = "AEAD_CHACHA20_POLY1305"
)

// List of AEAD ciphers: key size in bytes and constructor
var aeadList = map[string]struct {
	KeySize int
	New     func([]byte) (Cipher, error)
}{
	aeadAes128Gcm:        {16, aesGCM},
	aeadAes192Gcm:        {24, aesGCM},
	aeadAes256Gcm:        {32, aesGCM},
	aeadChacha20Poly1305: {32, chacha20Poly1305},
}

func makeAESGCM(key []byte) (cipher.AEAD, error) {
	blk, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(blk)
}

// aesGCM creates a new Cipher with a pre-shared key. len(psk) must be
// one of 16, 24, or 32 to select AES-128/196/256-GCM.
func aesGCM(psk []byte) (Cipher, error) {
	switch l := len(psk); l {
	case 16, 24, 32: // AES 128/196/256
	default:
		return nil, aes.KeySizeError(l)
	}
	return &metaCipher{psk: psk, makeAEAD: makeAESGCM}, nil
}

// chacha20Poly1305 creates a new Cipher with a pre-shared key. len(psk)
// must be 32.
func chacha20Poly1305(psk []byte) (Cipher, error) {
	if len(psk) != chacha20poly1305.KeySize {
		return nil, KeySizeError(chacha20poly1305.KeySize)
	}
	return &metaCipher{psk: psk, makeAEAD: chacha20poly1305.New}, nil
}

func NewAEADCipher(name string, key []byte, password string) (Cipher, error) {
	name = strings.ToUpper(name)
	switch name {
	case "CHACHA20-IETF-POLY1305":
		name = aeadChacha20Poly1305
	case "AES-128-GCM":
		name = aeadAes128Gcm
	case "AES-192-GCM":
		name = aeadAes192Gcm
	case "AES-256-GCM":
		name = aeadAes256Gcm
	}

	if choice, ok := aeadList[name]; ok {
		if key == nil || len(key) == 0 {
			key = kdf(password, choice.KeySize)
		}
		if len(key) != choice.KeySize {
			return nil, KeySizeError(choice.KeySize)
		}
		aead, err := choice.New(key)
		return aead, err
	}
	return nil, ErrCipherNotSupported
}
