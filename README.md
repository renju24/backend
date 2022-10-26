# API v1
Backend содержит в себе HTTP REST API и WebSocket-сервер.

HTTP REST API будет использоваться только для двух функций:
1. Регистрация
2. Авторизация

После успешной регистрации/авторизации установится WebSocket-соединение и дальнейшее общение с сервером будет происходить только по протоколу WebSocket.

# HTTP REST API
### Обработка ошибок от сервера
В случае ошибки сервер вернет HTTP статус код не равный 200 и объект ошибки в виде JSON:
```json
{
  "status": 0,
  "error": {
    "code": 105,
    "description": "user not found"
  }
}
```
где
status = 0 означает ошибку (в случае успеха status будет 1)

error - объект ошибки, где code - это числовой код ошибки, а description - описание ошибки на английском языке.

### Метод регистрации

```POST https://renju24.com/api/v1/sign_up```

Тело запроса (JSON):
```json
{
	"username": "username",
	"email": "email@email.com",
	"password": "password1",
	"repeated_password": "password1"
}
```
Тело ответа (JSON):
```json
{
	"status": 1,
	"token": "<JWT_TOKEN>"
}
```
Полученный JWT_TOKEN нужно надежно сохранить в мобильном приложении и добавлять его в заголовок Authorization: Bearer <JWT_TOKEN> при последующих запросах.

На фронте ничего делать не надо, так как сервер сам установит токен в защищенную http-only куку.

Возможные коды ошибок и их описания:
```
Код 100 (invalid JSON body)             Неправильный JSON в теле запроса. Ошибка на стороне приложения или сайта.
Код 101 (internal server error)         Ошибка на стороне сервера (например, если база данных упала).
Код 200 (username is required)          При регистрации не прислали username.
Код 202 (email is required)             При регистрации не прислали email.
Код 204 (password is required)          При регистрации не прислали пароль.
Код 205 (repeated_password is required) При регистрации не прислали повторный пароль.
Код 201 (username is already taken)     Пользователь с таким username уже зарегистрирован.
Код 203 (email is already taken)        Пользователь с таким email уже зарегистрирован.
Код 206 (invalid username length)       Невалидная длина username. Разрешено от 4 до 32 символов.
Код 207 (invalid email)                 Невалидный email. Например, если не содержит символ @.
Код 208 (invalid email length)          Невалидная длина email. Разрешено от 5 до 84 символов.
Код 209 (invalid password length)       Невалидная длина пароля. Разрешено от 8 до 64 символов.
Код 210 (invalid password character)    Пароль содержит недопустимые символы. Разрешены только латиница и цифры.
Код 211 (missing letter character)      Пароль должен содержать хотя бы одну букву.
Код 214 (missing digit character)       Пароль должен содержать хотя бы одну цифру.
Код 215 (passwords are not equal)       Пароли не совпадают.
```
### Метод авторизации (вход)

```POST https://renju24.com/api/v1/sign_in```

Тело запроса (JSON):
```json
{
	"login": "username или email",
	"password": "пароль"
}
```
Тело ответа (JSON):
```json
{
	"status": 1,
	"token": "<JWT_TOKEN>"
}
```

Возможные коды ошибок и их описания:
```
Код 100 (invalid JSON body)     Неправильный JSON в теле запроса. Ошибка на стороне приложения или сайта.
Код 101 (internal server error) Ошибка на стороне сервера (например, если база данных упала).
Код 103 (invalid credentials)   Неправильный пароль.
Код 105 (user not found)        Пользователь с таким логином не найден.
```
