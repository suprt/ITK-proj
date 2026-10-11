# trading

Сейчас для теста работоспособности использовать  
```text 
grpcurl -plaintext -d '{"email":"a@b.c","password":"secret123"}' \
localhost:50051 user.v1.UserService/RegisterUseCase
```

Логин админа
```text 
grpcurl -plaintext -d '{"email":"admin@local","password":"admin12345"}' \
localhost:50051 user.v1.UserService/Login
```

Логин юзера
```text
grpcurl -plaintext -d '{"email":"user@local","password":"user12345"}' \
localhost:50051 user.v1.UserService/Login
```

GetSelfProfile с чьим-либо токеном
```text
grpcurl -plaintext -H 'authorization: Bearer <TOKEN>' \
-d '{}' localhost:50051 user.v1.UserService/GetSelfProfile
```

GetUserByID — админ читает обычного юзера
```text
grpcurl -plaintext -H 'authorization: Bearer <ADMIN_ACCESS>' \
-d '{"user_id":"<USER_ID>"}' \
localhost:50051 user.v1.UserService/GetUserByID
```

GetUserByID — обычный юзер пытается читать кого-то → PermissionDenied
```text
grpcurl -plaintext -H 'authorization: Bearer <USER_ACCESS>' \
-d '{"user_id":"<ADMIN_ID>"}' \
localhost:50051 user.v1.UserService/GetUserByID
```

GetUserByID без токена → Unauthenticated

```text
grpcurl -plaintext -d '{"user_id":"<USER_ID>"}' \
localhost:50051 user.v1.UserService/GetUserByID
```