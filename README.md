[![Go Coverage](https://github.com/JMURv/effectiveMobile/wiki/coverage.svg)](https://raw.githack.com/wiki/JMURv/sso/coverage.html)

### Локальный запуск
1. Создать `.env` файл по примеру из `.env.example`
2. `git clone https://github.com/JMURv/effective_mobile`
3. `go run cmd/main.go`

### Запуск в docker
1. Создать `.env` файл по примеру из `.env.example`
2. Из директории с созданным файлом выполнить команду:

```docker run -d -p 8080:8080 -p 8081:8081 -v /path/to/your/.env:/app/.env jmurv/effective_mobile:latest```

Note: на 8081 порту крутится внешнее API, на которое по заданию требуется делать запрос