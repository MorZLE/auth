# auth

## GRPC сервис авторизации
### Описание
Проект "auth" представляет собой GRPC сервис, который обеспечивает функциональность авторизации и регистрации пользователей для приложений.

Он также предоставляет возможность проверки прав доступа пользователей, создания и удаления администраторов.

## Установка
Склонируйте репозиторий на свой локальный компьютер.

Перейдите в директорию проекта.

Конфигурация
Откройте файл config.yaml и настройте следующие параметры:

``` yaml
env: "local"  # Окружение проекта
storage_path: "./storage/auth.db"  # Путь к файлу базы данных
token_ttl: 1h  # Время жизни токена доступа
grpc:
  port: 4044  # Порт для gRPC-сервера
  timeout: 5s  # Таймаут для gRPC-запросов
```
## Запуск
Запустите базу данных, если это требуется.

Запустите сервер авторизации, используя следующую команду:
```go
go run main.go --config=config.yaml
```
Сервер авторизации будет запущен и будет доступен для использования.
## Использование
Для использования сервиса авторизации, вы можете взаимодействовать с ним через GRPC-интерфейс, используя соответствующие методы для регистрации, аутентификации, проверки прав доступа и управления администраторами.

### Примеры запросов GRPC
Регистрация нового пользователя:
```go
syntax = "proto3";

message RegisterRequest {
  string email = 1;
  string password = 2;
}


```
### Аутентификация пользователя:

```go
syntax = "proto3";

message LoginRequest {
  string login = 1;  
  string password = 2;
  int32 app_id = 3; 
}

message LoginResponse {
  string token = 1;
}

```

### Добавление приложения
```go
message AddAppRequest{
  string name = 1;
  string secret = 2;
  string key = 3;
}

message AddAppResponse{
  int32 app_id = 1;
}

```
### Проверка прав доступа пользователя:
```go
syntax = "proto3";

message IsAdminRequest {
  string user_id = 1;
  
}

message IsAdminResponse {
  bool is_admin = 1;
  int32 lvl = 2;
}

```
### Создание администратора:

```go
syntax = "proto3";

message CreateAdminRequest {
  string login = 1;
  int32 lvl = 2;
  string key = 3;
}



```

### Удаление администратора:
```go
syntax = "proto3";

message DeleteAdminRequest {
  string login = 1;
  string key = 2;
}

```
