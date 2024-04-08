# Golang send USDC example

This example shows how to send USDC with MPCVault's API using Golang.

Before running this example, you need to execute the following:

```bash

# Install the protocol compiler plugins for Go
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

# Pull protobuf repository
git clone https://github.com/mpcvault/mpcvaultapis.git

# Generate gRPC code
protoc --go_out=./ --go-grpc_out=./ ./mpcvaultapis/mpcvault/platform/v1/*.proto
```# MPCVault

#FIRST

go get github.com/aws/aws-lambda-go/lambda
go install github.com/aws/aws-lambda-go/cmd/build-lambda-zip@latest

#SECOND IN POWERSHELL

$env:GOOS = "linux"
$env:GOARCH = "amd64"
$env:CGO_ENABLED = "0"
go build -tags lambda.norpc -o bootstrap main.go
~\Go\Bin\build-lambda-zip.exe -o myFunction.zip bootstrap


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
		Amount:               "100",                                        // 1 USDC
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