package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"
)

func main() {
	// создаём новый приватный RSA-ключ длиной 4096 бит
	// обратите внимание, что для генерации ключа и сертификата
	// используется rand.Reader в качестве источника случайных данных
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		log.Fatal(err)
	}

	//Saving public key to file
	file, err := os.OpenFile(`.\..\agent\public.key`, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	pem.Encode(file, &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(&privateKey.PublicKey),
	})
	file.Close()

	//Saving private key to file
	file, err = os.OpenFile(`.\..\server\private.key`, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	pem.Encode(file, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	file.Close()
}
