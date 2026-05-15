# Uber Go Style Guide — Покрытие линтерами

Полный анализ всех рекомендаций из [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)
с указанием наличия автоматического линтера в `golangci-lint` или других инструментах.

---

## Guidelines

### Структуры данных и интерфейсы

| # | Рекомендация | Линтер / Инструмент | Покрытие |
|---|---|---|---|
| 1 | **Pointers to Interfaces** — не используй указатели на интерфейсы | `uberlint: ifaceptr` | Полное |
| 2 | **Verify Interface Compliance** — проверяй `var _ Interface = (*Type)(nil)` | Нет | — |
| 3 | **Receivers and Interfaces** — value vs pointer receivers | `recvcheck` (частично, проверяет консистентность receiver types) | Частичное |
| 4 | **Zero-value Mutexes are Valid** — не используй `new(sync.Mutex)`, не встраивай mutex | `govet: copylocks` (ловит копирование mutex); нет линтера для запрета встраивания | Частичное |
| 5 | **Copy Slices and Maps at Boundaries** — копируй при получении/возврате | Нет | — |
| 6 | **Defer to Clean Up** — используй defer для освобождения ресурсов | Нет | — |
| 7 | **Channel Size is One or None** — каналы с буфером 0 или 1 | `uberlint: chansize` | Полное |
| 8 | **Start Enums at One** — начинай enum с `iota + 1` | `uberlint: enumstart` | Полное |

---

### Работа со временем

| # | Рекомендация | Линтер / Инструмент | Покрытие |
|---|---|---|---|
| 9 | **Use `time.Time` for instants of time** | Нет | — |
| 10 | **Use `time.Duration` for periods of time** | Нет | — |
| 11 | **Use `time.Time` and `time.Duration` with external systems** — включай единицы в имена полей | `uberlint: timefield` | Частичное |

---

### Обработка ошибок

| # | Рекомендация | Линтер / Инструмент | Покрытие |
|---|---|---|---|
| 12 | **Error Types** — выбирай правильный тип ошибки (errors.New / fmt.Errorf / custom) | Нет (паттерн проектирования) | — |
| 13 | **Error Wrapping** — используй `%w` для оборачивания ошибок | `wrapcheck` (проверяет что ошибки из внешних пакетов обёрнуты); `errorlint` (проверяет правильное использование `%w`, `errors.Is`/`errors.As`) | Полное |
| 14 | **Error Naming** — `Err` prefix для vars, `Error` suffix для типов | `errname` (проверяет что sentinel errors имеют префикс `Err`, а типы ошибок — суффикс `Error`) | Полное |
| 15 | **Handle Errors Once** — не логируй и возвращай одновременно | Нет | — |
| 16 | **Handle Type Assertion Failures** — всегда используй `comma-ok` | `forcetypeassert` (флагает type assertions без comma-ok) | Полное |

---

### Управление потоком

| # | Рекомендация | Линтер / Инструмент | Покрытие |
|---|---|---|---|
| 17 | **Don't Panic** — избегай panic в продакшн-коде | `uberlint: nopanic` | Полное |
| 18 | **Use `sync/atomic` typed values** — используй `atomic.Bool`, `atomic.Int64` и т.д. (с Go 1.19 в stdlib, вместо `go.uber.org/atomic`) | `uberlint: atomicstd` | Полное |
| 19 | **Avoid Mutable Globals** — избегай изменяемых глобальных переменных | `gochecknoglobals` (частично, флагает любые глобальные vars, не только mutable) | Частичное |
| 20 | **Avoid Embedding Types in Public Structs** — не встраивай типы в экспортируемые структуры | `uberlint: publicembed` | Полное |
| 21 | **Avoid Using Built-In Names** — не затеняй предопределённые идентификаторы | `predeclared` (флагает затенение predeclared identifiers); `govet: shadow` | Полное |
| 22 | **Avoid `init()`** — избегай init-функций | `gochecknoinits` (флагает наличие init()); revive: `add-constant` правило `gochecknoinits` | Полное |
| 23 | **Exit in Main** — os.Exit/log.Fatal только в main() | `revive: deep-exit` (флагает os.Exit/log.Fatal вне main/init) | Полное |
| 24 | **Exit Once** — только один вызов os.Exit в main | `revive: deep-exit` (частично) | Частичное |

---

### Сериализация

| # | Рекомендация | Линтер / Инструмент | Покрытие |
|---|---|---|---|
| 25 | **Use field tags in marshaled structs** — добавляй теги для JSON/YAML полей | `musttag` (проверяет наличие тегов у экспортируемых полей структур, которые маршалятся) | Полное |

---

### Горутины

| # | Рекомендация | Линтер / Инструмент | Покрытие |
|---|---|---|---|
| 26 | **Don't fire-and-forget goroutines** — не запускай goroutines без управления жизненным циклом | Нет (`go.uber.org/goleak` — тестовый инструмент, не линтер) | — |
| 27 | **Wait for goroutines to exit** — используй WaitGroup или done-channel | Нет | — |
| 28 | **No goroutines in `init()`** — не запускай горутины в init | Частично `gochecknoinits` (если init запрещён, то и goroutine в нём тоже) | Частичное |

