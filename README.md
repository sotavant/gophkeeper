# proto

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/user.proto -I=proto

# todo
1. обновление данных на сервере по таймауту