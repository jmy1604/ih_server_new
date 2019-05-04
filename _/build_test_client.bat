call set_go_path.bat
go install ih_server_new/src/test_client
if errorlevel 1 goto exit

go build -i -o ../bin/test_client.exe ih_server_new/src/test_client
if errorlevel 1 goto exit

if errorlevel 0 goto ok

:exit
echo build test_client failed!!!!!!!!!!!!!!!!!!!

:ok
echo build test_client ok