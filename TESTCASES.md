# Тест кейсы для задания 2.1

---

# 1. POST /api/1/item

## TC-POST-01 Успешное создание объявления

**Предусловия:**
1. В системе есть seller с id=111289
**Шаги:**
1. Отправить POST запрос с body:
```json
{
  "sellerID": 111289,
  "name": "test",
  "price": 1000,
  "statistics": {
    "likes": 10,
    "viewCount": 10,
    "contacts": 10
  }
}
```

**Ожидаемый результат:**
- 200 OK
- В ответе присутствует id
- Данные совпадают с отправленными

---

## TC-POST-02 Отсутствует sellerID

**Шаги:**
1. Отправить POST запрос с body:
```json
{
  "name": "test",
  "price": 1000,
  "statistics": {
    "likes": 10,
    "viewCount": 10,
    "contacts": 10
  }
}
```

**Ожидаемый результат:**
- 400 Bad Request

---

## TC-POST-03 Неверный тип данных price

**Шаги:**
1. Отправить POST запрос с body:
```json
{
  "sellerID": 111289,
  "name": "test",
  "price": "1sb",
  "statistics": {
    "likes": 10,
    "viewCount": 10,
    "contacts": 10
  }
}
```

**Ожидаемый результат:**
- 400 Bad Request

---

## TC-POST-04 Нулевая цена(price)

**Шаги:**
1. Отправить POST запрос с body:
```json
{
  "sellerID": 111289,
  "name": "test",
  "price": 0,
  "statistics": {
    "likes": 10,
    "viewCount": 10,
    "contacts": 10
  }
}
```

**Ожидаемый результат:**
- 400 Bad Request

---

## TC-POST-05 Очень большое значение price

**Шаги:**
1. Отправить POST запрос с body:
```json
{
  "sellerID": 111289,
  "name": "test",
  "price": 9223372036854775808,
  "statistics": {
    "likes": 10,
    "viewCount": 10,
    "contacts": 10
  }
}
```

**Ожидаемый результат:**
- 400 Bad Request

---

## TC-POST-06 Пустое имя

**Шаги:**
1. Отправить POST запрос с body:
```json
{
  "sellerID": 111289,
  "name": "",
  "price": 100,
  "statistics": {
    "likes": 10,
    "viewCount": 10,
    "contacts": 10
  }
}
```

**Ожидаемый результат:**
- 400 Bad Request 

---

## TC-POST-07 Создание объявления с нулевой статистикой

**Шаги:**
1. Отправить POST запрос с body:
```json
{
  "sellerID": 111289,
  "name": "test",
  "price": 100,
  "statistics": {
    "likes": 0,
    "viewCount": 0,
    "contacts": 0
  }
}
```

**Ожидаемый результат:**
- 200 OK 

# 2. GET /api/1/item/{id}

## TC-GET-01 — Получение существующего объявления

**Предусловия:**
1. В системе есть объявление с id = 6eb9b015-b2cb-4502-8b61-c58903106073

---

## TC-GET-02 — Несуществующий id
**Тип:** Негативный  
**Техника:** Equivalence Partitioning  

**Ожидаемый результат:**
- 404 Not Found

---

## TC-GET-03 — Некорректный id
**Тип:** Негативный  
**Техника:** Equivalence Partitioning  

---

## TC-GET-04 — Пустой id
**Тип:** Негативный  
**Техника:** Equivalence Partitioning  

---

# 3. DELETE /api/2/item/{id}

## TC-DEL-01 — Успешное удаление
**Тип:** Позитивный  
**Техника:** State Transition  

---

## TC-DEL-02 — Повторное удаление
**Тип:** Corner case  
**Техника:** State Transition  

---

## TC-DEL-03 — Несуществующий id
**Тип:** Негативный  
**Техника:** Equivalence Partitioning  

---

## TC-DEL-04 — Некорректный id
**Тип:** Негативный  
**Техника:** Equivalence Partitioning  

---

# 4. GET /api/1/statistic/{id}

## TC-STAT-01 — Получение статистики
**Тип:** Позитивный  
**Техника:** Equivalence Partitioning  

---

## TC-STAT-02 — Несуществующий id
**Тип:** Негативный  
**Техника:** Equivalence Partitioning  

---

## TC-STAT-03 — Проверка структуры ответа
**Тип:** Нефункциональный  
**Техника:** Error Guessing  

**Проверки:**
- Поля: likes, viewCount, contacts
- Тип данных: integer

---

# 5. GET /api/1/{sellerID}/item

## TC-SELLER-01 — Получение объявлений пользователя
**Тип:** Позитивный  
**Техника:** Equivalence Partitioning  

---

## TC-SELLER-02 — Пользователь без объявлений
**Тип:** Corner case  
**Техника:** Equivalence Partitioning  

---

## TC-SELLER-03 — Некорректный sellerID
**Тип:** Негативный  
**Техника:** Equivalence Partitioning  

---

# Нефункциональные тесты

## TC-NF-01 — Время ответа
**Тип:** Нефункциональный  
- Ответ API < 500 мс

---

## TC-NF-02 — Нагрузочное тестирование
**Тип:** Нефункциональный  
- 100+ запросов без 500 ошибок

---

## TC-NF-03 — Проверка заголовков
**Тип:** Нефункциональный  
- Content-Type: application/json

