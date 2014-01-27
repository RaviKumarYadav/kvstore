package main

import "fmt"
import "net"
import "os"

/*
	Wait for Request and Respond by sending  the same 
	content received at Server and then close the 
	connection with that client.
*/

const BUFF_SIZE int32 = 512

// Create a global "Map" (Key-Value Store) , so that it is available to all clients
// and the content in it resides untill Server is "ON" , it will help for clients to close their
// connection and reconnect again  to fetch stored data
var kvs map[string]string = make(map[string]string)

func main(){

service := ":1201"

// Resolve Server Address
tcpAddr, err := net.ResolveTCPAddr("tcp4",service)
checkError(err)

listener, err := net.ListenTCP("tcp",tcpAddr)
checkError(err)

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
	
	var buf[BUFF_SIZE] byte
		
	for{
		// Read upto BUFF_SIZE bytes
		n,err := conn.Read(buf[0:])
	
		if err != nil{
			conn.Write(buf[]("Error in 'reading' data at server."))
			return
		}
		
		// -------------------
		// Logic at Server end
		// -------------------
			
		request := buf[0:n]
		comm := strings.Split(request, " ")
		
		var err error
		
		if comm[0] == "set"{
			kvs[comm[1]] = comm[2]
			_, err = conn.Write(comm[1] , " got added successfully.")
		}else if comm[0] == "get"{
			value,status := kvs[comm[1]]
			if status == true {
				_, err = conn.Write(value)
			}else{
				_, err = conn.Write("Error in 'get'.")
				return
			}
		}else if comm[0] == "delete"{
			delete(kvs,comm[1])
			_, err = conn.Write(comm[1]," got deleted.")
		}
		
		//fmt.Println(string(buf[0:]))
		//_, err2 := conn.Write(buf[0:n])
		
		if err != nil{
			conn.Write("Error Occurred somewhere.")
			return
		}
	}

}

func checkError(err error){
	if err != nil{
		fmt.Fprintf(os.Stderr,"Fatal error: %s ",err.Error())
		os.Exit(1)
	}
}
