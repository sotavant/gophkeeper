# proto

protoc --go_out=proto --go_opt=paths=source_relative --go-grpc_out=proto --go-grpc_opt=paths=source_relative proto/user.proto -I=proto

# todo
1. обновление данных на сервере по таймауту

# files
https://dev.to/dimk00z/grpc-file-transfer-with-go-1nb2

# questions

1. Шифрование
2. Клиент и Сервер в одном репозитории
3. tls, зашивать сертификат в репозиторий? (https://dev.to/techschoolguru/how-to-secure-grpc-connection-with-ssl-tls-in-go-4ph)