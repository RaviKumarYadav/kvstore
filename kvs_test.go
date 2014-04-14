package kvs_testing

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"testing"
)

type data_pkt struct {
	Status bool
	Msg    string
}

var kvs map[string]string

/*********************************************************************
	Starts a Server for each test case 
*********************************************************************/
func startServer(portNo int) {

	service := ":" + strconv.Itoa(portNo)

	// Resolve Server Address
	tcpAddr, err := net.ResolveTCPAddr("ip4", service)
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		// Go Routine to call method in a new Thread
		go handleClient(conn)
	}
}

/*********************************************************************
	Logic for Server
*********************************************************************/
func handleClient(conn net.Conn) {

	for {
		var err error
		var request string

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
				response.Msg = value
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

func startClient(portNo int, msg string) data_pkt {

	webAddr := strings.TrimSpace("localhost:" + strconv.Itoa(portNo))

	// Resolve Server Address
	tcpAddr, err := net.ResolveTCPAddr("tcp", webAddr)
	checkError(err)

	// Create / Dial TCP connection
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)

	input := strings.ToLower(strings.TrimSpace(msg))

	comm := strings.Split(input, " ")

	var result data_pkt

	if comm[0] == "get" || comm[0] == "set" || comm[0] == "delete" {

		// Send Request
		err = gob.NewEncoder(conn).Encode(input)
		checkError(err)

		// Read Response
		err := gob.NewDecoder(conn).Decode(&result)
		checkError(err)

	}

	return result

}

/*********************************************************************
	Scenario : Testing for With 1000 Clients , All will send "Set"
	Outcome  : Server should be appropriately updated
*********************************************************************/
func TestWith_1000_Clients(t *testing.T) {

	fmt.Println("\n\n******************************************************")
	fmt.Println("Testing with 1000 Clients using simple 'Set' and 'Get' ")
	fmt.Println("******************************************************")

	kvs = make(map[string]string)
	serverPortNo := 1201

	go startServer(serverPortNo)

	// Sending Data
	fmt.Println("Sending 1000 'Set' instructions...")
	for i := 0; i < 1000; i++ {
		startClient(serverPortNo, "set "+strconv.Itoa(i)+" "+strconv.Itoa(i+100))
	}

	// Receiving Data
	fmt.Println("Receiving Same Data by 'Get'...")
	for i := 0; i < 1000; i++ {
		pkt_rcvd := startClient(serverPortNo, "get "+strconv.Itoa(i))

		if pkt_rcvd.Status == false {
			t.Errorf("\n'Get' failed for ", i)
		}
	}
}

/****************************************************************************************
	Scenario : Sending 100 "Set" followed by 100 "Set" to update existing key-value
	Outcome  : Server should be appropriately updated
****************************************************************************************/
func TestWith_Updates(t *testing.T) {

	fmt.Println("\n\n***********************************************")
	fmt.Println("Testing with Clients updating existing values ")
	fmt.Println("***********************************************")

	serverPortNo := 1201
	clear_kvs(serverPortNo, 1000)

	// Sending Data
	fmt.Println("Sending 1000 'Set' instructions...")
	for i := 0; i < 100; i++ {
		startClient(serverPortNo, "set "+strconv.Itoa(i)+" "+strconv.Itoa(i+100))
	}

	for i := 0; i < 100; i++ {
		startClient(serverPortNo, "set "+strconv.Itoa(i)+" "+strconv.Itoa(i+1000))
	}

	// Receiving Data
	fmt.Println("Receiving Same Data by 'Get'...")
	for i := 0; i < 100; i++ {
		pkt_rcvd := startClient(serverPortNo, "get "+strconv.Itoa(i))

		value, _ := strconv.Atoi(pkt_rcvd.Msg)

		if pkt_rcvd.Status == false || value != (i+1000) {
			t.Errorf("\nGet failed for ", i)
		}
	}
}

/****************************************************************************************
	Scenario : Testing with few clients deleting existing values
	Outcome  : Server should not respond with deleted values
****************************************************************************************/
func TestWith_Deletes(t *testing.T) {

	fmt.Println("\n\n****************************************************")
	fmt.Println("Testing with few clients deleting existing values ")
	fmt.Println("****************************************************")

	serverPortNo := 1201
	clear_kvs(serverPortNo, 100)

	// Sending Data
	fmt.Println("Sending 1000 'Set' instructions...")
	for i := 0; i < 100; i++ {
		startClient(serverPortNo, "set "+strconv.Itoa(i)+" "+strconv.Itoa(i+100))
	}

	for i := 50; i < 60; i++ {
		startClient(serverPortNo, "delete "+strconv.Itoa(i))
	}

	// Receiving Data
	fmt.Println("Receiving Same Data by 'Get'...")
	for i := 50; i < 60; i++ {
		pkt_rcvd := startClient(serverPortNo, "get "+strconv.Itoa(i))

		if (i >= 50 && i < 60) && pkt_rcvd.Status == true {
			t.Errorf("\nFetched Deleted Values for ", i)
		} else if i >= 60 && (pkt_rcvd.Status == false || pkt_rcvd.Msg != strconv.Itoa(i+100)) {
			t.Errorf("\nFetched Incorrect Values for ", i)
		}
	}

}

/****************************************************************************************
	Scenario : Mixed Testing with insertion , deletion , modification of few values
	Outcome  : Server should not respond with appropriate values
****************************************************************************************/
func TestWith_Mixed_Ops(t *testing.T) {

	fmt.Println("\n\n**************************************************************************")
	fmt.Println("Mixed Testing with insertion , deletion , modification of few values ")
	fmt.Println("**************************************************************************")

	serverPortNo := 1201
	clear_kvs(serverPortNo, 100)

	// Sending Data
	fmt.Println("Sending 1000 'Set' instructions...")
	for i := 0; i < 100; i++ {
		startClient(serverPortNo, "set "+strconv.Itoa(i)+" "+strconv.Itoa(i+100))
	}

	for i := 50; i < 70; i++ {
		startClient(serverPortNo, "delete "+strconv.Itoa(i))
	}

	for i := 80; i < 100; i++ {
		startClient(serverPortNo, "set "+strconv.Itoa(i)+" "+strconv.Itoa(i+1000))
	}

	// Receiving Data
	fmt.Println("Receiving Same Data by 'Get'...")
	for i := 0; i < 100; i++ {
		pkt_rcvd := startClient(serverPortNo, "get "+strconv.Itoa(i))

		if (i >= 50 && i < 70) && pkt_rcvd.Status == true {
			t.Errorf("\nFetched Deleted Values for ", i)
		} else if (i >= 80 && i < 100) && (pkt_rcvd.Status == false || pkt_rcvd.Msg != strconv.Itoa(i+1000)) {
			t.Errorf("\nFetched Incorrect Values for ", i)
		}
	}

}

func clear_kvs(serverPortNo int, upper_limit int) {

	for i := 0; i < upper_limit; i++ {
		startClient(serverPortNo, "delete "+strconv.Itoa(i))
	}

}

// Display Error (if any) and then close client program
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s ", err.Error())
		fmt.Println("In CheckError")
		os.Exit(1)
	}
}
