set PROTO_PATH=.\rental\api
set GO_OUT_PATH=.\rental\api\gen\v1

mkdir %GO_OUT_PATH%

protoc -I=%PROTO_PATH% --go_out=plugins=grpc,paths=source_relative:%GO_OUT_PATH% rental.proto
protoc -I=%PROTO_PATH% --grpc-gateway_out=paths=source_relative,grpc_api_configuration=%PROTO_PATH%\rental.yaml:%GO_OUT_PATH% rental.proto