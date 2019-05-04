set GOPATH=%cd%/../../..

md proto

code_generator.exe -c login_db.json -d ../src/db_gen -p proto/login_db.proto

cd ../third_party/protobuf
protoc.exe --go_out=../../src/db_gen/login_db --proto_path=../../db_define/proto login_db.proto
cd ../../../db_define