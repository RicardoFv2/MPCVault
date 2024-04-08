package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/golang/protobuf/ptypes/wrappers"
	apipb "go.mpcvault.com/go.mpcvault.com/genproto/mpcvaultapis/platform/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"log"
)

const APIEndpoint string = "api.mpcvault.com:443"
const APIToken = "RB/obe/9p6sXh4t+PM5wflmqrigf3D4/GCYcE7uua2g="
const VaultUUID = "417559a4-4b65-4e81-beff-cdeb47d94090"
const CallbackClientSignerPublicKey = "AAAAC3NzaC1lZDI1NTE5AAAAIJf/bxaO9FLWy4b06wh5xYXoefW97wfIkit+Gbe5h53Y"

func EVMSendERC20(ctx context.Context, client apipb.PlatformAPIClient) (string, error) {
	// Prepare and execute your request here
	// Simplified for brevity - implement your request logic
	req := apipb.CreateSigningRequestRequest{}
	createSigningType := &apipb.CreateSigningRequestRequest_EvmSendErc20{}
	req.Type = createSigningType
	createSigningType.EvmSendErc20 = &apipb.EVMSendERC20{
		ChainId:              137,
		From:                 "0xACe27d1AEc109c5cDd2C5fE3Ad6bfe281A2F27fE",
		To:                   "0x7C696d20DD81056735F5A86170D5d0Abf2566483",
		TokenContractAddress: "0xa2F4eB0E3838dc6726138f6202D06b213f22f291",
		Amount:               "1",
		GasFee: &apipb.EVMGas{
			GasLimit: &wrappers.StringValue{Value: "175000"},
		},
		Nonce: &wrappers.Int64Value{Value: 23},
	}
	req.Notes = &wrappers.StringValue{Value: "Ricardo desde Lambda"}

	res, err := client.CreateSigningRequest(ctx, &req) // Use the passed client parameter
	if err != nil {
		log.Fatalf("err: %v", err)
	}

	b, _ := json.MarshalIndent(res, "", " ")
	fmt.Println(string(b))
	return string(b), nil
}

func Handler(ctx context.Context, apiGatewayEvent events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Initialize gRPC connection
	creds := credentials.NewClientTLSFromCert(nil, "")
	conn, err := grpc.Dial(APIEndpoint, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
		return events.APIGatewayProxyResponse{}, fmt.Errorf("failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	grpcClient := apipb.NewPlatformAPIClient(conn)
	header := metadata.New(map[string]string{"x-mtoken": APIToken})
	grpcCtx := metadata.NewOutgoingContext(ctx, header)

	result, err := EVMSendERC20(grpcCtx, grpcClient)
	if err != nil {
		return events.APIGatewayProxyResponse{}, fmt.Errorf("failed to execute EVMSendERC20: %v", err)
	}

	// Set CORS headers
	headers := map[string]string{
		"Access-Control-Allow-Origin":  "*", // Allow requests from any origin
		"Access-Control-Allow-Headers": "Content-Type",
		"Access-Control-Allow-Methods": "OPTIONS, POST",
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    headers,
		Body:       fmt.Sprintf("Lambda execution successful: %s", result),
	}, nil
}

func main() {
	lambda.Start(Handler)
}
