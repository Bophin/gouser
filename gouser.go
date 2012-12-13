package main

import (
	"crypto/sha256"
	"os"
	"io"
	"fmt"
	"strings"
	"flag"
)

type userData struct {
	id string
	salt string
	hash string
}

var pwdfile string = "pwd.txt"

func main() {
	flag.Parse()
	args := flag.Args()

	//f, _ := os.OpenFile(pwdfile, os.O_CREATE, 0644)
	//f.Close()
	
	
	if len(args) != 2 {
		fmt.Printf("Incorrect usage, write 'gouser [user] [pwd]'\n")
		os.Exit(0)
	}

	username := args[0]
	password := args[1]

	if addUser(username, password) {
		fmt.Printf("user added: %v\n", username)
	} else {
		fmt.Printf("user already exists: %v\n", username)
		checkPwd(username, password)
	}
}

func addUser(username, password string) bool {

	if userExists(username) {
		return false
	}
	
	salt := genSalt()
	print(salt + "salt from generation\n")
	hash := getHash(salt, password)

	result := fmt.Sprintf("%v:%v:%v\n", username, salt, hash)
	fmt.Printf("%v", result)

	flags := os.O_WRONLY | os.O_APPEND | os.O_CREATE
	file, err := os.OpenFile(pwdfile, flags, 0644)
	if err != nil {
		panic(err)
	}

	file.Write([]byte(result))
	err = file.Close()
	if err != nil {
		panic(err)
	}
	return true
}

func genSalt() string {
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
	return fmt.Sprintf("%x", salt)
}

func userExists(username string) bool {
	users := getUsers()

	for _, val := range users {
		if val.id == username {
			return true
		}
	}
	return false
}

func readFile(name string) (string, error) {
	file, err := os.OpenFile(pwdfile, os.O_RDONLY | os.O_CREATE, 0644)
	if err != nil {
		return "", err
	}
	
	b := make([]byte, 1000)
	for {
		n, err := file.Read(b)
		if err == io.EOF {
			return string(b), nil
		} else if err != nil {
			return "", err
		}

		if n < len(b) {
			break
		} else {
			b2 := make([]byte, cap(b)*2, cap(b)*2)
			for i := range b {
				b2[i] = b[i]
			}
			b = b2
		}
	}
	file.Close()
	return string(b), nil
}

func getUsers() []userData {
	//Splits per row
	f_str, err := readFile(pwdfile)
	if err != nil {
		panic(err)
	}
	rows := strings.Split(f_str, "\n")

	users := make([]userData, len(rows))
	for i, val := range rows {
		row := strings.Split(val, ":")
		if len(row) == 3 {
			users[i] = userData{row[0], row[1], row[2]}
		}
	}
	return users
}

func checkPwd(username, password string) bool {
	users := getUsers()

	for _, val := range users {
		//Splits rows by :
		if val.id == username {
			print(val.salt + " = salt from file\n")
			hash := getHash(val.salt, password)
			if hash == val.hash {
				print("success! logging in..\n")
			} else {
				fmt.Printf("%v !=\n%v\n", hash, val.hash)
			}
		}
	}
	return false
}

func getHash(salt, pwd string) string {
	hash := sha256.New()
	io.WriteString(hash, salt + pwd)
	return fmt.Sprintf("%x", hash.Sum(nil))
}
