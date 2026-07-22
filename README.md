# 💳 Simple Wallet API

Учебный REST API сервис для управления электронными кошельками, пользователями, переводами и историей транзакций с Jwt аутентификацией.

---

## ⚡️ Запуск проекта

### Вариант 1: Через Docker Compose (Рекомендуемый)

1. Клонируйте репозиторий:
   ```bash
   git clone https://github.com/VLKasabiev/simple-wallet.git
   cd simple-wallet
   ```

2. Запустите приложения и базу данных:
   ```bash
   docker-compose up -d --build
   ```

3. Сервис станет доступен по адресу: `http://localhost:8080`

---

### Вариант 2: Локальный запуск (Go)

1. Установите зависимости и запустите приложение:
   ```bash
   go mod download
   go run cmd/main.go
   ```

---

## 🚀 Основные возможности

- **Аутентификация и безопасность:** Регистрация, авторизация с вычислением JWT-токенов, защита эндпоинтов авторизационным middleware (`Bearer Token`).
- **Управление пользователями:** Создание, получение профилей пользователей.
- **Управление кошельками:**
  - Создание кошельков для пользователя (с валютами MVP: `RUB`, `USD`, `EUR`).
  - Просмотр кошелька и баланса.
- **Финансовые операции (Транзакции):**
  - **Deposit** — пополнение баланса.
  - **Withdraw** — списание средств с проверкой остатка (баланс не уходит в минус).
  - **Transfer** — перевод средств между двумя кошельками.
  - При недостатке средств операция списания/перевода фиксируется со статусом `failed`, а баланс остаётся неизменным.
  - Фильтрация истории транзакций по типу (`type`) и статусу (`status`) а также сортировка по датам по возрастанию (`created_at_desc`) и убыванию (`created_at_desc`).
- **Health Check:** Эндпоинт проверки работоспособности сервиса (`/health`).

---

## 🛠 Технологический стек

