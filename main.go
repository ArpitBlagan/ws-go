package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type Info struct{
	Email string `json:"email"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var jwtKey = []byte("my_secret_key")

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

func handleLogin(w http.ResponseWriter ,r *http.Request){
	var info Info
	if err:=json.NewDecoder(r.Body).Decode(&info);err!=nil{
		fmt.Println("Error while decoding the request body")
		http.Error(w,"Error while decoding :(",http.StatusInternalServerError)
		return
	}
	// check in DB if the user exits let's say in our case it exits everytime.

	expirationTime := time.Now().Add(15 * time.Minute)
	claims := &Claims{
		Username: info.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  expirationTime,
		HttpOnly: true,
	})
}


func checkJwt(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	tokenStr := cookie.Value
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if !token.Valid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	w.Write([]byte(fmt.Sprintf("Welcome, %s!", claims.Username)))
}


func main(){
	fmt.Println("Hey from the Go server")
	// Here we can have our DB connection....

	router:=mux.NewRouter()
	// router.HandleFunc("/ws",handleWebSocket).Methods("GET")
	router.HandleFunc("/login",handleLogin).Methods("POST")
	router.HandleFunc("/check",checkJwt).Methods("GET")
	err:=http.ListenAndServe(":8000",router)
	if(err!=nil){
		fmt.Println("Error while creating a server :(",err)
		return
	}
	fmt.Println("Server is running on 9000 port :)")
}