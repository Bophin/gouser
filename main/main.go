package main

import (
	"github.com/Bophin/gouser"
	"fmt"
	"os"
	"flag"
)
var command = "check"

func main() {
	username, password := initProgram()

	switch command {
		case "check": 
			if gouser.CheckPwd(username, password) {
				print("Success! logging in..\n")
			} else {
				print("Failure.. Invalid username or password.\n")
			}
		case "add": 
			if gouser.AddUser(username, password) {
				fmt.Printf("user added: %v\n", username)
			} else {
				fmt.Printf("user already exists: %v\n", username)
			}
		default: 
			fmt.Printf("Error, this should not be possible\n")
			os.Exit(1)
	}
}

//Parses command line arguments and flags, and returns username and password
//that will be used for the program.
func initProgram() (username, password string) {
	pwdflag := flag.String("f", "pwd.txt", "passwordfile to use")
	cmdflag := flag.Bool("a", false, "sets add mode")
	flag.Parse()
	args := flag.Args()
	
	gouser.SetPwdFile(*pwdflag)
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