---

### Performance

| # | Рекомендация | Линтер / Инструмент | Покрытие |
|---|---|---|---|
| 29 | **Prefer strconv over fmt** — strconv быстрее fmt | `perfsprint` (флагает fmt.Sprint/Sprintf когда можно использовать strconv) | Полное |
| 30 | **Avoid repeated string-to-byte conversions** — не конвертируй `[]byte("str")` в цикле | `uberlint: stringbytes` | Полное |
| 31 | **Prefer Specifying Container Capacity** — указывай capacity для slices и maps | `prealloc` (предлагает prealloc для slices перед append); `makezero` (флагает slices без начальной capacity) | Частичное |

---

### Style

| # | Рекомендация | Линтер / Инструмент | Покрытие |
|---|---|---|---|
| 32 | **Avoid overly long lines** — мягкий лимит 99 символов | `lll` (настраиваемый лимит длины строки) | Полное |
| 33 | **Be Consistent** — будь последовательным | Нет (субъективно) | — |
| 34 | **Group Similar Declarations** — группируй const/var/type/import | `grouper` (анализирует группировку const/var/type/import declarations) | Полное |
| 35 | **Import Group Ordering** — stdlib отдельно от сторонних | `goimports` (автоматически форматирует группы импортов) | Полное |
| 36 | **Package Names** — lowercase, не plural, не util/common | `revive: var-naming` (частично, проверяет naming conventions) | Частичное |
| 37 | **Function Names** — MixedCaps, underscores только в тестах | `revive: var-naming` (частично) | Частичное |
| 38 | **Import Aliasing** — alias только при конфликтах | `importas` (настраиваемые правила для import aliasing) | Частичное |
| 39 | **Function Grouping and Ordering** — группировка по receiver, порядок вызовов | Нет | — |
| 40 | **Reduce Nesting** — уменьшай вложенность, early return | `revive: early-return` (предлагает ранний возврат); `nestif` (флагает глубокую вложенность if) | Частичное |
| 41 | **Unnecessary Else** — убирай лишний else | `revive: early-return` (предлагает убрать else через early return) | Полное |
| 42 | **Top-level Variable Declarations** — не указывай тип в var, если он совпадает | `uberlint: vartype` | Полное |
| 43 | **Prefix Unexported Globals with _** — `_defaultPort` вместо `defaultPort` | `uberlint: globalprefix` | Полное |
| 44 | **Embedding in Structs** — embedded types сверху, пустая строка-разделитель | `uberlint: embedlayout` | Частичное |
| 45 | **Local Variable Declarations** — предпочитай `:=` вместо `var s = "foo"` | `uberlint: localvar` | Полное |
| 46 | **nil is a valid slice** — возвращай nil вместо `[]int{}` | `uberlint: nilslice` | Полное |
| 47 | **Reduce Scope of Variables** — сокращай область видимости | Нет | — |
| 48 | **Avoid Naked Parameters** — добавляй комментарии `/* ... */` для неочевидных параметров | `uberlint: nakedparams` | Частичное |
| 49 | **Use Raw String Literals to Avoid Escaping** — используй backticks | `uberlint: rawstring` | Полное |

---

### Initializing Structs

| # | Рекомендация | Линтер / Инструмент | Покрытие |
|---|---|---|---|
| 50 | **Use Field Names to Initialize Structs** — всегда указывай имена полей | `govet` (проверяет unkeyed composite literals, сам гайд говорит "enforced by go vet") | Частичное (для внешних пакетов) |
| 51 | **Omit Zero Value Fields in Structs** — опускай нулевые поля | `uberlint: zerofields` | Полное |
| 52 | **Use `var` for Zero Value Structs** — `var user User` вместо `user := User{}` | `uberlint: zerovar` | Полное |
| 53 | **Initializing Struct References** — `&T{}` вместо `new(T)` | `uberlint: newref` | Полное |
| 54 | **Initializing Maps** — `make()` для пустых, литералы для фиксированных | `uberlint: mapinit` | Частичное |

---

### Printf

| # | Рекомендация | Линтер / Инструмент | Покрытие |
|---|---|---|---|
| 55 | **Format Strings outside Printf** — делай format strings `const` | `uberlint: constprintf` | Полное |
| 56 | **Naming Printf-style Functions** — функции Printf-стиля должны заканчиваться на `f` | `govet: printf` (проверяет format strings для известных Printf-функций, кастомные через `-printfuncs`) | Полное |

---

### Patterns

