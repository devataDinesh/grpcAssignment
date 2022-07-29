package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpcImplementation/proto"
	"io"
	"time"
)

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		fmt.Println("could not connect: ", err)
	}
	defer cc.Close()

	c := proto.NewCalculatorServiceClient(cc)

	//Unary Method
	callSum(c)

	//Server-Side Streaming method
	callPrimeNumbers(c)

	//Client-Side Streaming method
	callComputeAverage(c)
}

func callSum(c proto.CalculatorServiceClient) {

	fmt.Println("Calling calculator.Sum() Method")

	res, err := c.Sum(context.Background(), &proto.SumRequest{
		Num1: 2,
		Num2: 3,
	})

	if err != nil {
		fmt.Println("Error in the Sum() ", err.Error())
		return
	}

	fmt.Println("Result from the Sum() is : ", res.Result)

}

func callPrimeNumbers(c proto.CalculatorServiceClient) {

	fmt.Println("\nCalling calculator.PrimeNumbers() Method")

	res, err := c.PrimeNumbers(context.Background(), &proto.PrimeNumbersRequest{
		Num: 20,
	})

	if err != nil {
		fmt.Println("Error in the PrimeNumbers() ", err.Error())
		return
	}

	fmt.Print("Result from the PrimeNumbers() is : ")
	for {
		msg, err := res.Recv()
		if err == io.EOF { //we have reached to the end of the file
			break
		}
		if err != nil {
			fmt.Println("error while receiving server stream : ", err)
		}
		fmt.Print(msg.Result, "\t")
	}
}

func callComputeAverage(c proto.CalculatorServiceClient) {

	fmt.Println("\n\nCalling calculator.ComputeAverage() Method")

	stream, err := c.ComputeAverage(context.Background())

	if err != nil {
		fmt.Println("Error in ComputeAverage() ", err)
	}

	var requests []*proto.ComputeAverageRequest

	requests = append(requests, &proto.ComputeAverageRequest{
		Num: 2,
	})
	requests = append(requests, &proto.ComputeAverageRequest{
		Num: 3,
	})
	requests = append(requests, &proto.ComputeAverageRequest{
		Num: 4,
	})
	requests = append(requests, &proto.ComputeAverageRequest{
		Num: 5,
	})
	requests = append(requests, &proto.ComputeAverageRequest{
		Num: 6,
	})
	requests = append(requests, &proto.ComputeAverageRequest{
		Num: 7,
	})
	requests = append(requests, &proto.ComputeAverageRequest{
		Num: 8,
	})

	for _, request := range requests {
		stream.Send(request)
		time.Sleep(100 * time.Millisecond)
	}
	response, err := stream.CloseAndRecv()
	if err != nil {
		fmt.Println("Error while receiving response from server : ", err)
	}
	fmt.Println("Result from the ComputeAverage() : ", response.GetResult())

}
