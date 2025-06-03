# goCalcProj

Сервис-калькулятор на основе gRPC с параллельной обработкой инструкций, написанный на Go. Проект поддерживает базовые арифметические операции и предоставляет Swagger UI для изучения.

## Возможности

- Параллельная обработка инструкций вычислений с использованием горутин.
- Поддержка операций +, - и *.
- API через gRPC и http.
- Swagger UI для документации API.
- Развёртывание с помощью Docker через docker-compose.

## Структура проекта

| Директория/Файл                | Описание                                                                 |
|--------------------------------|--------------------------------------------------------------------------|
| `/api/`                        | Содержит `calculator.proto`. |
| `/cmd/client/`                 | Клиентская часть для взаимодействия с сервером (`main.go`).              |
| `/cmd/server/`                 | Серверная часть, точка входа для gRPC и HTTP (`main.go`).                |
| `/gen/openapi/`                | Спецификация API (`calculator.swagger.json`) для Swagger UI.             |
| `/gen/pb/`                     | Сгенерированные protobuf-структуры (`calculator.pb.go`).                 |
| `/internal/calculator/`        | Бизнес-логика: `builder.go`, `resolver.go`, `graph.go`, `types.go`. |
| `/Dockerfile`                  | Конфигурация для сборки Docker-образа `calc-server`.                     |
| `/docker-compose.yml`          | Настройка сервисов `calc-server` и `swagger-ui`.                         |
| `/go.mod`, `/go.sum`           | Управление зависимостями.                            |
| `/README.md`                   | Описание проекта.                                                   |

### Основные компоненты
- **gRPC-сервис (`CalculatorService`)**: Определён в `api/calculator.proto`. Метод `ProcessInstructions` принимает список инструкций и возвращает результаты.
- **gRPC-gateway**: В `cmd/server/main.go` преобразует gRPC-запросы в HTTP, позволяя клиентам (например, Swagger UI) отправлять запросы через RESTful API.
- **Swagger UI**: Развёрнут через Docker, монтирует `calculator.swagger.json` для отображения и тестирования API.
- **Бизнес-логика**: Реализована в `internal/calculator`:
  - `builder.go`: Обрабатывает инструкции, управляет переменными и результатами.
  - `graph.go`: Реализует алгоритм поиска циклических зависимостей (DFS).
  - `types.go`: Определяет вспомогательные типы для логики.
  - `resolver.go`: Выполняет преобразование операнда в значение.

### **Пример работы**:  
### 1
input
```json
[
  { "type": "calc", "op": "+", "var": "x", "left": 1,  "right": 2 },
  { "type": "print", "var": "x" }
]
```

output
```json
{ "items": [
    { "var": "x","value": 3}
  ]
}
```

### 2
input
```json
[
  { "type": "calc", "op": "+", "var": "x",   "left": 10,  "right": 2  },
  { "type": "print",             "var": "x"                     },
  { "type": "calc", "op": "-", "var": "y",   "left": "x",  "right": 3  },
  { "type": "calc", "op": "*", "var": "z",   "left": "x",  "right": "y" },
  { "type": "print",             "var": "w"                     },
  { "type": "calc", "op": "*", "var": "w",   "left": "z",  "right": 0  }
]
```

output
```json
{ "items": [
    {"var": "x","value": 12},
    {"var": "w","value": 0}
  ]
}
```

### 3
input
```json
[
  { "type": "calc", "op": "+", "var": "x",        "left": 10,   "right": 2    },
  { "type": "calc", "op": "*", "var": "y",        "left": "x",  "right": 5    },
  { "type": "calc", "op": "-", "var": "q",        "left": "y",  "right": 20   },
  { "type": "calc", "op": "+", "var": "unusedA",  "left": "y",  "right": 100  },
  { "type": "calc", "op": "*", "var": "unusedB",  "left": "unusedA", "right": 2 },
  { "type": "print",             "var": "q"                        },
  { "type": "calc", "op": "-", "var": "z",        "left": "x",  "right": 15   },
  { "type": "print",             "var": "z"                        },
  { "type": "calc", "op": "+", "var": "ignoreC",  "left": "z",  "right": "y"  },
  { "type": "print",             "var": "x"                        }
]
```

output
```json
{ "items": [
    {"var": "q","value": 40},
    {"var": "z","value": -3},
    {"var": "x","value": 12}
  ]
}
```