| # | Рекомендация | Линтер / Инструмент | Покрытие |
|---|---|---|---|
| 57 | **Test Tables** — используй table-driven tests | Нет (паттерн проектирования) | — |
| 58 | **Avoid Unnecessary Complexity in Table Tests** — не усложняй тест-таблицы | Нет | — |
| 59 | **Parallel Tests** — `tt := tt` для `t.Parallel()` | `tparallel` (частично, проверяет использование t.Parallel); `paralleltest` (проверяет что t.Parallel вызывается) | Частичное |
| 60 | **Functional Options** — паттерн Option для опциональных аргументов | Нет (паттерн проектирования) | — |

---

## Сводка

| Категория | Всего правил | Есть линтер | Частично | Нет линтера |
|---|---|---|---|---|
| Структуры данных и интерфейсы | 8 | 3 | 2 | 3 |
| Работа со временем | 3 | 0 | 1 | 2 |
| Обработка ошибок | 5 | 3 | 0 | 2 |
| Управление потоком | 8 | 6 | 2 | 0 |
| Сериализация | 1 | 1 | 0 | 0 |
| Горутины | 3 | 0 | 1 | 2 |
| Performance | 3 | 2 | 1 | 0 |
| Style | 18 | 9 | 6 | 3 |
| Initializing Structs | 5 | 3 | 2 | 0 |
| Printf | 2 | 2 | 0 | 0 |
| Patterns | 4 | 0 | 1 | 3 |
| **ИТОГО** | **60** | **29** | **16** | **15** |

---

## Рекомендуемая конфигурация golangci-lint

Для максимального покрытия Uber Go Style Guide:

```yaml
version: "2"

linters:
  enable:
    # Базовые (рекомендованы самим Uber)
    - errcheck
    - goimports
    - revive
    - govet
    - staticcheck

    # uberlint — custom module plugin with all uberlint analyzers.
    - uberlint

    # Дополнительные для покрытия Uber Style Guide
    - errname          # Error Naming (#14)
    - errorlint        # Error Wrapping (#13)
    - wrapcheck        # Error Wrapping — external packages (#13)
    - forcetypeassert  # Handle Type Assertion Failures (#16)
    - predeclared      # Avoid Using Built-In Names (#21)
    - gochecknoinits   # Avoid init() (#22)
    - musttag          # Use field tags in marshaled structs (#25)
    - perfsprint       # Prefer strconv over fmt (#29)
    - prealloc         # Prefer Specifying Container Capacity (#31)
    - makezero         # Prefer Specifying Container Capacity (#31)
    - lll              # Avoid overly long lines (#32)
    - grouper          # Group Similar Declarations (#34)
    - importas         # Import Aliasing (#38)
    - nestif           # Reduce Nesting (#40)
  settings:
    custom:
      uberlint:
        type: "module"
        description: "Uber Go Style Guide linter"
    lll:
      line-length: 99
    revive:
      rules:
        - name: deep-exit       # Exit in Main (#23)
        - name: early-return    # Reduce Nesting + Unnecessary Else (#40, #41)
    govet:
      enable:
        - shadow                # Avoid Using Built-In Names (#21)
        - copylocks             # Zero-value Mutexes (#4)
```

---

## Правила покрытые uberlint (20 реализовано)

| # | Правило | Analyzer ID |
|---|---|---|
| 1 | Pointers to Interfaces | `ifaceptr` |
| 7 | Channel Size is One or None | `chansize` |
| 8 | Start Enums at One | `enumstart` |
| 17 | Don't Panic | `nopanic` |
| 18 | Use sync/atomic typed values | `atomicstd` |
| 30 | Avoid repeated string-to-byte conversions | `stringbytes` |
| 42 | Top-level Variable Declarations | `vartype` |
| 43 | Prefix Unexported Globals with _ | `globalprefix` |
| 46 | nil is a valid slice | `nilslice` |
| 49 | Use Raw String Literals to Avoid Escaping | `rawstring` |
| 52 | Use var for Zero Value Structs | `zerovar` |
| 53 | Initializing Struct References | `newref` |
| 11 | time.Time/Duration with external systems | `timefield` |
| 20 | Avoid Embedding Types in Public Structs | `publicembed` |
| 44 | Embedding in Structs | `embedlayout` |
| 45 | Local Variable Declarations | `localvar` |
| 48 | Avoid Naked Parameters | `nakedparams` |
| 51 | Omit Zero Value Fields in Structs | `zerofields` |
| 54 | Initializing Maps | `mapinit` |
| 55 | Format Strings outside Printf | `constprintf` |

---

## Правила БЕЗ линтера (ниша для uberlint Phase 3+)

Следующие 15 правил не имеют автоматического линтера — кандидаты для следующей фазы:

1. Verify Interface Compliance
2. Copy Slices and Maps at Boundaries
3. Defer to Clean Up
4. Use time.Time for instants
5. Use time.Duration for periods
6. Error Types (proper pattern choice)
7. Handle Errors Once
8. Don't fire-and-forget goroutines
9. Wait for goroutines to exit
10. Function Grouping and Ordering
11. Reduce Scope of Variables
12. Test Tables
13. Avoid Unnecessary Complexity in Table Tests
14. Functional Options
15. Be Consistent (subjective)
