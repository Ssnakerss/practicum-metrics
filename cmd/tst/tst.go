package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
)

var secretkey = []byte("secret key")

func generateRandom(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

const (
	password = "x35k9f"
	msg      = `0ba7cd8c624345451df4710b81d1a349ce401e61bc7eb704ca` +
		`a84a8cde9f9959699f75d0d1075d676f1fe2eb475cf81f62ef` +
		`f701fee6a433cfd289d231440cf549e40b6c13d8843197a95f` +
		`8639911b7ed39a3aec4dfa9d286095c705e1a825b10a9104c6` +
		`be55d1079e6c6167118ac91318fe`
)

func core(s string) string {
	fmt.Println("hello from core >", s)
	return s
}

func mw1(f func(string) string) func(string) string {
	ff := func(s string) string {
		fmt.Println("Hello from mw1 >", s)
		return f(s) + "+mw1"
	}
	return ff
}

func mw2(f func(string) string) func(string) string {
	ff := func(s string) string {
		fmt.Println("Hello from mw2 >", s)
		return f(s) + "+mw2"
	}
	return ff
}

func main() {

	fmt.Println(mw2(mw1(core))("hi"))

}

func decodeSample() {
	// допишите код
	// 1) получите ключ из password, используя sha256.Sum256
	key := sha256.Sum256([]byte(password))
	fmt.Println("1")
	// 2) создайте aesblock и aesgcm
	aesblock, err := aes.NewCipher(key[:])
	fmt.Println("2")
	if err != nil {
		log.Fatal(err)
	}
	aesgcm, err := cipher.NewGCM(aesblock)
	fmt.Println("3")
	if err != nil {
		log.Fatal(err)
	}
	// 3) получите вектор инициализации aesgcm.NonceSize() байт с конца ключа
	nonce := key[len(key)-aesgcm.NonceSize():]
	fmt.Println("4>", nonce)
	// 4) декодируйте сообщение msg в двоичный формат

	encrypted, err := hex.DecodeString(msg)

	fmt.Println("5>", encrypted)
	if err != nil {
		log.Fatal(err)
	}
	result, err := aesgcm.Open(nil, nonce, encrypted, nil)
	fmt.Println("6>", result)
	if err != nil {
		log.Fatal(err)
	}
	// 5) расшифруйте и выведите данные
	fmt.Printf("result: %s", string(result[:]))
}

func gcm() {
	src := []byte("Ключики от сердца твоего")
	fmt.Printf("original: %s\n", src)

	key, err := generateRandom(2 * aes.BlockSize)
	if err != nil {
		log.Fatal("error:", err.Error())
		return
	}
	//NewChiper сщздает и возвращает новый chiper.Block
	//Ключевым аргументом должен быть ключ AES 16 24 или 32 байта
	//для выбора AES-128, AES-192 или AES-256
	aesblock, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal("aes error", err.Error())
	}

	// NewGCM возвращает заданный 128-битный блочный шифр
	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		log.Fatal("aes error", err.Error())
	}
	//создаем вектор инициализации
	nonce, err := generateRandom(aesgcm.NonceSize())
	if err != nil {
		log.Fatal("aes error", err.Error())
	}
	//шифруем
	dst := aesgcm.Seal(nil, nonce, src, nil)
	fmt.Println(">>>>>")
	fmt.Println(">>>>>")
	fmt.Printf("%s", dst)
	fmt.Println(">>>>>")
	fmt.Println("<<<<<")
	//дешефруем
	src2, err := aesgcm.Open(nil, nonce, dst, nil)
	if err != nil {
		log.Fatal("aes error", err.Error())
	}

	fmt.Printf("decrypted: %s\n", src2)
}

func checkSign() {
	var (
		data []byte // декодированное сообщение с подписью
		id   uint32 // значение идентификатора
		err  error
		sign []byte // HMAC-подпись от идентификатора
	)
	msg := "048ff4ea240a9fdeac8f1422733e9f3b8b0291c969652225e25c5f0f9f8da654139c9e21"

	// допишите код
	// 1) декодируйте msg в data
	data, err = hex.DecodeString(msg)
	if err != nil {
		return
	}
	// 2) получите идентификатор из первых четырёх байт,
	//    используйте функцию binary.BigEndian.Uint32
	id = binary.BigEndian.Uint32(data[:4])
	// 3) вычислите HMAC-подпись sign для этих четырёх байт
	h := hmac.New(sha256.New, secretkey)
	h.Write(data[:4])
	sign = h.Sum(nil)
	// ...

	if hmac.Equal(sign, data[4:]) {
		fmt.Println("Подпись подлинная. ID:", id)
	} else {
		fmt.Println("Подпись неверна. Где-то ошибка")
	}
}
