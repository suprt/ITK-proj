# user-service


## Тесты

### Юнит-тесты

```bash
# use cases
go test ./internal/application/... -v

# gRPC-хендлеры (mockgen)
go test ./internal/adapters/inbound/grpc/... -v
```
### Интеграционные (testcontainers)

**Репозитории (Postgres):**  
Проверяют SQL, транзакции, маппинг ошибок БД.

```bash
go test ./internal/adapters/outbound/postgres/... -v
```

**gRPC end-to-end:**  
Поднимают реальный gRPC-сервер (bufconn) + реальный Postgres. Проверяют полный цикл: клиент → интерсепторы → хендлер → use case → БД.

```bash
go test ./test/e2e/... -v
```

## Примеры команд

### Регистрация

```bash
grpcurl -plaintext -d '{
  "email": "user@local",
  "username": "user",
  "password": "secret123"
}' localhost:50051 user.v1.UserService/Register
```

Ответ:

```json
{ "userId": "..." }
```

### Логин

```bash
grpcurl -plaintext -d '{
  "email": "user@local",
  "password": "secret123"
}' localhost:50051 user.v1.UserService/Login
```

Ответ:

```json
{
  "userId": "...",
  "accessToken": "eyJ...",
  "refreshToken": "...",
  "accessExpiresIn": "900"
}
```

### Свой профиль

```bash
ACCESS=<access_token>

grpcurl -plaintext -H "authorization: Bearer $ACCESS" \
  -d '{}' localhost:50051 user.v1.UserService/GetSelfProfile
```

### Публичный профиль (без авторизации)

```bash
grpcurl -plaintext -d '{
  "user_id": "<USER_ID>"
}' localhost:50051 user.v1.UserService/GetPublicProfile
```
### Обновление своего профиля

```bash
ACCESS=<access_token>

# Меняем только username
grpcurl -plaintext -H "authorization: Bearer $ACCESS" -d '{
  "username": "new_username"
}' localhost:50051 user.v1.UserService/UpdateProfile

# Меняем только email
grpcurl -plaintext -H "authorization: Bearer $ACCESS" -d '{
  "email": "new@local"
}' localhost:50051 user.v1.UserService/UpdateProfile

# Меняем оба поля
grpcurl -plaintext -H "authorization: Bearer $ACCESS" -d '{
  "email": "new@local",
  "username": "new_username"
}' localhost:50051 user.v1

### Смена пароля

```bash
grpcurl -plaintext -H "authorization: Bearer $ACCESS" -d '{
  "current_password": "secret123",
  "new_password": "newsecret123"
}' localhost:50051 user.v1.UserService/ChangePassword
```

### Обновление refresh-токена

```bash
REFRESH=<refresh_token>

grpcurl -plaintext -d "{
  \"refresh_token\": \"$REFRESH\"
}" localhost:50051 user.v1.UserService/RefreshToken
```

### Админские операции

```bash
ADMIN_ACCESS=<admin_access_token>

# Получить полный профиль юзера
grpcurl -plaintext -H "authorization: Bearer $ADMIN_ACCESS" -d '{
  "user_id": "<USER_ID>"
}' localhost:50051 user.v1.UserService/GetUserByID

# Выдать роль
grpcurl -plaintext -H "authorization: Bearer $ADMIN_ACCESS" -d '{
  "user_id": "<USER_ID>",
  "role": "ROLE_TRADER"
}' localhost:50051 user.v1.UserService/AddRole

# Забрать роль
grpcurl -plaintext -H "authorization: Bearer $ADMIN_ACCESS" -d '{
  "user_id": "<USER_ID>",
  "role": "ROLE_TRADER"
}' localhost:50051 user.v1.UserService/RemoveRole

# Сбросить пароль (возвращает временный пароль)
grpcurl -plaintext -H "authorization: Bearer $ADMIN_ACCESS" -d '{
  "user_id": "<USER_ID>"
}' localhost:50051 user.v1.UserService/AdminResetPassword

# Изменить профиль (email/username) чужого юзера
grpcurl -plaintext -H "authorization: Bearer $ADMIN_ACCESS" -d '{
  "user_id": "<USER_ID>",
  "username": "renamed_by_admin"
}' localhost:50051 user.v1.UserService/AdminUpdateProfile

# Или email
grpcurl -plaintext -H "authorization: Bearer $ADMIN_ACCESS" -d '{
  "user_id": "<USER_ID>",
  "email": "newemail@local"
}' localhost:50051 user.v1.UserService/AdminUpdateProfile
```