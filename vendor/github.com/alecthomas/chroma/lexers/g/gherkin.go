package g

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Gherkin lexer.
var Gherkin = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Gherkin",
		Aliases:   []string{"cucumber", "Cucumber", "gherkin", "Gherkin"},
		Filenames: []string{"*.feature", "*.FEATURE"},
		MimeTypes: []string{"text/x-gherkin"},
	},
	gherkinRules,
))

func gherkinRules() Rules {
	stepKeywords := `^(\s*)(하지만|조건|먼저|만일|만약|단|그리고|그러면|那麼|那么|而且|當|当|前提|假設|假设|假如|假定|但是|但し|並且|并且|同時|同时|もし|ならば|ただし|しかし|かつ|و |متى |لكن |عندما |ثم |بفرض |اذاً |כאשר |וגם |בהינתן |אזי |אז |אבל |Якщо |Унда |То |Припустимо, що |Припустимо |Онда |Но |Нехай |Лекин |Когато |Када |Кад |К тому же |И |Задато |Задати |Задате |Если |Допустим |Дадено |Ва |Бирок |Аммо |Али |Але |Агар |А |І |Și |És |Zatati |Zakładając |Zadato |Zadate |Zadano |Zadani |Zadan |Youse know when youse got |Youse know like when |Yna |Ya know how |Ya gotta |Y |Wun |Wtedy |When y'all |When |Wenn |WEN |Và |Ve |Und |Un |Thì |Then y'all |Then |Tapi |Tak |Tada |Tad |Så |Stel |Soit |Siis |Si |Sed |Se |Quando |Quand |Quan |Pryd |Pokud |Pokiaľ |Però |Pero |Pak |Oraz |Onda |Ond |Oletetaan |Og |Och |O zaman |Når |När |Niin |Nhưng |N |Mutta |Men |Mas |Maka |Majd |Mais |Maar |Ma |Lorsque |Lorsqu'|Kun |Kuid |Kui |Khi |Keď |Ketika |Když |Kaj |Kai |Kada |Kad |Jeżeli |Ja |Ir |I CAN HAZ |I |Ha |Givun |Givet |Given y'all |Given |Gitt |Gegeven |Gegeben sei |Fakat |Eğer ki |Etant donné |Et |Então |Entonces |Entao |En |Eeldades |E |Duota |Dun |Donitaĵo |Donat |Donada |Do |Diyelim ki |Dengan |Den youse gotta |De |Dato |Dar |Dann |Dan |Dado |Dacă |Daca |DEN |Când |Cuando |Cho |Cept |Cand |Cal |But y'all |But |Buh |Biết |Bet |BUT |Atès |Atunci |Atesa |Anrhegedig a |Angenommen |And y'all |And |An |Ama |Als |Alors |Allora |Ali |Aleshores |Ale |Akkor |Aber |AN |A také |A |\* )`

	featureKeywords := `^(기능|機能|功能|フィーチャ|خاصية|תכונה|Функціонал|Функционалност|Функционал|Фича|Особина|Могућност|Özellik|Właściwość|Tính năng|Trajto|Savybė|Požiadavka|Požadavek|Osobina|Ominaisuus|Omadus|OH HAI|Mogućnost|Mogucnost|Jellemző|Fīča|Funzionalità|Funktionalität|Funkcionalnost|Funkcionalitāte|Funcționalitate|Functionaliteit|Functionalitate|Funcionalitat|Funcionalidade|Fonctionnalité|Fitur|Feature|Egenskap|Egenskab|Crikey|Característica|Arwedd)(:)(.*)$`

	featureElementKeywords := `^(\s*)(시나리오 개요|시나리오|배경|背景|場景大綱|場景|场景大纲|场景|劇本大綱|劇本|剧本大纲|剧本|テンプレ|シナリオテンプレート|シナリオテンプレ|シナリオアウトライン|シナリオ|سيناريو مخطط|سيناريو|الخلفية|תרחיש|תבנית תרחיש|רקע|Тарих|Сценарій|Сценарио|Сценарий структураси|Сценарий|Структура сценарію|Структура сценарија|Структура сценария|Скица|Рамка на сценарий|Пример|Предыстория|Предистория|Позадина|Передумова|Основа|Концепт|Контекст|Założenia|Wharrimean is|Tình huống|The thing of it is|Tausta|Taust|Tapausaihio|Tapaus|Szenariogrundriss|Szenario|Szablon scenariusza|Stsenaarium|Struktura scenarija|Skica|Skenario konsep|Skenario|Situācija|Senaryo taslağı|Senaryo|Scénář|Scénario|Schema dello scenario|Scenārijs pēc parauga|Scenārijs|Scenár|Scenaro|Scenariusz|Scenariul de şablon|Scenariul de sablon|Scenariu|Scenario Outline|Scenario Amlinellol|Scenario|Scenarijus|Scenarijaus šablonas|Scenarij|Scenarie|Rerefons|Raamstsenaarium|Primer|Pozadí|Pozadina|Pozadie|Plan du scénario|Plan du Scénario|Osnova scénáře|Osnova|Náčrt Scénáře|Náčrt Scenáru|Mate|MISHUN SRSLY|MISHUN|Kịch bản|Konturo de la scenaro|Kontext|Konteksts|Kontekstas|Kontekst|Koncept|Khung tình huống|Khung kịch bản|Háttér|Grundlage|Geçmiş|Forgatókönyv vázlat|Forgatókönyv|Fono|Esquema do Cenário|Esquema do Cenario|Esquema del escenario|Esquema de l'escenari|Escenario|Escenari|Dis is what went down|Dasar|Contexto|Contexte|Contesto|Condiţii|Conditii|Cenário|Cenario|Cefndir|Bối cảnh|Blokes|Bakgrunn|Bakgrund|Baggrund|Background|B4|Antecedents|Antecedentes|All y'all|Achtergrond|Abstrakt Scenario|Abstract Scenario)(:)(.*)$`

	examplesKeywords := `^(\s*)(예|例子|例|サンプル|امثلة|דוגמאות|Сценарији|Примери|Приклади|Мисоллар|Значения|Örnekler|Voorbeelden|Variantai|Tapaukset|Scenarios|Scenariji|Scenarijai|Příklady|Példák|Príklady|Przykłady|Primjeri|Primeri|Piemēri|Pavyzdžiai|Paraugs|Juhtumid|Exemplos|Exemples|Exemplele|Exempel|Examples|Esempi|Enghreifftiau|Ekzemploj|Eksempler|Ejemplos|EXAMPLZ|Dữ liệu|Contoh|Cobber|Beispiele)(:)(.*)$`

	return Rules{
		"comments": {
			{`\s*#.*$`, Comment, nil},
		},
		"featureElements": {
			{stepKeywords, Keyword, Push("stepContentStack")},
			Include("comments"),
			{`(\s|.)`, NameFunction, nil},
		},
		"featureElementsOnStack": {
			{stepKeywords, Keyword, Pop(2)},
			Include("comments"),
			{`(\s|.)`, NameFunction, nil},
		},
		"examplesTable": {
			{`\s+\|`, Keyword, Push("examplesTableHeader")},
			Include("comments"),
			{`(\s|.)`, NameFunction, nil},
		},
		"examplesTableHeader": {
			{`\s+\|\s*$`, Keyword, Pop(2)},
			Include("comments"),
			{`\\\|`, NameVariable, nil},
			{`\s*\|`, Keyword, nil},
			{`[^|]`, NameVariable, nil},
		},
		"scenarioSectionsOnStack": {
			{featureElementKeywords, ByGroups(NameFunction, Keyword, Keyword, NameFunction), Push("featureElementsOnStack")},
		},
		"narrative": {
			Include("scenarioSectionsOnStack"),
			{`(\s|.)`, NameFunction, nil},
		},
		"tableVars": {
			{`(<[^>]+>)`, NameVariable, nil},
		},
		"numbers": {
			{`(\d+\.?\d*|\d*\.\d+)([eE][+-]?[0-9]+)?`, LiteralString, nil},
		},
		"string": {
			Include("tableVars"),
			{`(\s|.)`, LiteralString, nil},
		},
		"pyString": {
			{`"""`, Keyword, Pop(1)},
			Include("string"),
		},
		"stepContentRoot": {
			{`$`, Keyword, Pop(1)},
			Include("stepContent"),
		},
		"stepContentStack": {
			{`$`, Keyword, Pop(2)},
			Include("stepContent"),
		},
		"stepContent": {
			{`"`, NameFunction, Push("doubleString")},
			Include("tableVars"),
			Include("numbers"),
			Include("comments"),
			{`(\s|.)`, NameFunction, nil},
		},
		"tableContent": {
			{`\s+\|\s*$`, Keyword, Pop(1)},
			Include("comments"),
			{`\\\|`, LiteralString, nil},
			{`\s*\|`, Keyword, nil},
			{`"`, LiteralString, Push("doubleStringTable")},
			Include("string"),
		},
		"doubleString": {
			{`"`, NameFunction, Pop(1)},
			Include("string"),
		},
		"doubleStringTable": {
			{`"`, LiteralString, Pop(1)},
			Include("string"),
		},
		"root": {
			{`\n`, NameFunction, nil},
			Include("comments"),
			{`"""`, Keyword, Push("pyString")},
			{`\s+\|`, Keyword, Push("tableContent")},
			{`"`, NameFunction, Push("doubleString")},
			Include("tableVars"),
			Include("numbers"),
			{`(\s*)(@[^@\r\n\t ]+)`, ByGroups(NameFunction, NameTag), nil},
			{stepKeywords, ByGroups(NameFunction, Keyword), Push("stepContentRoot")},
			{featureKeywords, ByGroups(Keyword, Keyword, NameFunction), Push("narrative")},
			{featureElementKeywords, ByGroups(NameFunction, Keyword, Keyword, NameFunction), Push("featureElements")},
			{examplesKeywords, ByGroups(NameFunction, Keyword, Keyword, NameFunction), Push("examplesTable")},
			{`(\s|.)`, NameFunction, nil},
		},
	}
}
