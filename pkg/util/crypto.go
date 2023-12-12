package util

import (
	"crypto/aes"
	"crypto/cipher"
)

func AESDecrypt(encrypted []byte, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	size := block.BlockSize()
	mode := cipher.NewCBCDecrypter(block, iv[:size])

	decrypted := make([]byte, len(encrypted))
	mode.CryptBlocks(decrypted, encrypted)

	decrypted = pkcs5UnPadding(decrypted)

	return decrypted, nil
}

func pkcs5UnPadding(decrypted []byte) []byte {
	length := len(decrypted)
	unPadding := int(decrypted[length-1])
	return decrypted[:(length - unPadding)]
}
