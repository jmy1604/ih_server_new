set GOPATH=%cd%/../../..

md proto

code_generator.exe -c game_db.json -d ../src/db_gen -p proto/game_db.proto

cd ../third_party/protobuf
protoc.exe --go_out=../../src/db_gen/game_db --proto_path=../../db_define/proto game_db.proto
cd ../../../db_define