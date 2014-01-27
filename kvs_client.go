package main

import (
"net"
"os"
"fmt"
"io/ioutil"
"bufio"
"strings"
)



func welcome_Msg(){
	fmt.Println("\n ---------------------------")
	fmt.Println("\n Welcome to Client Interface")
	fmt.Println("\n----------------------------")
}

func closing_Msg(){
	fmt.Println("\n ------------------")
	fmt.Println("\n We are closing !!!")
	fmt.Println("\n-------------------")
}


// Read data from terminal
func read_Terminal() string{
	buffer_io := bufio.NewReader(os.Stdin)
	input, err := buffer_io.ReadString('\n')
	checkError(err)
	return input
}


func main(){

// Welcome Message
welcome_Msg()


// Ask for Server Details like IP Address / Website (with Port No.)
fmt.Println("\nPlease provide Server Details\n")
fmt.Println("eg :- www.google.com:80 or 74.125.200.105:80 :- ")

webAddr := strings.TrimSpace(read_Terminal())


// Resolve Server Address
tcpAddr,err := net.ResolveTCPAddr("tcp4",webAddr)
checkError(err)


// Create / Dial TCP connection
conn,err := net.DialTCP("tcp",nil,tcpAddr)
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
for{
	fmt.Println("\nType command :- ")
	input := strings.ToLower(strings.TrimSpace(read_Terminal()))
	
	if input == "close"{
		break
	}
	
	comm := strings.Split(input, " ")
	fmt.Println("0 ",comm[0], " 1 ",comm[1]," 2 ",comm[2])
	

	if comm[0] == "get" || comm[0] == "set" || comm[0] == "delete" {
		// Send Request
		_, err = conn.Write([]byte(input))
		checkError(err)

		// Read Response
		result,err := ioutil.ReadAll(conn)
		checkError(err)

		// Printed Message
		fmt.Println(string(result))	
	}else{
		// Print error and Ask for new user input
		fmt.Println("\nInvalid Command '",input[0],"'")
	}

	
}

closing_Msg()

os.Exit(0)


//abc := strings.Split(input, " ")

//fmt.Println("\n",abc)

/*
        
        
        
	if len(os.Args) != 2{
		fmt.Fprintf(os.Stderr,"Usage: %s hostname\n",os.Args[0])
		fmt.Println("Usage: ",os.Args[0], "hostname")
		os.Exit(1)
	}

service := os.Args[1]

// Resolve Server Address
tcpAddr,err := net.ResolveTCPAddr("tcp4",service)
checkError(err)

// Create / Dial TCP connection
conn,err := net.DialTCP("tcp",nil,tcpAddr)
checkError(err)

// Sent Request
_, err = conn.Write([]byte("HEAD / HTTP/1.0\r\n\r\n"))
checkError(err)

// Read Response
result,err := ioutil.ReadAll(conn)
checkError(err)

// Printed Message
fmt.Println(string(result))

*/



}


// Display Error (if any) and then close client program
func checkError(err error){
	if err != nil{
		// panic(err)
		fmt.Fprintf(os.Stderr,"Fatal error: %s ",err.Error())
		os.Exit(1)
	}
}
