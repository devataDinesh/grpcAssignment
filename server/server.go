package main

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"grpcImplementation/proto"
	"io"
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
		if isPrime {
			response.Send(&proto.PrimeNumbersResponse{
				Result: i,
			})
		}
	}
	return nil
}

func (*server) ComputeAverage(request proto.CalculatorService_ComputeAverageServer) error {
	fmt.Println("*** ComputeAverage Method was called ***")

	var sum int64 = 0
	var cnt int64 = 0

	for {
		msg, err := request.Recv()
		if err == io.EOF { //we have finished reading client stream
			return request.SendAndClose(&proto.ComputeAverageResponse{
				Result: sum / cnt,
			})
		}
		if err != nil {
			fmt.Println("Error while reading client stream : ", err)
		}

		sum = sum + msg.Num
		cnt = cnt + 1
	}
}

func (*server) FindMaxNumber(stream proto.CalculatorService_FindMaxNumberServer) error {

	fmt.Println("*** FindMaxNumber Method was called ***")

	var max int64 = 0

	for {
		request, err := stream.Recv()

		if err == io.EOF {
			return nil
		}
		if err != nil {
			fmt.Println("error while receiving data from  client : ", err)
			return err
		}
		if max < request.GetNum() {
			max = request.GetNum()
			sendErr := stream.Send(&proto.FindMaxAverageResponse{
				Result: max,
			})
			if sendErr != nil {
				fmt.Println("Error while sending Response to the client ", err)
				return err
			}
		}
	}
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
