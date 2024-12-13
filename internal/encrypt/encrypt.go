package encrypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
)

type Coder struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func (c *Coder) GenerateKeys(pathToPublicKey string, pathToPrivateKey string) error {
	whoAmI := "encrypt.GenerateKeys" // имя функции, которая вызывается

	// создаём новый приватный RSA-ключ длиной 4096 бит
	// обратите внимание, что для генерации ключа и сертификата
	// используется rand.Reader в качестве источника случайных данных
	var err error
	c.privateKey, err = rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return fmt.Errorf("%s: %f", whoAmI, err)
	}
	c.publicKey = &c.privateKey.PublicKey
	if pathToPrivateKey != "" && pathToPublicKey != "" {
		//Saving public key to file
		file, err := os.OpenFile(pathToPublicKey, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			return fmt.Errorf("%s: %f", whoAmI, err)
		}
		pem.Encode(file, &pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: x509.MarshalPKCS1PublicKey(c.publicKey),
		})
		file.Close()

		//Saving private key to file
		file, err = os.OpenFile(pathToPrivateKey, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			return fmt.Errorf("%s: %f", whoAmI, err)
		}
		pem.Encode(file, &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(c.privateKey),
		})
		file.Close()
	}
	return nil
}

func (c *Coder) LoadPrivateKey(pathToPrivateKey string) error {
	whoAmI := "encrypt.LoadPrivateKey" // имя функции, которая вызывается
	if pathToPrivateKey != "" {
		//Loading private key from file
		keyData, err := os.ReadFile(pathToPrivateKey)
		if err != nil {
			return fmt.Errorf("%s: %f", whoAmI, err)
		}
		pemBlock, _ := pem.Decode(keyData) // считываем pem блок
		if pemBlock == nil || pemBlock.Type != "RSA PRIVATE KEY" {
			return fmt.Errorf("%s: %f", whoAmI, errors.New("failed to decode PEM block or not a private key"))
		}

		c.privateKey, err = x509.ParsePKCS1PrivateKey(pemBlock.Bytes)
		if err != nil {
			return fmt.Errorf("%s: %f", whoAmI, err)
		}
		return nil
	} else {
		return fmt.Errorf("%s: %f", whoAmI, errors.New("path private key not defined"))
	}
}

func (c *Coder) LoadPublicKey(pathToPublicKey string) error {
	whoAmI := "encrypt.LoadPublicKey" // имя функции, которая вызывается
	if pathToPublicKey != "" {
		//Loading public key from file
		keyData, err := os.ReadFile(pathToPublicKey)
		if err != nil {
			return fmt.Errorf("%s: %f", whoAmI, err)
		}
		pemBlock, _ := pem.Decode(keyData) // считываем pem блок
		if pemBlock == nil || pemBlock.Type != "RSA PUBLIC KEY" {
			return fmt.Errorf("%s: %f", whoAmI, errors.New("failed to decode PEM block or not a public key"))
		}

		c.publicKey, err = x509.ParsePKCS1PublicKey(pemBlock.Bytes)
		if err != nil {
			return fmt.Errorf("%s: %f", whoAmI, err)
		}

		return nil
	} else {
		return fmt.Errorf("%s: %f", whoAmI, errors.New("path public key not defined"))
	}
}

func (c *Coder) Encrypt(data []byte) ([]byte, error) {
	whoAmI := "encrypt.Encrypt" // имя функции, которая вызывается

	if c.publicKey != nil {
		eData, err := rsa.EncryptOAEP(
			sha256.New(),
			rand.Reader,
			c.publicKey,
			data,
			nil)

		if err != nil {
			return nil, fmt.Errorf("%s: %f", whoAmI, err)
		}
		return eData, nil
	}
	return nil, fmt.Errorf("%s: %f", whoAmI, errors.New("public key not set"))
}

func (c *Coder) Decrypt(data []byte) ([]byte, error) {
	whoAmI := "encrypt.Decrypt" // имя функции, которая вызывается

	if c.publicKey != nil {
		eData, err := rsa.DecryptOAEP(
			sha256.New(),
			nil,
			c.privateKey,
			data,
			nil)

		if err != nil {
			return nil, fmt.Errorf("%s: %f", whoAmI, err)
		}
		return eData, nil
	}
	return nil, fmt.Errorf("%s: %f", whoAmI, errors.New("private key not set"))
}

//----------------------------------------------
