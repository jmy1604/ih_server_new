set GOPATH=%cd%/../../..

md proto

code_generator.exe -c rpc_db.json -d ../src/db_gen -p proto/rpc_db.proto

cd ../third_party/protobuf
protoc.exe --go_out=../../src/db_gen/rpc_db --proto_path=../../db_define/proto rpc_db.proto
cd ../../../db_define