* **Язык программирования:** Go (Golang)
* **Веб-фреймворк:** [Echo v4](https://echo.labstack.com/)
* **Аутентификация:** JWT (JSON Web Tokens)
* **База данных:** PostgreSQL (или любая SQL-база данных)
* **Контейнеризация:** Docker & Docker Compose
* **Контроль версий:** Git

---

## 📐 Архитектура и Сущности

### 1. User (Пользователь)
| Поле | Тип | Описание |
| :--- | :--- | :--- |
| `id` | `int` | Уникальный идентификатор |
| `name` | `string` | Имя пользователя |
| `email` | `string` | Уникальный Email |
| `created_at` | `time.Time` | Дата и время создания |
| `updated_at` | `time.Time` | Дата и время обновления |

### 2. Wallet (Кошелёк)
| Поле | Тип | Описание |
| :--- | :--- | :--- |
| `id` | `int` | Уникальный идентификатор |
| `user_id` | `int` | ID владельца кошелька |
| `currency` | `Currency` | Валюта (`RUB`, `USD`, `EUR`) |
| `balance` | `decimal.Decimal` | Текущий баланс кошелька |
| `created_at` | `time.Time` | Дата создания |
| `updated_at` | `time.Time` | Дата обновления |

### 3. Transaction (Транзакция)
| Поле | Тип | Описание |
| :--- | :--- | :--- |
| `id` | `int64` | Уникальный идентификатор |
| `wallet_id` | `int64` | ID кошелька |
| `type` | `TransactionType` | Тип операции (`deposit`, `withdraw`, `transfer_in`, `transfer_out`) |
| `amount` | `decimal.Decimal` | Сумма операции (должна быть > 0) |
| `status` | `TransactionStatus` | Статус операции (`success`, `failed`) |
| `description` | `string` | Описание / комментарий |
| `created_at` | `time.Time` | Дата проведения операции |

---

## 🛣 API Endpoints

### 🔓 Открытые эндпоинты (Public)

| Метод | Эндпоинт | Описание |
| :--- | :--- | :--- |
| `GET` | `/health` | Проверка статуса сервиса |
| `POST` | `/users` | Регистрация нового пользователя |
| `GET` | `/users` | Получение списка всех пользователей |
| `POST` | `/users/login` | Авторизация пользователя (возвращает JWT-токен) |

### 🔒 Защищенные эндпоинты (Protected — Требуется `Authorization: Bearer <token>`)

| Метод | Эндпоинт | Описание |
| :--- | :--- | :--- |
| `GET` | `/users/:id` | Получение информации о пользователе |
| `POST` | `/users/:id/wallets` | Создание нового кошелька для пользователя |
| `GET` | `/users/:id/wallets` | Получение всех кошельков пользователя |
| `GET` | `/wallets/:id` | Получение детальной информации о кошельке |
| `GET` | `/wallets/:id/balance` | Просмотр баланса кошелька |
| `POST` | `/wallets/:id/deposit` | Пополнение баланса кошелька |
| `POST` | `/wallets/:id/withdraw` | Списание средств с кошелька |
| `POST` | `/wallets/:id/transfer` | Перевод средств на другой кошелёк |
| `GET` | `/wallets/:id/transactions` | История транзакций (с возможностью фильтрации) |

---

## 🚦 Примеры использования (HTTP / cURL)

### 1. Проверка статуса сервиса (`GET /health`)
GET /health

---

### 2. Регистрация нового пользователя (`POST /users`)
```http
POST /users
Content-Type: application/json

{
  "name": "User1",
  "email": "user1@example.com",
  "password": "user1"
}
```

---

### 3. Получение списка всех пользователей (`GET /users`)
GET /users

---

### 4. Авторизация (`POST /users/login`)
```http
POST /users/login
Content-Type: application/json

{
  "email": "user1@example.com",
  "password": "user1"
}
```

> **Примечание:** Все последующие запросы к защищенным эндпоинтам должны содержать заголовок:  
> `Authorization: Bearer <your_jwt_token>`

---

### 5. Получение информации о пользователе (`GET /users/:id`)
```http
GET /users/1
Authorization: Bearer <your_jwt_token>
```

---

### 6. Создание кошелька для пользователя (`POST /users/:id/wallets`)
```http
POST /users/1/wallets
Authorization: Bearer <your_jwt_token>
Content-Type: application/json

{
  "currency": "RUB"
}
```

---

### 7. Получение кошельков пользователя (`GET /users/:id/wallets`)
```http
GET /users/1/wallets
Authorization: Bearer <your_jwt_token>
```

---

### 8. Детальная информация о кошельке (`GET /wallets/:id`)
```http
GET /wallets/1
Authorization: Bearer <your_jwt_token>
```

---

### 9. Просмотр баланса кошелька (`GET /wallets/:id/balance`)
```http
GET /wallets/1/balance
Authorization: Bearer <your_jwt_token>
```

---

### 10. Пополнение кошелька (`POST /wallets/:id/deposit`)
```http
POST /wallets/1/deposit
Authorization: Bearer <your_jwt_token>
Content-Type: application/json

{
  "amount": "1000.50"
}
```

---

### 11. Списание средств (`POST /wallets/:id/withdraw`)
```http
POST /wallets/1/withdraw
Authorization: Bearer <your_jwt_token>
Content-Type: application/json

{
  "amount": "300.00"
}
```

---

### 12. Перевод между кошельками (`POST /wallets/:id/transfer`)
```http
POST /wallets/1/transfer
Authorization: Bearer <your_jwt_token>
Content-Type: application/json

{
  "to_wallet_id": 2,
  "amount": "200.00",
  "description": "Payment for services"
}
```

---

### 13. История транзакций (`GET /wallets/:id/transactions`)
```http
GET /wallets/1/transactions?type=withdraw&status=success
Authorization: Bearer <your_jwt_token>
```

---
