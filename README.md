Key Value Store
===============

Store
-----
Key-Value store is simply implemented as Map[string]string , it is kept at server.


kvs_client.go
----------

kvs_client.go is the source file which contains implementation for Client. It connects to given server (using port number) and composes requests for adding , deleting , updating , fetching data from/to kvstore.

kvs_server.go
----------

kvs_server.go is the source file which contains implementation for Server. It serves each client request coming to it .

There are many different types of requets that can be exchanged between Client-Server.

Request Types :- 
-------------

  Set Request (Add new record)
>   * It will add given key with the given value
>   * Syntax "Set    key    value"

  Get Request (Fetch record)
>   * It will fetch value for the given key
>   * Syntax "Get    key"

  Delete Request (Delete record)
>   * It will delete record with the given key
>   * Syntax "Delete    key"


Test Cases
-----------

* Simple Set-Get Requests
    * First fire 1000 Set Requests with unique keys using 1000 different clients.
    * Then Fetch all the values using previous Keys.
* Updates Request
    * First Set 100 unique key-values.
    * Update value of few of the keys (eg 10 keys).
    * Fire 100 Get requests to cross-check.
* Deletes Request
    * First Set 100 unique key-values.
    * Delete value of few of the keys (eg 10 keys).
    * Fire 100 Get requests to cross-check that Deleted records are now present.
* Mixed Request
    * First Set 100 unique key-values.
    * Update value of few of the keys (eg 10 keys).
    * Delete value of few of the keys (eg 10 keys).
    * Fire 100 Get requests to cross-check.

How to Run Test Cases
-----------------------

```sh
cd github.com/RaviKumarYadav/kvstore
go test
```



License
----

[IIT Bombay]

[zmq]:http://zeromq.org/
[IIT Bombay]:http://www.cse.iitb.ac.in/
