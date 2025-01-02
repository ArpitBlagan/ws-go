package main

import (
	"context"
	"encoding/json"
	"fmt"
	proto "go-ws/protoc"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)
var client proto.ExampleClient
func main(){
	fmt.Println("Runing client :)")
	router:=mux.NewRouter()
	conn,err:=grpc.Dial("localhost:9000",grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err!=nil{
		fmt.Println(err)
		panic(err)
	}
	client=proto.NewExampleClient(conn)
	router.HandleFunc("/get",handleGet).Methods("GET")
	errS:=http.ListenAndServe(":8000",router)
	if errS!=nil{
		fmt.Println("Error while creating a server",errS)
		return;
	}
	defer conn.Close()
}

func handleGet(w http.ResponseWriter,r *http.Request){
	req:=&proto.HelloReq{Something:"hello from client"}
	res,err:=client.HelloFunc(context.TODO(),req)
	if err != nil {
        log.Fatalf("Error calling HelloFunc: %v", err)
		http.Error(w,"Error while making a grpc call",http.StatusInternalServerError)
		return;
    }
	fmt.Println(res)
	ress,_:=json.Marshal(res)
	w.Header().Set("Content-type","application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(ress)
}
//*** For Client Side Streaming....
// for client side stream we will have the stream as a parameter right
// so the code will look like this:
//req:=[]*proto.HelloReq{{Something:""},{Somthing:""}}
//stream,err:=client.ServerReply(context.TODO())
//for _,re :=range(req){
// err:=stream.Send(re)
// if err!=nil{
// fmt.Println(err)
//return;}
//}
//res,err:=stream.CloseAndRecv()