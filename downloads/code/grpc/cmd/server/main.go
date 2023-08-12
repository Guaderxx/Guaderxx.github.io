package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/Guaderxx/grpc-demo/pkg/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	protos.UnimplementedAddressBookServiceServer
}

func (s *server) AddPerson(ctx context.Context, p *protos.Person) (*protos.Status, error) {
	name := p.GetName()
	if name == "ader" {
		return &protos.Status{
			Ok: true,
		}, nil
	}
	return nil, status.Errorf(codes.FailedPrecondition, "don't support person with this name")
}

func (s *server) ListPeople(req *protos.ListReq, ser protos.AddressBookService_ListPeopleServer) error {
	num := req.GetListnum()
	var i int32 = 0
	for i = 0; i < num; i++ {
		res := &protos.AddressBook{
			People: []*protos.Person{&protos.Person{
				Name:  time.Now().String(),
				Id:    i,
				Email: "",
			}},
		}
		ser.Send(res)
	}
	return nil
}

func (s *server) AddPeople(ser protos.AddressBookService_AddPeopleServer) error {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for {
			pp, err := ser.Recv()
			if err == io.EOF {
				wg.Done()
				break
			}
			handleRecv(pp)
		}
	}()

	for i := int32(0); i < 10; i++ {
		ser.Send(&protos.Status{
			Ok: true,
		})
	}

	wg.Wait()
	fmt.Println("SToS over")
	return nil
}

func handleRecv(p *protos.Person) {
	fmt.Println("name:", p.GetName(), " ID: ", p.GetId())
}

func main() {
	listener, _ := net.Listen("tcp", ":8082")
	ser := grpc.NewServer()
	protos.RegisterAddressBookServiceServer(ser, new(server))
	ser.Serve(listener)
}
