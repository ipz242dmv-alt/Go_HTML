package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

// ══════════════════════════════════════════════════════════════════════════════
//  Типи даних для шаблонів
// ══════════════════════════════════════════════════════════════════════════════

type PageData struct {
	Title string
	Page  string
}
type FamousPerson struct {
	Name     string
	Activity string
	Years    string
}
type HomeData struct {
	PageData
	Famous []FamousPerson
}
type HistoryEvent struct {
	Year  string
	Event string
}
type HistoryData struct {
	PageData
	Events []HistoryEvent
}
type Attraction struct {
	Icon        string
	Name        string
	Description string
	Address     string
	Color       string
}
type AttractionsData struct {
	PageData
	Attractions []Attraction
}
type Stat struct {
	Value string
	Label string
}
type FactsData struct {
	PageData
	Stats []Stat
}
type ContactData struct {
	PageData
	IsPost      bool
	HasError    bool
	ErrorMsg    string
	FormName    string
	FormEmail   string
	FormSubject string
	FormMessage string
}

// ══════════════════════════════════════════════════════════════════════════════
//  Кешування шаблонів (Fix Issue #1)
// ══════════════════════════════════════════════════════════════════════════════

var templates map[string]*template.Template

// initTemplates парсить усі шаблони один раз при старті програми
func initTemplates() {
	pages := []string{"home", "history", "attractions", "facts", "contact"}
	templates = make(map[string]*template.Template, len(pages))
	for _, page := range pages {
		templates[page] = template.Must(template.ParseFiles(
			"templates/layout.html",
			"templates/"+page+".html",
		))
	}
}

// render виконує шаблон для заданої сторінки
func render(w http.ResponseWriter, page string, data any) {
	tmpl, ok := templates[page]
	if !ok {
		http.Error(w, "Сторінку не знайдено", http.StatusNotFound)
		return
	}
	if err := tmpl.ExecuteTemplate(w, "layout.html", data); err != nil {
		http.Error(w, "Помилка рендерингу", http.StatusInternalServerError)
	}
}

// ══════════════════════════════════════════════════════════════════════════════
//  Обробники сторінок
// ══════════════════════════════════════════════════════════════════════════════

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	data := HomeData{
		PageData: PageData{Title: "Головна", Page: "home"},
		Famous: []FamousPerson{
			{"Сергій Корольов", "Конструктор ракетної техніки, засновник практичної космонавтики", "1907–1966"},
			{"Кароль Шимановський", "Польський композитор і піаніст", "1882–1937"},
			{"Борис Тен (Микола Хомичевський)", "Поет, перекладач, перекладач «Іліади»", "1897–1983"},
			{"Іван Огієнко (митрополит Іларіон)", "Мовознавець, церковний діяч, перекладач Біблії", "1882–1972"},
			{"Святослав Ріхтер", "Піаніст, один з найвидатніших у XX столітті", "1915–1997"},
			{"Микола Щорс", "Радянський військовий діяч, командарм", "1895–1919"},
		},
	}
	render(w, "home", data)
}

func historyHandler(w http.ResponseWriter, r *http.Request) {
	data := HistoryData{
		PageData: PageData{Title: "Історія", Page: "history"},
		Events: []HistoryEvent{
			{"884", "Заснування Житомира — за легендою, дружинником Аскольда і Діра на ім'я Житомир"},
			{"1240", "Перша літописна згадка — місто спустошене військами хана Батия"},
			{"1320", "Житомир захоплений литовським князем Гедиміном, входить до складу Великого князівства Литовського"},
			{"1444", "Місто отримує самоврядування за Магдебурзьким правом"},
			{"1569", "Після Люблінської унії Житомир входить до складу Речі Посполитої"},
			{"1648", "Під час Визвольної війни Хмельницького місто неодноразово переходить з рук у руки"},
			{"1667", "За Андрусівським перемир'ям Житомир відходить до Польщі"},
			{"1778", "Місто стає центром Волинського воєводства"},
			{"1795", "Після третього поділу Польщі Житомир входить до складу Російської імперії"},
			{"1804", "Житомир стає губернським містом Волинської губернії"},
			{"1917–1920", "Місто переходить під контроль різних сил у роки Громадянської війни"},
			{"1941–1943", "Окупація нацистськими військами у роки Другої світової війни"},
			{"1991", "Незалежність України — Житомир стає обласним центром"},
			{"2022", "Місто витримує ракетні удари в ході повномасштабного вторгнення Росії"},
		},
	}
	render(w, "history", data)
}

