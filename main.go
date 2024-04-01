package main

import (
	"context"
	"encoding/json"
	"fmt"
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

var grpcClient apipb.PlatformAPIClient

func main() {
	creds := credentials.NewClientTLSFromCert(nil, "")
	conn, err := grpc.Dial(APIEndpoint, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	header := metadata.New(map[string]string{
		"x-mtoken": APIToken,
	})
	ctx := metadata.NewOutgoingContext(context.Background(), header)

	grpcClient = apipb.NewPlatformAPIClient(conn)

	EVMSendERC20(ctx)
}

func EVMSendERC20(ctx context.Context) {
	req := apipb.CreateSigningRequestRequest{}
	createSigningType := &apipb.CreateSigningRequestRequest_EvmSendErc20{}
	req.Type = createSigningType
	createSigningType.EvmSendErc20 = &apipb.EVMSendERC20{
		ChainId:              137,                                          
		From:                 "0xACe27d1AEc109c5cDd2C5fE3Ad6bfe281A2F27fE", 
		To:                   "0x7C696d20DD81056735F5A86170D5d0Abf2566483", 
		TokenContractAddress: "0xbb15C694f4EFADd1628515280770e5E6f871A625", 
		Amount:               "50",                                    
		GasFee:               nil,                                          
		Nonce:                &wrappers.Int64Value{Value: 40},
	}
	req.Notes = &wrappers.StringValue{Value: "Ricardo"} 

	res, err := grpcClient.CreateSigningRequest(ctx, &req)
	if err != nil {
		log.Fatalf("err: %v", err)
	}

	b, _ := json.MarshalIndent(res, "", " ")
	fmt.Println(string(b))
}

// EVMSendERC20 send erc20 token with client signer
func EVMSendERC20WithClientSigner(ctx context.Context) {
	req := apipb.CreateSigningRequestRequest{}
	createSigningType := &apipb.CreateSigningRequestRequest_EvmSendErc20{}
	req.Type = createSigningType
	createSigningType.EvmSendErc20 = &apipb.EVMSendERC20{
		ChainId:              137,                                          // polygon chain id
		From:                 "0xACe27d1AEc109c5cDd2C5fE3Ad6bfe281A2F27fE", // sender address
		To:                   "0x7C696d20DD81056735F5A86170D5d0Abf2566483", // receiver address
		TokenContractAddress: "0x3c499c542cEF5E3811e1192ce70d8cC03d5c3359", // USDC contract address on polygon
		Amount:               "100",                                    // 1 USDC
		GasFee:               nil,                                          // leave nil to use auto gas settings
	}
	req.Notes = &wrappers.StringValue{Value: "sending 1 USDC for testing"} // setting transaction notes
	req.VaultUuid = &wrappers.StringValue{Value: VaultUUID}
	req.CallbackClientSignerPublicKey = &wrappers.StringValue{Value: CallbackClientSignerPublicKey}

	res, err := grpcClient.CreateSigningRequest(ctx, &req)
	if err != nil {
		log.Fatalf("err: %v", err)
	}

	b, _ := json.MarshalIndent(res, "", " ")
	fmt.Println(string(b))
}
