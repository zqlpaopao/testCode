package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)
/**
借鉴 "github.com/wenzhenxi/gorsa"
https://blog.csdn.net/u011142688/article/details/79380129

******/

func main() {
	var publicStr = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCk9GDYL9+9l8Jly9OFFbOC4uaQ
PD2XIBZXIQy+tPsuBmSg19HpsTi609hhxkqIKHoou3722Ks9dgUxA1A5g8uOwr+j
wPdXhJiNFK5sjuV1EdHlykiQwI+gSHku2R4lmaThPSWJ8T9YiN7ILoRaj86EnHDQ
DhBqg4JD8q+/dkEyQwIDAQAB
-----END PUBLIC KEY-----`


	block, _ := pem.Decode([]byte(publicStr))
	if block == nil {

	}
	// x509 parse public key
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		fmt.Println(err)
	}

	//fmt.Println( pub.(*rsa.PublicKey))
	//
	var p = pub.(*rsa.PublicKey)
	//fmt.Println(sha256.New().Size() )
	//fmt.Println(p.Size())
	encryptedBytes, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		p,
		[]byte("hellow"),
		nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("encrypted bytes: ", string(encryptedBytes))



	//解密
	var privatekey =`-----BEGIN PRIVATE KEY-----
MIICdgIBADANBgkqhkiG9w0BAQEFAASCAmAwggJcAgEAAoGBAKT0YNgv372XwmXL
04UVs4Li5pA8PZcgFlchDL60+y4GZKDX0emxOLrT2GHGSogoeii7fvbYqz12BTED
UDmDy47Cv6PA91eEmI0UrmyO5XUR0eXKSJDAj6BIeS7ZHiWZpOE9JYnxP1iI3sgu
hFqPzoSccNAOEGqDgkPyr792QTJDAgMBAAECgYBOCeNrQ7LpQkvQ1w45zxt/F5OW
tzk4LxECpXsfGgYfLx0aTyBbG+HH2YNsNmB6bBPnA1U8uSLCT/yCxJuGqkh5fgsr
mHWfySG1PiGZ0PD0pwZVsmw0jrfQgOFF95r3mZJ1204OGeLzlZKAjB7JAumtBA16
QLgHDDYYnaczyyz/EQJBANh7QNM16JMwr3Adk80Gu5Zlbg2dqbIbVgkIjximFm0W
PfPe4Uvp/bxo7Irx9A+kCFmgVF4mPQN5S1bUBNaiN4kCQQDDESPPHw3vyIhTzLu0
pQWG3ACA5d0y1JHpF8ZLllOmpAcJqrgA+5NJczkbuJpr5fSBi4MbvLOaoPiDBghm
GRxrAkAU4E3wEFLNXvSMK04Fh5CvgDiMt5eVxW0Wkey6w8mF24895VB0sav2b2fg
PlT67Sag/gUkzyszGo9ZYDjXOe2BAkEAn2c7Pv9ekSrrFKfCYB1WRd00YCD3QJlq
3vL5rT0cAJob0i97C/qJYsVQzrFtJ20UAGS0cA8lKeAPFGrypBQzHwJATSfnG1OR
K4ASKy13HCBTlRoyCL4/7uI6eAjcfywAQLr+dPhE2Nmo8xn4ga5qHQvk3qfZWwWa
Qq6c3cCDlKKY0A==
-----END PRIVATE KEY-----`
	blocks, _ := pem.Decode([]byte(privatekey))
	if block == nil {

	}
	pri, err := x509.ParsePKCS1PrivateKey(blocks.Bytes)
	if err == nil {

	}
	fmt.Println(pri)
	pri2, err := x509.ParsePKCS8PrivateKey(blocks.Bytes)
	if err != nil {

	}
	 fmt.Println(pri2.(*rsa.PrivateKey))

	encryptedBytess, err := rsa.DecryptOAEP(

		sha256.New(),
		rand.Reader,
		pri2.(*rsa.PrivateKey),
		[]byte(encryptedBytes),
		nil)
	if err != nil {
		panic(err)
	}

fmt.Println(string(encryptedBytess))
}