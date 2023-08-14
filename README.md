# authorization-api

curl -i -H "Content-Type: application/json" -d '{"login":"user1","password":"password"}' -X POST http://localhost:3000/auth/signup && echo

curl -i -H "Content-Type: application/json" -d '{"login":"user1","password":"pasword"}' -X POST http://localhost:3000/auth/signin && echo

curl -i -H "Content-Type: application/json" -d '{"login":"user1","password":"password"}' -X POST http://localhost:3000/auth/signin && echo

curl -i -H "X-Token: <token_here>" -X GET http://localhost:3000/auth/history && echo

curl -i -H "X-Token: <token_here>" -X DELETE http://localhost:3000/auth/history && echo