package security

import "testing"

var (
	privateKey, publicKey []byte
)
var RSA = &RsaSecurity{}

func tearup() {
	privateKey, publicKey, _ = GenerateRASPrivateAndPublicKeys()
	RSA, _ = NewRsaSecurityFromStringKey(string(publicKey), string(privateKey))

}

// 公钥加密私钥解密
func Test_PubENCTYPTPriDECRYPT(t *testing.T) {
	tearup()
	pubenctypt, err := RSA.PublicKeyEncrypt([]byte(`hello world`))
	if err != nil {
		t.Error(err)
	}

	pridecrypt, err := RSA.PrivateKeyDecrypt(pubenctypt)
	if err != nil {
		t.Error(err)
	}
	if string(pridecrypt) != `hello world` {
		t.Error(`不符合预期`)
	}
}

// 公钥解密私钥加密
func Test_PriENCTYPTPubDECRYPT(t *testing.T) {
	tearup()
	prienctypt, err := RSA.PrivateKeyEncrypt([]byte(`hello world`))
	if err != nil {
		t.Error(err)
	}
	pubdecrypt, err := RSA.PublicKeyDecrypt(prienctypt)
	if err != nil {
		t.Error(err)
	}
	if string(pubdecrypt) != `hello world` {
		t.Error(`不符合预期`)
	}
}

// 公钥解密私钥加密
func Test_priEPubD(t *testing.T) {
	data := "hello word"
	tearup()
	enData, err := PrivateKeyEncrypt([]byte(data), RSA.privateKey)
	if err != nil {
		t.Error(err)
	}
	deData, err := PublicKeyDecrypt(enData, RSA.publicKey)
	if err != nil {
		t.Error(err)
	}
	if string(deData) != data {
		t.Error(`不符合预期`)
	}
}
