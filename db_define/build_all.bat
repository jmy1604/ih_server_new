set GOPATH=%cd%/../../../
code_generator.exe -c json/game_db.json -d ../src/db_gen
code_generator.exe -c json/login_db.json -d ../src/db_gen
code_generator.exe -c json/rpc_db.json -d ../src/db_gen