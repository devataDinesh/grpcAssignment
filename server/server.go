package main

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"grpcImplementation/proto"
	"net"
)

type server struct {
	proto.UnimplementedCalculatorServiceServer
}

func (*server) Sum(ctx context.Context, request *proto.SumRequest) (*proto.SumResponse, error) {

	fmt.Println("*** Sum Method was called ***")

	if request.GetNum1() == 0 && request.GetNum2() == 0 {
		return nil, errors.New("both the Numbers are Zero")
	}

	return &proto.SumResponse{
		Result: request.GetNum1() + request.GetNum2(),
	}, nil
}

func (*server) PrimeNumbers(request *proto.PrimeNumbersRequest, response proto.CalculatorService_PrimeNumbersServer) error {
	fmt.Println("*** PrimeNumbers Method was called ***")

	if request.GetNum() == 0 {
		return errors.New("number is Zero")
	}

	if request.GetNum() == 1 {
		return errors.New("number is One")
	}

	if request.GetNum() == 2 {
		response.Send(&proto.PrimeNumbersResponse{
			Result: 0,
		})
		return nil
	}

	var i int64
	for i = 2; i < request.GetNum(); i++ {
		isPrime := checkIsPrime(i)
		fmt.Println("Value of i is ", i, "isPrime", isPrime)
		if isPrime {
			response.Send(&proto.PrimeNumbersResponse{
				Result: i,
			})
		}
	}
	return nil
}

func checkIsPrime(num int64) bool {
	var j int64
	for j = 2; j < num-1; j++ {
		if num%j == 0 {
			return false
		}
	}
	return true

}

func main() {

	fmt.Println("Starting the server!!!")

	listen, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		fmt.Println("Server Failed to Listen: ", err)
	}

	s := grpc.NewServer()
	proto.RegisterCalculatorServiceServer(s, &server{})

	// //Register reflection service on gRPC server
	// reflection.Register(s)

	if err = s.Serve(listen); err != nil {
		fmt.Println("Server Failed to Serve : ", err)
	}
}