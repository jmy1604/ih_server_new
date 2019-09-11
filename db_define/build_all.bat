set GOPATH=%cd%/../../..
code_generator.exe -c json/game_db.json -d ../src/db_gen -p ../third_party/protobuf/protoc.exe
code_generator.exe -c json/login_db.json -d ../src/login_server/db -p ../third_party/protobuf/protoc.exe
code_generator.exe -c json/rpc_db.json -d ../src/rpc_server/db -p ../third_party/protobuf/protoc.exe