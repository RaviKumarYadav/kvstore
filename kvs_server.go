package main

import "fmt"
import "net"
import "os"
import "strings"
import "encoding/gob"

const BUFF_SIZE int64 = 1024

//************************* KVStore MAP *********************************************
// Create a global "Map" (Key-Value Store) , so that it is available to all clients
// and the content in it resides untill Server is "ON" , it will help for clients to close their
// connection and reconnect again  to fetch stored data
//***********************************************************************************
var kvs map[string]string = make(map[string]string)

type data_pkt struct {
	Status bool
	Msg    string
}

func main() {

	// Server Port Number
	service := ":1201"

	// Resolve Server Address
	tcpAddr, err := net.ResolveTCPAddr("ip4", service)
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	// Display Message that "Server started successfully"
	fmt.Println("\nServer started Successfully!!!\nPort Number is '", service, "'")

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		// Go Routine to call method in a new Thread
		go handleClient(conn)

	}

}

func handleClient(conn net.Conn) {

	for {
		var err error
		var request string

		fmt.Println("We are Listening ... ")

		err = gob.NewDecoder(conn).Decode(&request)

		if err != nil {
			conn.Write([]byte("Error in 'reading' data at server."))
			return
		}

		// -------------------
		// Logic at Server
		// -------------------
		comm := strings.Split(request, " ")
		response := data_pkt{}

		if comm[0] == "set" {
			kvs[comm[1]] = comm[2]
			response.Status = true
			response.Msg = kvs[comm[1]] + " got added successfully."
			err = gob.NewEncoder(conn).Encode(response)
		} else if comm[0] == "get" {
			value, status := kvs[comm[1]]
			if status == true {
				response.Status = true
				response.Msg = comm[1] + " --> " + value
				err = gob.NewEncoder(conn).Encode(response)
			} else {
				response.Status = false
				response.Msg = "Error!!! \nNo key exists."
				err = gob.NewEncoder(conn).Encode(response)
			}
		} else if comm[0] == "delete" {
			temp, ok := kvs[comm[1]]

			if temp == "" && ok == false {
				response.Status = false
				response.Msg = comm[1] + " does not exists in K.V.Store ."
				err = gob.NewEncoder(conn).Encode(response)
			} else {
				delete(kvs, comm[1])
				response.Status = true
				response.Msg = comm[1] + " got deleted."
				err = gob.NewEncoder(conn).Encode(response)
			}
		}

		checkError(err)

		if err != nil {
			conn.Write([]byte("Error Occurred in Server somewhere."))
			return
		}

	}

}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error in Server : %s ", err.Error())
		os.Exit(1)
	}
}
