package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"strings"
)

type data_pkt struct {
	Status bool
	Msg    string
}

func welcome_Msg() {
	fmt.Println("\n ---------------------------")
	fmt.Println("\n Welcome to Client Interface")
	fmt.Println("\n----------------------------")
}

func closing_Msg() {
	fmt.Println("\n ------------------")
	fmt.Println("\n We are closing !!!")
	fmt.Println("\n-------------------")
}

// Read data from terminal
func read_Terminal() string {
	buffer_io := bufio.NewReader(os.Stdin)
	input, err := buffer_io.ReadString('\n')
	checkError(err)
	return input
}

func main() {

	// Welcome Message
	welcome_Msg()

	// Ask for Server Details like IP Address / Website (with Port No.)
	fmt.Println("\nPlease provide Server Details\n")
	fmt.Println("eg :- www.google.com:80 or 74.125.200.105:80 :- ")

	webAddr := strings.TrimSpace(read_Terminal())

	// Resolve Server Address
	tcpAddr, err := net.ResolveTCPAddr("tcp", webAddr)
	checkError(err)

	// Create / Dial TCP connection
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)

	// Display options to Client for usage
	fmt.Println("\n\t--------------------------------------------------------")
	fmt.Println("\n\tType following commands for playing with Key-Value Store")
	fmt.Println("\n\t--------------------------------------------------------")
	fmt.Println("\n1. To store a key-value 			--> set <space> <key> <space> <value>")
	fmt.Println("\n2. To fetch a value (based on key) 	--> get <space> <key>")
	fmt.Println("\n3. To delete a key-value (based on key)	--> delete <space> <key>")
	fmt.Println("\n4. To close the program			--> close")
	fmt.Println("\n\n\n")

	// Ask for options while 
	for {
		fmt.Print("\nType command :- ")
		input := strings.ToLower(strings.TrimSpace(read_Terminal()))

		if input == "close" || input == "exit" {
			break
		}

		comm := strings.Split(input, " ")

		if comm[0] == "get" || comm[0] == "set" || comm[0] == "delete" {

			size := len(comm)

			if comm[0] == "get" {
				if size != 2 {
					fmt.Println("Invalid Command for 'get' \n Correct Command \t 'get <space> <key>' ")
					continue
				}
			} else if comm[0] == "set" {
				if size != 3 {
					fmt.Println("Invalid Command for 'set' \n Correct Command \t 'set <space> <key> <space> <value>' ")
					continue
				}
			} else if comm[0] == "delete" {
				if size != 2 {
					fmt.Println("Invalid Command for 'delete' \n Correct Command \t 'delete <space> <key>' ")
					continue
				}
			}

			// Send Request
			err = gob.NewEncoder(conn).Encode(input)
			checkError(err)

			var result data_pkt //string

			// Read Response
			err := gob.NewDecoder(conn).Decode(&result)
			checkError(err)

			// Print Response
			fmt.Println("\nResponse ", result)
		} else {
			// Print error and Ask for new user input
			fmt.Println("\nInvalid Command '", comm[0], "'")
		}

	}

	closing_Msg()

	os.Exit(0)

}

// Display Error (if any) and then close client program
func checkError(err error) {
	if err != nil {
		// panic(err)
		fmt.Fprintf(os.Stderr, "Fatal error: %s ", err.Error())
		fmt.Println("In CheckError")
		os.Exit(1)
	}
}
