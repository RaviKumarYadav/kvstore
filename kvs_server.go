package main

import "fmt"
import "net"
import "os"
import "strings"
import "encoding/gob"


const BUFF_SIZE int64 = 1024

// Create a global "Map" (Key-Value Store) , so that it is available to all clients
// and the content in it resides untill Server is "ON" , it will help for clients to close their
// connection and reconnect again  to fetch stored data
var kvs map[string]string = make(map[string]string)

func main(){

service := ":1201"

// Resolve Server Address
tcpAddr, err := net.ResolveTCPAddr("ip4",service)
checkError(err)

listener, err := net.ListenTCP("tcp",tcpAddr)
checkError(err)

// Display Message that "Server started successfully"
fmt.Println("\nServer started Successfully!!!\nPort Number is '",service,"'")

// Run Server Program forever
for{
	conn, err := listener.Accept()
	if err != nil{
		continue
	}
	
	// Go Routine to call method in a new Thread
	go handleClient(conn)
	
}

// No Exit in Server
// os.Exit(0)

}


func handleClient(conn net.Conn){

	
	// close connection on exit from method
	defer conn.Close()
	

		
	for{
	
		//var buf[] byte = make([]byte, 1024)//BUFF_SIZE)
		
		var err error
		var request string
		
		// Read upto BUFF_SIZE bytes
		//n,err := conn.Read(buf[0:])
		err = gob.NewDecoder(conn).Decode(&request)
	
		if err != nil{
			conn.Write([]byte("Error in 'reading' data at server."))
			return
		}
		
		// -------------------
		// Logic at Server
		// -------------------
			
		// Display the Request Received
		fmt.Println("\nRequest received :- ",request,"\n")
		
		comm := strings.Split(request, " ")
		
		
		if comm[0] == "set"{
			kvs[comm[1]] = comm[2]
			//_, err = conn.Write([]byte(kvs[comm[1]] + " got added successfully."))
			err = gob.NewEncoder(conn).Encode(kvs[comm[1]] + " got added successfully.")
		}else if comm[0] == "get"{
			value,status := kvs[comm[1]]
			if status == true {
				err = gob.NewEncoder(conn).Encode(comm[1] + " --> " + value)
				//_, err := conn.Write([]byte(value))
			}else{
				err = gob.NewEncoder(conn).Encode("Error!!! \nNo key exists.")
				//_, err := conn.Write([]byte("Error!!! \nNo key exists."))
			}
		}else if comm[0] == "delete"{
			delete(kvs,comm[1])
			err = gob.NewEncoder(conn).Encode(comm[1] + " got deleted.")
			//_, err := conn.Write([]byte(comm[1] + " got deleted."))
		}
		
		checkError(err)
				
		if err != nil{
			conn.Write([]byte("Error Occurred in Server somewhere."))
			return
		}
		
	}

}

func checkError(err error){
	if err != nil{
		fmt.Fprintf(os.Stderr,"Fatal error in Server : %s ",err.Error())
		os.Exit(1)
	}
}
