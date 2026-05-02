# Метрика Чепина — `ValidateAccountNumber`

## Описание метрики

Метрика Чепина (Chapin's metric) оценивает качество использования переменных в программном модуле. Каждая переменная относится к одной из четырёх категорий:

| Категория | Обозначение | Описание |
|---|---|---|
| Input | **P** | Переменные, которые только читаются (входные параметры, внешние константы, возвращаемые значения) |
| Modified | **M** | Переменные, которые изменяются (записываются) внутри функции |
| Control | **C** | Переменные, используемые в условиях ветвления или управления потоком |
| Spurious | **T** | Паразитные переменные, не влияющие на логику |

**Формула:**

```
Q = P + 2M + 3C + 0.5T
```

---

## Анализируемая функция

**Файл:** `service/subscriber/validation.go`
**Функция:** `ValidateAccountNumber(accountNumber string) error`

```go
func ValidateAccountNumber(accountNumber string) error {
    runes := []rune(accountNumber)
    length := len(runes)

    if length < accountNumberMinLength {
        return ErrAccountNumberTooShort
    }
    if length > accountNumberMaxLength {
        return ErrAccountNumberTooLong
    }

    if !isUpperRussianLetter(runes[0]) {
        return ErrAccountNumberStartNotLetter
    }

    if !isDigit(runes[length-1]) {
        return ErrAccountNumberEndNotDigit
    }

    hasHyphen := false
    prevHyphen := false

    for _, r := range runes {
        switch {
        case isDigit(r):
            prevHyphen = false
        case isUpperRussianLetter(r):
            prevHyphen = false
        case r == '-':
            if prevHyphen {
                return ErrAccountNumberConsecutiveHyphens
            }
            hasHyphen = true
            prevHyphen = true
        default:
            return ErrAccountNumberInvalidChar
        }
    }

    if !hasHyphen {
        return ErrAccountNumberNoHyphen
    }

    return nil
}
```

---

## Классификация переменных

### P — входные / только читаемые (P = 8)

| Переменная | Тип | Обоснование |
|---|---|---|
| `accountNumber` | `string` | Входной параметр, только читается |
| `ErrAccountNumberTooShort` | `error` | Только возвращается, не изменяется |
| `ErrAccountNumberTooLong` | `error` | Только возвращается, не изменяется |
| `ErrAccountNumberStartNotLetter` | `error` | Только возвращается, не изменяется |
| `ErrAccountNumberEndNotDigit` | `error` | Только возвращается, не изменяется |
| `ErrAccountNumberConsecutiveHyphens` | `error` | Только возвращается, не изменяется |
| `ErrAccountNumberInvalidChar` | `error` | Только возвращается, не изменяется |
| `ErrAccountNumberNoHyphen` | `error` | Только возвращается, не изменяется |

### M — модифицируемые (M = 1)

| Переменная | Тип | Обоснование |
|---|---|---|
| `runes` | `[]rune` | Присваивается однократно, не участвует напрямую в условиях ветвления |

### C — управляющие (C = 6)

| Переменная | Тип | Где используется в условии |
|---|---|---|
| `length` | `int` | `if length < accountNumberMinLength`, `if length > accountNumberMaxLength`, `runes[length-1]` |
| `hasHyphen` | `bool` | `if !hasHyphen` |
| `prevHyphen` | `bool` | `if prevHyphen` |
| `r` | `rune` | `case isDigit(r)`, `case isUpperRussianLetter(r)`, `case r == '-'` |
| `accountNumberMinLength` | `const int` | `if length < accountNumberMinLength` |
| `accountNumberMaxLength` | `const int` | `if length > accountNumberMaxLength` |

### T — паразитные (T = 0)

Паразитных переменных нет.

---

## Вычисление

```
Q = P + 2M + 3C + 0.5T
Q = 8 + 2×1 + 3×6 + 0.5×0
Q = 8 + 2 + 18 + 0
Q = 28
```

---

## Результат

| Категория | Количество | Вклад |
|---|---|---|
| P (входные) | 8 | 8 |
| M (модифицируемые) | 1 | 2 |
| C (управляющие) | 6 | 18 |
| T (паразитные) | 0 | 0 |
| **Q** | | **28** |

---

## Интерпретация

| Диапазон Q | Оценка сложности |
|---|---|
| 0–10 | Низкая сложность |
| 11–20 | Средняя сложность |
| **21–30** | **Умеренно высокая сложность** |
| > 30 | Высокая сложность |

**Вывод:** Значение Q = 28 соответствует умеренно высокой сложности. Основной вклад вносят управляющие переменные (18 из 28), что отражает разветвлённую логику валидации: 5 независимых условий возврата ошибки и цикл с переключением по типу символа. Семь переменных ошибок в категории P увеличивают метрику, однако не усложняют логику управления потоком.
