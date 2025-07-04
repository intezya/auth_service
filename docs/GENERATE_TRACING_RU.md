# generate_tracing.go - Генератор трейсинг-обёрток

Утилита для автоматической генерации обёрток с трейсингом для Go структур и их методов.

<!-- TOC -->
- [Возможности](#возможности)
- [Установка и использование](#установка-и-использование)
    - [Основные флаги](#основные-флаги)
    - [Примеры использования](#примеры-использования)
- [Специальные возможности для gRPC](#специальные-возможности-для-grpc)
- [Структура сгенерированного кода](#структура-сгенерированного-кода)
- [Алгоритм работы](#алгоритм-работы)
- [Ограничения](#ограничения)
- [Интеграция](#интеграция)
<!-- /TOC -->


## Возможности

- Автоматическое создание обёрток с трейсингом для любых структур
- Поддержка интерфейсов из разных пакетов
- Корректная обработка named импортов (например, `authpb "github.com/example/auth"`)
- Автоматическое встраивание `mustEmbed` интерфейсов для gRPC сервисов
- Поддержка вложенных интерфейсов в wrapped структурах
- Умное управление импортами (добавляет только используемые)
- Рекурсивный поиск структур в директориях

## Установка и использование

### Основные флаги

```bash
# Обязательные
-struct string            # Имя структуры для обёртывания
-file string              # Go файл со структурой (опционально, можно использовать поиск)

# Конфигурация интерфейса
-interface string         # Имя интерфейса (по умолчанию = имя структуры)
-interface-pkg string     # Пакет интерфейса (если отличается от пакета структуры)

# Конфигурация структуры
-struct-pkg string        # Пакет структуры (если отличается от выходного)

# Выходные параметры
-output string            # Выходной файл (по умолчанию: <struct>_tracing.go)
-package string           # Имя выходного пакета (по умолчанию: как у входного)

# Трейсинг
-tracer string            # Импорт пакета трейсера (по умолчанию: github.com/intezya/auth_service/internal/infrastructure/metrics/tracer)

# Поиск
-search-dir string        # Директория для поиска структуры (по умолчанию: ".")
-recursive                # Рекурсивный поиск в поддиректориях

# Встраивание интерфейсов
-embed string             # Список интерфейсов для встраивания (через запятую)
-embed-pkg string         # Пакеты для встраиваемых интерфейсов (через запятую)

# Настройка конструктора
-use-constructor bool     # Использовать конструктор исходной структуры (по умолчанию: да)
-constructor-name string  # Название конструктора исходной структуры (по умолчанию: NewИмяСтруктурыСБольшойБуквы)

# Отладка
-verbose                  # Подробный вывод
```

## Примеры использования

### 1. Простая структура в том же пакете

```bash
go run ./tools/generate_tracing.go -struct=UserService -file=service.go
```

### 2. gRPC контроллер с named импортами

```bash
go run ./tools/gen_tracing.go \
  -struct=authController \
  -interface=AuthServiceServer \
  -interface-pkg=github.com/intezya/auth_service/protos/go/auth \
  -file=./internal/adapters/grpc/controller.go \
  -output=./internal/adapters/grpc/controller_tracing.go \
  -embed=UnimplementedAuthServiceServer
```

### 3. Структура из другого пакета

```bash
go run ./tools/generate_tracing.go \
  -struct=DatabaseRepo \
  -interface=Repository \
  -struct-pkg=github.com/example/internal/repo \
  -interface-pkg=github.com/example/internal/domain \
  -output=repo_tracing.go
```

### 4. Поиск структуры в директории

```bash
go run ./tools/generate_tracing.go \
  -struct=PaymentService \
  -search-dir=./internal/services \
  -recursive
```

### 5. Множественные встраиваемые интерфейсы

```bash
go run ./tools/generate_tracing.go \
  -struct=myGrpcServer \
  -interface=MyServiceServer \
  -interface-pkg=github.com/example/proto \
  -embed="UnimplementedMyServiceServer,ValidationInterface" \
  -embed-pkg="github.com/example/proto,github.com/example/validation"
```

(также посмотрите taskfile.yaml в корне репозитория)

## Специальные возможности для gRPC

### Автоматическое встраивание mustEmbed интерфейсов

Генератор автоматически обнаруживает и встраивает `mustEmbed` интерфейсы из gRPC сервисов:

```go
// Исходная структура
type authController struct {
    authpb.UnimplementedAuthServiceServer  // есть изначально
    authService service.AuthService
}

// Сгенерированная обёртка
type authControllerWithTracing struct {
    wrapped authpb.AuthServiceServer
    authpb.UnimplementedAuthServiceServer  // автоматически встроен
}
```

### Корректная обработка named импортов

Генератор правильно обрабатывает named импорты:

```go
// Входной файл
import (
    authpb "github.com/intezya/auth_service/protos/go/auth"
)

// Сгенерированный файл
import (
    authpb "github.com/intezya/auth_service/protos/go/auth"  // сохраняет alias
)

func (t *authControllerWithTracing) Register(
    ctx context.Context, 
    request *authpb.AuthenticationRequest,  // использует правильный alias
) (*authpb.Empty, error) {
    // ...
}
```

## Структура сгенерированного кода

Генератор создаёт:

1. **Обёрточную структуру** с полем `wrapped` типа интерфейса
2. **Конструктор** `New<Struct>WithTracing`
3. **Методы-обёртки** для всех публичных методов оригинальной структуры
4. **Трейсинг** для методов с `context.Context` в качестве первого параметра
5. **Встроенные интерфейсы** при необходимости

### Пример сгенерированного кода

```go
// Code generated by tracing-gen. DO NOT EDIT.

package grpc

import (
    "context"
    tracer "github.com/intezya/auth_service/internal/infrastructure/metrics/tracer"
    authpb "github.com/intezya/auth_service/protos/go/auth"
)

type authControllerWithTracing struct {
    wrapped authpb.AuthServiceServer
    authpb.UnimplementedAuthServiceServer
}

func NewAuthControllerWithTracing(wrapped authpb.AuthServiceServer) authpb.AuthServiceServer {
    return &authControllerWithTracing{
        wrapped: wrapped,
    }
}

func (t *authControllerWithTracing) Register(ctx context.Context, request *authpb.AuthenticationRequest) (*authpb.Empty, error) {
    ctx, span := tracer.StartSpan(ctx, "AuthController.Register")
    defer span.End()
    return t.wrapped.Register(ctx, request)
}
```

## Алгоритм работы

1. **Парсинг входного файла** - анализ AST для поиска структуры и её методов
2. **Анализ импортов** - сохранение named импортов и их псевдонимов
3. **Поиск интерфейсов** - определение целевого интерфейса и встраиваемых
4. **Генерация кода** - создание обёрточной структуры с трейсингом
5. **Оптимизация импортов** - удаление неиспользуемых импортов
6. **Форматирование** - применение `go fmt` к результату

## Ограничения

- Обрабатывает только публичные методы (начинающиеся с заглавной буквы)
- Трейсинг добавляется только к методам с `context.Context` первым параметром
- Не поддерживает generic типы (Go 1.18+)
- Требует, чтобы структура реализовывала указанный интерфейс

## Интеграция

```bash
# В Taskfile
version: '3'

tasks:
  generate_tracing:
    desc: "Generate tracing wrappers for services"
    cmds:
      - |
        go run ./tools/generate_tracing.go \
          --struct=authService \
          --interface=AuthService \
          --file=./internal/application/service/auth_service.go \
          --output=./internal/application/service/auth_service_tracing.go
    
# В go:generate директивах
//go:generate go run ./tools/generate_tracing.go -struct=UserService -file=service.go ...
```
