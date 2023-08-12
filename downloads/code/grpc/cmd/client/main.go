package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/Guaderxx/grpc-demo/pkg/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial(":8082", opts...)
	if err != nil {
		panic("grcp Dial err")
	}
	defer conn.Close()

	client := protos.NewAddressBookServiceClient(conn)
	ctx := context.Background()
	// HandleAddPerson(ctx, client)
	// HandleListPeople(ctx, client)
	HandleAddPeople(ctx, client)
}

func HandleAddPerson(ctx context.Context, cli protos.AddressBookServiceClient) {
	per := protos.Person{Name: "aader"}
	status, err := cli.AddPerson(ctx, &per)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("Status: ", status.Ok)
}

func HandleListPeople(ctx context.Context, cli protos.AddressBookServiceClient) {
	res, err := cli.ListPeople(ctx, &protos.ListReq{Listnum: 10})
	if err != nil {
		log.Println(err)
		return
	}
	for {
		tmp, err := res.Recv()
		if err == io.EOF {
			break
		}
		fmt.Println("Tmp: ", tmp.GetPeople())
	}
	fmt.Println("ListPeople over")
}

func HandleAddPeople(ctx context.Context, cli protos.AddressBookServiceClient) {
	res, err := cli.AddPeople(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := int32(1); i < 11; i++ {
		res.Send(&protos.Person{
			Name: time.Now().GoString(),
			Id:   i,
		})
	}
	res.CloseSend()

	for {
		tmp, err := res.Recv()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Println(err)
				return
			}
		}
		fmt.Println("Ttmp: ", tmp.GetOk())
	}
	fmt.Println("Over")
}