func attractionsHandler(w http.ResponseWriter, r *http.Request) {
	data := AttractionsData{
		PageData: PageData{Title: "Пам'ятки", Page: "attractions"},
		Attractions: []Attraction{
			{"⛪", "Кафедральний собор святого Михайла", "Одна з найстаріших православних споруд міста, збудована у XVIII столітті. Є духовним і архітектурним символом Житомира.", "майдан Соборний", ""},
			{"🕍", "Костел святої Софії", "Барочний римо-католицький костел, збудований у 1737–1751 роках. Один з найкрасивіших архітектурних пам'ятників міста.", "вул. Михайлівська, 10", "green"},
			{"🏛️", "Краєзнавчий музей", "Один з найстаріших музеїв України (1865 р.). Зберігає понад 200 000 експонатів з природи, археології та етнографії Житомирщини.", "вул. Михайлівська, 1", "orange"},
			{"🚀", "Музей космонавтики ім. С. П. Корольова", "Присвячений уродженцю Житомира — видатному конструктору Сергію Корольову. Містить унікальні експонати ракетно-космічної техніки.", "вул. Дмитрівська, 5", "purple"},
			{"🌉", "Замковий міст і скелі на Тетереві", "Мальовниче місце над річкою Тетерів із гранітними скелями і видом на старе місто. Популярне місце прогулянок та фотосесій.", "набережна Тетереву", "teal"},
			{"🎭", "Обласний музичний театр", "Провідний театральний заклад Житомирщини. Репертуар включає опери, оперети, мюзикли та балетні вистави.", "майдан Соборний, 2", ""},
		},
	}
	render(w, "attractions", data)
}

func factsHandler(w http.ResponseWriter, r *http.Request) {
	data := FactsData{
		PageData: PageData{Title: "Факти", Page: "facts"},
		Stats: []Stat{
			{"884", "рік заснування міста"},
			{"~263 000", "мешканців (2023 р.)"},
			{"61 км²", "площа міста"},
			{"10+", "закладів вищої освіти"},
			{"5", "театрів та музеїв"},
			{"140 км", "до Києва"},
			{"1240", "рік першої літописної згадки"},
			{"3", "річки протікають через місто"},
		},
	}
	render(w, "facts", data)
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	data := ContactData{
		PageData: PageData{Title: "Контакти", Page: "contact"},
	}
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "<p>Помилка: %v</p>", err)
			return
		}
		data.IsPost = true
		data.FormName = r.PostForm.Get("name")
		data.FormEmail = r.PostForm.Get("email")
		data.FormSubject = r.PostForm.Get("subject")
		data.FormMessage = r.PostForm.Get("message")
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
	}
	render(w, "contact", data)
}

// ══════════════════════════════════════════════════════════════════════════════
//  main
// ══════════════════════════════════════════════════════════════════════════════

func main() {
	initTemplates() // ← шаблони парсяться один раз при старті

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/history", historyHandler)
	http.HandleFunc("/attractions", attractionsHandler)
	http.HandleFunc("/facts", factsHandler)
	http.HandleFunc("/contact", contactHandler)

	addr := "localhost:8080"
	fmt.Println("╔═══════════════════════════════════════╗")
	fmt.Println("║   Сайт «Місто Житомир» запущено!      ║")
	fmt.Println("╠═══════════════════════════════════════╣")
	fmt.Printf("║   Адреса: http://%s       ║\n", addr)
	fmt.Println("║   Сторінки:                           ║")
	fmt.Println("║     /            — Головна            ║")
	fmt.Println("║     /history     — Історія            ║")
	fmt.Println("║     /attractions — Пам'ятки           ║")
	fmt.Println("║     /facts       — Факти              ║")
	fmt.Println("║     /contact     — Контакти           ║")
	fmt.Println("║   Ctrl+C для зупинки                  ║")
	fmt.Println("╚═══════════════════════════════════════╝")
	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Println("Помилка сервера:", err)
	}
}
