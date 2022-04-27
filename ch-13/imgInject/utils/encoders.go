package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

func encodeDecode(input []byte, key string) []byte {
	var bArr = make([]byte, len(input))
	for i := 0; i < len(input); i++ {
		bArr[i] += input[i] ^ key[i%len(key)]
	}
	return bArr
}

//XorEncode returns encoded byte array
func XorEncode(decode []byte, key string) []byte {
	return encodeDecode(decode, key)
}

//XorDecode returns decoded byte array
func XorDecode(encode []byte, key string) []byte {
	return encodeDecode(encode, key)
}

//AESencrypt returns encoded byte array
func AESencrypt(decode []byte, key string) []byte {
	keyPass := []byte(createHash(key))
	c, err := aes.NewCipher(keyPass)
	if err != nil {
		fmt.Println(err)
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil{
		fmt.Println(err)
	}
	nonce := make([]byte, gcm.NonceSize())
	return gcm.Seal(nonce, nonce, decode, nil)
}

//AESdecrypt returns decoded byte array
func AESdecrypt(encode []byte, key string) []byte {
	keyPass := []byte(createHash(key))
	block, err := aes.NewCipher(keyPass)
	if err != nil {
		fmt.Println(err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Println(err)
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := encode[:nonceSize], encode[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		fmt.Println(err)
	}
	return plaintext


	// gcm, err := cipher.NewGCM(c)
	// if err != nil{
	// 	fmt.Println(err)
	// }
	// nonceSize := gcm.NonceSize()
	// if len(encode) < nonceSize {
	// 	fmt.Println(err)
	// }
	// fmt.Print("GOT HERE\n")
	// //nonce, ciphertext := encode[:nonceSize], encode[nonceSize:]
	// nonce := make([]byte, gcm.NonceSize())
	// fmt.Print("NOW HERE\n")
	// plaintext, err := gcm.Open(nil, nonce, encode, nil)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// return plaintext
}

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}
