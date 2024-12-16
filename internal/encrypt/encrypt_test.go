package encrypt

import "testing"

func TestCoder_GenerateKeys(t *testing.T) {
	e := Coder{}

	err := e.GenerateKeys("public_key.pem", "private_key.pem")
	if err != nil {
		t.Errorf("Coder.GenerateKeys() error = %v, wantErr %v", err, false)
	}

	err = e.LoadPrivateKey("private_key.pem")
	if err != nil {
		t.Errorf("Coder.LoadPrivateKey() error = %v, wantErr %v", err, false)
	}

	err = e.LoadPublicKey("public_key.pem")
	if err != nil {
		t.Errorf("Coder.LoadPublicKey() error = %v, wantErr %v", err, false)
	}

	data := "Hello, world!"
	encodedData, err := e.Encrypt([]byte(data))
	if err != nil {
		t.Errorf("Coder.Encrypt() error = %v, wantErr %v", err, false)
	}
	decodedData, err := e.Decrypt(encodedData)
	if err != nil {
		t.Errorf("Coder.Decrypt() error = %v, wantErr %v", err, false)
	}
	if string(decodedData) != data {
		t.Errorf("Coder.Decrypt() error = %v, wantErr %v", err, false)
	}

}
