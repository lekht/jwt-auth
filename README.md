# authorization-api

Приложение реализует методы аутентификации. Для запуска приложения склонируйте его к себе на пк и выполните следующую команду

```go run ./cmd/service/service.go``` 

## SignUp
Метод создания учетной записи. Принимает json структуру с полями логина и пароля. Для вызова метода выполните команду в терминале.
```
curl -i -H "Content-Type: application/json" -d '{"login":"user1","password":"password"}' -X POST http://localhost:3000/auth/signup && echo
```

## SignIn с неверным паролем
Метод попытки авторизации с неверным паролем. Необходим для обновления аудита по пользователю и проверки функционала по блокировке пользователя.
```
curl -i -H "Content-Type: application/json" -d '{"login":"user1","password":"pasword"}' -X POST http://localhost:3000/auth/signin && echo
```

## SignIn с правильным паролем
Для успешной аутентификации выполните следующую команду
```
curl -i -H "Content-Type: application/json" -d '{"login":"user1","password":"password"}' -X POST http://localhost:3000/auth/signin && echo
```

## Получение аудита
Возращает весь аудит по пользователю. Принимает header X-Token, который возвращается в теле SignIn запроса.
```
curl -i -H "X-Token: <token_here>" -X GET http://localhost:3000/auth/history && echo
```

## Очищение аудита
Очищает аудит пользователя. Принимает header X-Token, который возвращается в теле SignIn запроса 
```curl -i -H "X-Token: <token_here>" -X DELETE http://localhost:3000/auth/history && echo
```