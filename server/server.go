package main

import (
	"context"
	"fmt"
	proto "go-ws/protoc"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct{
	proto.UnimplementedExampleServer
}

func main(){
	fmt.Println("Hello from the server :)")
	listener,tcpErr:=net.Listen("tcp",":9000")
	if tcpErr!=nil{
		panic(tcpErr)
	}
	srv:=grpc.NewServer()
	proto.RegisterExampleServer(srv,&server{})
	reflection.Register(srv)
	if e:=srv.Serve(listener);e!=nil{
		fmt.Println(e)
		panic(e)
	}
}

func (s *server)HelloFunc(c context.Context,req *proto.HelloReq)(*proto.HelloRes,error){
	fmt.Println("request coming",req.Something)
	fmt.Println("let's reply to request coming from the client")
	return &proto.HelloRes{Reply: "Response from the server where we are handling the service (Grpc)"},nil
}

// Server side streaming example (* first need to update the .proto file :-) *)

// func (s *server)HelloFunc(req *proto.HelloReq,stream proto.ServerReplyServer) error{
// 	fmt.Println(req.Something)
// 	res:=[]*proto.HelloRes{
// 		{Reply: "Cool"},
// 		{Reply: "Cool"},
// 	}
// 	for _,msg:=range res {
// 		err:=stream.send(msg);
// 		if err!=nil{
// 			return err
// 		}
// 	}
// for {
	// msg,err:=stream.Recv()
	// if err==io.EOF{
	// 	break;
	// }
	// we can do some stuff with the provided msg and then return some response
	// fmt.println(msg)
	// stream.send("Coolness on its peak");
//	}
// 	return nil
// }
