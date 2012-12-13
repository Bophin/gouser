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

var pwdfile = "pwd.txt"
var command = "check"

func main() {
	username, password := initProgram()

	switch command {
		case "check": 
			if CheckPwd(username, password) {
				print("Success! logging in..\n")
			} else {
				print("Failure.. Invalid username or password.\n")
			}
		case "add": 
			if AddUser(username, password) {
				fmt.Printf("user added: %v\n", username)
			} else {
				fmt.Printf("user already exists: %v\n", username)
			}
		default: 
			fmt.Printf("Error, this should not be possible\n")
			os.Exit(1)
	}
}

func SetPwdFile(name string) {
	pwdfile = name
}

//Parses command line arguments and flags, and returns username and password
//that will be used for the program.
func initProgram() (username, password string) {
	pwdflag := flag.String("f", "pwd.txt", "passwordfile to use")
	cmdflag := flag.Bool("a", false, "sets add mode")
	flag.Parse()
	args := flag.Args()
	
	pwdfile = *pwdflag
	if *cmdflag {
		command = "add"
	}

	if len(args) != 2 {
		fmt.Printf("Incorrect usage, write 'gouser [user] [pwd]'\n")
		os.Exit(0)
	}

	username = args[0]
	password = args[1]
	return
}

//Tries to add user to the passwordfile. Returns false if the username
//already exists in the file. Otherwise it generates password salt and hash
//which are stored together with the username in the file. 
func AddUser(username, password string) bool {

	if userExists(username) {
		return false
	}
	
	salt := genSalt()
	hash := getHash(salt, password)

	result := fmt.Sprintf("%v:%v:%v\n", username, salt, hash)

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

//Checks the passwordfile for the given username.
func userExists(username string) bool {
	users := getUsers()

	for _, val := range users {
		if val.id == username {
			return true
		}
	}
	return false
}

//Reads the passwordfile and parses the information from it. The result is returned
//as a list of userData, with username, salt and hashed password.
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

//Reads 32 random bytes from /dev/random and returns them as a string
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

//Compares the username and password with passwordfile to see if
//it finds a match.
func CheckPwd(username, password string) bool {
	users := getUsers()

	for _, val := range users {
		//Splits rows by :
		if val.id == username {
			hash := getHash(val.salt, password)
			if hash == val.hash {
				return true
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

//Reads the file given by name and returns the whole file as a string
func readFile(name string) (string, error) {
	file, err := os.OpenFile(name, os.O_RDONLY | os.O_CREATE, 0644)
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