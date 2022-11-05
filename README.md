## Общая информация
Backend состоит из:
1. HTTP REST API ```https://renju24.com/api/v1/<api_method>```
2. WebSocket API ```wss://renju24.com/connection/websocket```

HTTP REST API используется только для следующих задач:
1. Регистрация: ```POST https://renju24.com/api/v1/sign_up```
2. Авторизация: ```POST https://renju24.com/api/v1/sign_in```
3. OAuth2 авторизация: ```GET https://renju24.com/api/v1/oauth2/<platform>/<service>```
4. Пинг сервера: ```GET https://renju24.com/api/v1/ping```

После успешной регистрации/авторизации и получения токена нужно установить WebSocket-соединение с сервером и дальше общаться через него.

Подробнее в [Документации](https://github.com/renju24/backend/wiki)
