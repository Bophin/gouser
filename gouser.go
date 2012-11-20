package main

import (
	"crypto/sha256"
	"os"
	"io"
	"fmt"
)

func main() {
	addUser("test", "password")
}

func addUser(username, password string) {
	
	hash := sha256.New()
	salt := make([]byte, 32)
	file, err := os.Open("/dev/random")
	if err != nil {
		panic(err)
	}

	_, err = file.Read(salt)
	if err != nil {
		panic(err)
	}
	err = file.Close()
	if err != nil {
		panic(err)
	}

	io.WriteString(hash, string(salt) + password)

	result := fmt.Sprintf("%v:%x:%x\n", username, salt, hash.Sum(nil))
	fmt.Printf("%v", result)

	flags := os.O_WRONLY | os.O_APPEND | os.O_CREATE
	file, err = os.OpenFile("pwd.txt", flags, 0644)
	if err != nil {
		panic(err)
	}

	file.Write([]byte(result))
	err = file.Close()
	if err != nil {
		panic(err)
	}

}
