[![Go Coverage](https://github.com/JMURv/effective_mobile/wiki/coverage.svg)](https://raw.githack.com/wiki/JMURv/effective_mobile/coverage.html)

### Локальный запуск
1. Создать `.env` файл по примеру из `.env.example`
2. `git clone https://github.com/JMURv/effective_mobile.git .`
3. `go run cmd/main.go`

### Запуск в docker
1. Создать `.env` файл по примеру из `.env.example`
2. Из директории с созданным файлом выполнить команду:

```docker run -d -p 8080:8080 -p 8081:8081 -v $(pwd)/.env:/app/.env jmurv/effective_mobile:latest```

Note: на 8081 порту крутится внешнее API, на которое по заданию требуется делать запрос

### Почему текст песни хранится в списке?
В данном случае текст песни хранится `в списке`, потому что требуется `пагинация по его частям`. Если бы текст был обычной строкой, то организовать пагинацию стало бы намного сложнее, так как это требовало бы разделения текста по разрыву строки `\n\n`