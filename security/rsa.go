package security

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

type RsaSecurity struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

func NewRsaSecurityFromRsaKey(publicKey *rsa.PublicKey, privateKey *rsa.PrivateKey) (result *RsaSecurity) {
	result = &RsaSecurity{}
	result.publicKey = publicKey
	result.privateKey = privateKey
	return
}
func NewRsaSecurityFromStringKey(publicKey, privateKey string) (result *RsaSecurity, err error) {
	result = &RsaSecurity{}
	if len(publicKey) > 0 {
		block, _ := pem.Decode([]byte(publicKey))
		if block == nil {
			return nil, errors.New("get public key error")
		}
		// x509 parse public key
		pub, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		result.publicKey = pub.(*rsa.PublicKey)
	}
	if len(privateKey) > 0 {
		block, _ := pem.Decode([]byte(privateKey))
		if block == nil {
			return nil, errors.New("get private key error")
		}
		result.privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			pri2, err := x509.ParsePKCS8PrivateKey(block.Bytes)
			if err != nil {
				return nil, err
			}
			result.privateKey = pri2.(*rsa.PrivateKey)
		}
	}
	return
}

func GenerateRASPrivateAndPublicKeys() (privateKey, publicKey []byte, err error) {
	pri, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}
	derTx := x509.MarshalPKCS1PrivateKey(pri)
	block := pem.Block{Type: "RSA PRIVATE KEY", Bytes: derTx}
	privateKey = pem.EncodeToMemory(&block)
	stream, err := x509.MarshalPKIXPublicKey(&pri.PublicKey)
	if err != nil {
		return nil, nil, err
	}
	block = pem.Block{Type: "RSA PUBLIC KEY", Bytes: stream}
	publicKey = pem.EncodeToMemory(&block)
	return privateKey, publicKey, nil
}

func (s *RsaSecurity) PublicKeyEncrypt(input []byte) (encryptedBlockBytes []byte, err error) {
	msgLen := len(input)
	h := sha256.New()
	rng := rand.Reader
	label := []byte("efucloud-encrypt")
	step := s.publicKey.Size() - 2*h.Size() - 2
	var encryptedBytes []byte
	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}
		encryptedBlockBytes, err := rsa.EncryptOAEP(h, rng, s.publicKey, input[start:finish], label)
		if err != nil {
			return nil, err
		}

		encryptedBytes = append(encryptedBytes, encryptedBlockBytes...)
	}

	return encryptedBytes, nil
}

func (s *RsaSecurity) PrivateKeyDecrypt(input []byte) (decryptedBytes []byte, err error) {
	msgLen := len(input)
	step := s.privateKey.PublicKey.Size()
	h := sha256.New()
	rng := rand.Reader
	label := []byte("efucloud-encrypt")
	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}
		decryptedBlockBytes, err := rsa.DecryptOAEP(h, rng, s.privateKey, input[start:finish], label)
		if err != nil {
			return nil, err
		}
		decryptedBytes = append(decryptedBytes, decryptedBlockBytes...)
	}

	return decryptedBytes, nil
}
