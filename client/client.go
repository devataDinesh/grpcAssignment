package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpcImplementation/proto"
	"io"
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
		//fmt.Println("hello")
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
