# Принципи програмування у проєкті «Місто Житомир»

Проєкт — веб-сайт про місто Житомир, реалізований мовою Go з використанням стандартної бібліотеки `net/http` та шаблонізатора `html/template`.

---

## 1. DRY (Don't Repeat Yourself)

Принцип «не повторюйся» реалізований через спільний layout-шаблон.

**Приклад:** замість копіювання однакового HTML (header, nav, footer) у кожну сторінку — використовується один файл `templates/layout.html`, а кожна сторінка визначає лише блок `{{ define "content" }}`.

```
// main.go — функція parseTemplate завантажує layout один раз для кожної сторінки
func parseTemplate(page string) *template.Template {
    tmpl := template.Must(template.ParseFiles(
        "templates/layout.html",
        "templates/"+page+".html",
    ))
    return tmpl
}
```

> **Файл:** [`main.go`, рядки 74–80](main.go#L74-80) — функція `parseTemplate`  
> **Файл:** [`templates/layout.html`](templates/layout.html) — єдиний layout для всього сайту

---

## 2. SRP — Single Responsibility Principle (Принцип єдиної відповідальності)

Кожна функція-обробник відповідає лише за одну сторінку. Жодна функція не змішує логіку кількох сторінок.

```go
// Кожен handler відповідає лише за свою сторінку:
func homeHandler(w http.ResponseWriter, r *http.Request) { ... }
func historyHandler(w http.ResponseWriter, r *http.Request) { ... }
func attractionsHandler(w http.ResponseWriter, r *http.Request) { ... }
func factsHandler(w http.ResponseWriter, r *http.Request) { ... }
func contactHandler(w http.ResponseWriter, r *http.Request) { ... }
```

> **Файл:** [`main.go`, рядки 83–220](main.go)

---

## 3. Separation of Concerns (Розділення відповідальностей)

Дані, логіка і представлення чітко розділені:

- **Дані** — структури (`HomeData`, `HistoryData`, `AttractionsData`, `FactsData`, `ContactData`) у `main.go`
- **Логіка** — обробники HTTP-запитів у `main.go`
- **Представлення** — HTML-шаблони у директорії `templates/`

```go
// Структури даних відокремлені від логіки обробників
type HomeData struct {
    PageData
    Famous []FamousPerson
}

type HistoryData struct {
    PageData
    Events []HistoryEvent
}
```

> **Файл:** [`main.go`, рядки 14–70](main.go) — всі типи даних  
> **Файл:** [`templates/home.html`](templates/home.html) — шаблон представлення

---

## 4. Composition over Inheritance (Композиція замість наслідування)

Go не має класичного наслідування. Замість нього використовується вбудовування (embedding) структур — `PageData` вбудовується у всі специфічні структури сторінок:

```go
type PageData struct {
    Title string
    Page  string
}

type HomeData struct {
    PageData          // вбудовування — аналог наслідування
    Famous []FamousPerson
}

type HistoryData struct {
    PageData
    Events []HistoryEvent
}
```

> **Файл:** [`main.go`, рядки 16–55](main.go)

---

## 5. KISS (Keep It Simple, Stupid)

Проєкт не використовує зовнішніх фреймворків чи залежностей — лише стандартна бібліотека Go. Це спрощує структуру, деплой і розуміння коду.

```go
import (
    "fmt"
    "html/template"
    "net/http"
    "strings"
)
```

> **Файл:** [`main.go`, рядки 1–8](main.go)

---

## 6. Fail Fast (Швидка відмова)

При помилці завантаження шаблону програма одразу завершується завдяки `template.Must`, а не продовжує роботу з невалідним станом:

```go
func parseTemplate(page string) *template.Template {
    tmpl := template.Must(template.ParseFiles(
        "templates/layout.html",
        "templates/"+page+".html",
    ))
    return tmpl
}
```

> **Файл:** [`main.go`, рядки 74–80](main.go)

---

## 7. Input Validation (Валідація вхідних даних)

У обробнику форми контактів реалізована послідовна валідація всіх обов'язкових полів перед обробкою даних:

```go
if strings.TrimSpace(data.FormName) == "" {
    data.HasError = true
    data.ErrorMsg = "Будь ласка, введіть ваше ім'я!"
} else if strings.TrimSpace(data.FormEmail) == "" {
    data.HasError = true
    data.ErrorMsg = "Будь ласка, введіть email!"
} else if !strings.Contains(data.FormEmail, "@") {
    data.HasError = true
    data.ErrorMsg = "Введіть коректний email (має містити @)!"
} else if strings.TrimSpace(data.FormMessage) == "" {
    data.HasError = true
    data.ErrorMsg = "Будь ласка, введіть повідомлення!"
}
```

> **Файл:** [`main.go`, рядки 195–208](main.go)

---

## 8. Convention over Configuration

Назви файлів шаблонів відповідають назвам маршрутів і сторінок за конвенцією, що усуває потребу в явній конфігурації:

- маршрут `/history` → обробник `historyHandler` → шаблон `templates/history.html`
- маршрут `/attractions` → обробник `attractionsHandler` → шаблон `templates/attractions.html`

```go
func parseTemplate(page string) *template.Template {
    tmpl := template.Must(template.ParseFiles(
        "templates/layout.html",
        "templates/"+page+".html",   // ім'я файлу = ім'я сторінки
    ))
    return tmpl
}
```

> **Файл:** [`main.go`, рядки 74–80](main.go)
