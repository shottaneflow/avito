# avito

Интеграционные тесты для API объявлений Avito QA.

## Структура проекта
```
.
├── internal/
│   ├── client/http/advertisement/   # HTTP-клиенты (create, get, delete)
│   ├── constants/path/              # Пути к эндпоинтам
│   ├── managers/advertisement/      # Бизнес-логика + models/
│   └── runner/                      # Базовый URL
├── tests/
│   └── scenarios/
│       ├── createAdvertisement/     # TC-POST
│       ├── getAdvertisement/        # TC-GET
│       ├── getStatistic/            # TC-STAT
│       ├── getSellerItems/          # TC-SELLER
│       └── e2e/                     # TC-E2E-01
├── .golangci.yml
├── BUGS.md
└── README.md
```

## Зависимости
```bash
go mod tidy
```

| Библиотека | Назначение |
|---|---|
| `github.com/stretchr/testify` | assertions, require |
| `github.com/ozontech/allure-go` | Allure-репорты |

## Запуск тестов
```bash
# Установить зависимости
go mod tidy

# Все тесты
go test ./tests/... -v

# Отдельные сценарии
go test ./tests/scenarios/createAdvertisement/... -v
go test ./tests/scenarios/getAdvertisement/... -v
go test ./tests/scenarios/getStatistic/... -v
go test ./tests/scenarios/getSellerItems/... -v
go test ./tests/scenarios/e2e/... -v


# Allure
allure serve ./allure-results
```

## Линтер и форматтер

### Запуск
```bash
# Проверка линтером
golangci-lint run ./...

# Только форматирование
gofmt -w .
goimports -w .
```

Конфигурация линтера — `.golangci.yml`.

## Известные баги

Подробно описаны в [BUGS.md](./BUGS.md).

| ID | Эндпоинт | Описание | Severity |
|---|---|---|---|
| BUG-01 | POST /item | Ответ содержит только строку вместо объекта | Critical |
| BUG-02 | POST /item | Отрицательный price принимается | Major |
| BUG-03 | POST /item | Отрицательная статистика принимается | Major |
| BUG-04 | POST /item | Поведение без поля statistics не определено | Major |
| BUG-05 | GET /item/{id} | Возвращает массив вместо объекта | Major |
| BUG-06 | GET /statistic/{id} | Возвращает массив вместо объекта | Minor |

## Teardown

Каждый suite, создающий объявления, удаляет их через `AfterAll` / `AfterEach` с помощью `DELETE /api/2/item/:id`. E2E тест выполняет удаление как явный последний шаг.