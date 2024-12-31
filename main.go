package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

//Upgrader function
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool{
		return true
	},
}

func handleWebSocket(w http.ResponseWriter,r *http.Request){
	conn,err:=upgrader.Upgrade(w,r,nil)
	if(err!=nil){
		fmt.Println("Not able to upgrade http to ws",err)
		return
	}
	defer conn.Close()

	fmt.Println("New WebSocket connection established")

	for {
		// Read message from client
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			break
		}

		fmt.Printf("Received message: %s\n", string(message))

		// Echo message back to client
		err = conn.WriteMessage(messageType, message)
		if err != nil {
			fmt.Println("Error writing message:", err)
			break
		}
	}
}

func main(){
	fmt.Println("Hey from the Go server")
	// Here we can have our DB connection....
	
	router:=mux.NewRouter()
	router.HandleFunc("/ws",handleWebSocket).Methods("GET")
	err:=http.ListenAndServe(":9000",router)
	if(err!=nil){
		fmt.Println("Error while creating a server :(")
		return
	}
	fmt.Println("Server is running on 9000 port :)")
}