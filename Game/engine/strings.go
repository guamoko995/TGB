package engine

import (
	"TelegramGameBot/Game/base"
	"fmt"
	"strings"
)

// Перемещает предмет с подробным текстовым результатом, как в случае
// успеха так и при ошибке.
func Place(obj base.Positioner, cont base.Conteiner) string {
	err := base.Place(obj, cont)
	if err != nil {
		return TitleAfterDot(err.Error())
	}
	o := Postfix(obj)
	return obj.Name() + " помещен" + o + " " + cont.Name("куда")
}

// Если текст представим в виде конкатенации ключа одного из приложений
// объекта с некоторым остатком строки, то возвращает true и остаток
// строки.
func ConsumeAppUs(h *TreeHandlers, text string) (Handler, string) {
	for nameApp, app := range h.Applications {
		if ok, endText := ConsumePrefixWords(nameApp, text); ok {
			return app, endText
		}
	}
	return nil, text
}

// Возвращает true если ключ одного из приложений объекта начинается с
// переданной строки.
func ChekPartAppUs(h *TreeHandlers, PartNameApp string) bool {
	for nameApp := range h.Applications {
		if ok, _ := ConsumePrefixWords(PartNameApp, nameApp); ok {
			return true
		}
	}
	return false
}

// Возвращает окончание существительного в роде переданного объекта.
func Postfix(n base.Namer) string {
	switch n.Info("род") {
	case "ж":
		return "а"
	case "с":
		return "о"
	case "м":
		return ""
	default:
		return "(а/о)"
	}
}

// Возвращает окончание прилагательного в роде переданного объекта.
func PostfixPril(n base.Namer) string {
	switch n.Info("род") {
	case "ж":
		return "ая"
	case "с":
		return "ое"
	case "м":
		return "ой"
	default:
		return "ой(ая/ое)"
	}
}

func InputFormat(input string) string {
	input = strings.ToLower(input)
	for {
		sNew := strings.ReplaceAll(input, "  ", " ")
		if sNew == input {
			break
		}
		input = sNew
	}
	return input
}

func OutputFormat(output string) string {
	if output != "" {
		output = TitleAfterDot(output + ".")
		output = strings.ReplaceAll(output, "..", ".")
		output = strings.ReplaceAll(output, "..", "...")
		output = strings.ReplaceAll(output, "?.", "?")
		output = strings.ReplaceAll(output, "!.", "!")
	}
	return output
}

// Если текст представим в виде конкатенации одной из форм имени объекта
// с некоторым остатком строки, то возвращает true и остаток строки.
func ConsumeNameObj(obj base.Namer, text string) (bool, string) {
	for _, name := range obj.Names() {
		if ok, endText := ConsumePrefixWords(name, text); ok {
			return true, endText
		}
	}
	return false, text
}

// Возвращает true если одна из форм имени объекта начинается с
// переданной строки.
func ChekPartName(obj base.Namer, PartName string) bool {
	for _, name := range obj.Names() {
		if ok, _ := ConsumePrefixWords(PartName, name); ok {
			return true
		}
	}
	return false
}

// Добавляет массив форм имени с префиксом в ключах.
func NamesAddWithPrefix(obj base.Namer, keysAndNames [][]string, prefix string) {
	for _, str := range keysAndNames {
		obj.AddName(base.Title{Form: prefix + str[0], Value: str[1]})
	}
}

// Поиск объекта с одной из форм имени которого начинается текст, форма
// имени обязательно должна быть отделена пробелом. Вторым значеием
// возвращается остаток текста без найденного объекта.
func ConsumePositionerFoundByName(text string, fp ...base.FindPosition) (base.Positioner, string) {
	criterion := func(obj base.Positioner) bool {
		ok, _ := ConsumeNameObj(obj, text)
		return ok
	}
	obj := base.Find(criterion, fp...)
	if obj == nil {
		return nil, text
	}
	ok, endText := ConsumeNameObj(obj, text)
	if !ok {
		panic("Объект, найденный по критерию не соответствует критерию.")
	}
	return obj, endText
}

// Возвращает true если в поисковой позиции fp найдется хоть один
// объект, любая форма имени которого начинается с PartName.
func FindObjByPartName(PartName string, fp ...base.FindPosition) bool {
	criterion := func(obj base.Positioner) bool {
		return ChekPartName(obj, PartName)
	}
	obj := base.Find(criterion, fp...)
	return obj != nil
}

// Соединяет переданные строки в одну с разделителем ", ", в качестве
// последнего разделителя используется " и ".
func List(ms ...string) string {
	l := len(ms)
	if l > 1 {
		return strings.Join(ms[:l-1], ", ") + " и " + ms[l-1]
	}
	if l == 1 {
		return ms[0]
	}
	return ""
}

// Делает первый символ и каждый, следующий за ". " заглавными.
func TitleAfterDot(text string) string {
	sent := strings.Split(text, ". ")
	newSent := make([]string, 0)
	for _, s := range sent {
		word := strings.SplitN(s, " ", 2)
		word[0] = strings.Title(word[0])
		newSent = append(newSent, strings.Join(word, " "))
	}
	return strings.Join(newSent, ". ")
}

// Если текст представим в виде конкатенации PrefixWords с некоторым
// остатком текста, в том числе нулевым, то возвращает true и остаток
// текста.
func ConsumePrefixWords(PrefixWords string, text string) (ok bool, endText string) {
	ok = strings.HasPrefix(text+" ", PrefixWords+" ")
	if ok {
		endText = strings.Replace(text, PrefixWords, "", 1)
		endText = strings.TrimLeft(endText, " ")
	}
	return
}

// Соединяет не пустые строки с разделителем. Cтрока из одних пробелов
// считается пустой.
func StrJoin(str []string, sep string) string {
	mstr := make([]string, 0)
	for _, s := range str {
		if strings.ReplaceAll(s, " ", "") != "" {
			mstr = append(mstr, s)
		}
	}
	return strings.Join(mstr, sep)
}

// Получает мир в котором находится объект.
func RootConteiner(d base.Positioner) *World {
	c := base.RootConteiner(d)
	if w, ok := c.(*World); ok {
		return w
	}
	panic("Тип RootConteiner не World")
}

// Строковое представление контейнера.
func StrConteiner(c base.Conteiner) string {
	ms := make([]string, 0)
	for _, obj := range c.Content() {
		ms = append(ms, obj.Name())
	}
	if len(ms) > 0 {
		return TitleAfterDot(c.Name("где") + " " + List(ms...))
	}
	return "Пуст" + PostfixPril(c) + " " + c.Name()
}

//Ошибка пустая строка
var ErrEmptyStr = fmt.Errorf("предмет не указан")
