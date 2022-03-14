package worldGame

import (
	"crypto/rand"
	"io"
	"math/big"
	"os"
	"sort"
	"strings"
)

// Возвращает массив с рунами русского алфавита.
func ABC() []rune {
	return []rune{'а', 'б', 'в', 'г', 'д', 'е', 'ж', 'з', 'и', 'й', 'к',
		'л', 'м', 'н', 'о', 'п', 'р', 'с', 'т', 'у', 'ф', 'х', 'ц', 'ч',
		'ш', 'щ', 'ъ', 'ы', 'ь', 'э', 'ю', 'я'}
}

// Генерирует случайную карту замены русского алфавита (без ё). Алфавит
// отображается в себя.
func GenCryptMap() map[rune]rune {
	abc := ABC()
	l := int64(len(abc))
	tab := make(map[rune]rune)
	for _, b := range ABC() {
		big, _ := rand.Int(rand.Reader, big.NewInt(int64(l)))
		r := int((*big).Int64())
		l--
		tab[b] = abc[r]
		abc = append(abc[:r], abc[r+1:]...)
	}
	return tab
}

// Возвращает текст из файла.
func strFromFile(fileName string) string {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	data := make([]byte, 64)

	var n int
	var str string
	for {
		n, err = file.Read(data)
		if err == io.EOF { // если конец файла
			break // выходим из цикла
		}
		str += string(data[:n])
	}
	str += string(data[:n])
	return str
}

// Условный текст, с которым можно работать с помощью виртуального
// текстового редактора.
type QwestText interface {
	Replace(original, new rune)
	Count(R rune) int
	CountAll() int
	Reset()
	Up()
	Down()
	Print() string
	Crypt(cryptMap map[rune]rune)
}

// Массив условных текстов также является условным текстом.
type MQT []QwestText

func (mt MQT) Replace(original, new rune) {
	for i := range mt {
		mt[i].Replace(original, new)
	}
}

func (mt MQT) Count(R rune) int {
	resp := 0
	for i := range mt {
		resp += mt[i].Count(R)
	}
	return resp
}

func (mt MQT) CountAll() int {
	resp := 0
	for i := range mt {
		resp += mt[i].CountAll()
	}
	return resp
}

func (mt MQT) Reset() {
	for i := range mt {
		mt[i].Reset()
	}
}

func (mt MQT) Up() {
	for i := range mt {
		mt[i].Up()
	}
}

func (mt MQT) Down() {
	for i := range mt {
		mt[i].Down()
	}
}

func (mt MQT) Print() string {
	str := make([]string, 0)
	for i := range mt {
		str = append(str, mt[i].Print())
	}
	return strings.Join(str, " ")
}

func (mt MQT) Crypt(cryptMap map[rune]rune) {
	for i := range mt {
		mt[i].Crypt(cryptMap)
	}
}

// Обычный текст, с которым можно работать с помощью виртуального
// текстового редактора.
type QText struct {
	Real     string
	Original string // Хранит исходный текст для отката изменений.
}

// Заменяет символы в тексте.
func (t *QText) Replace(original, new rune) {
	t.Real = strings.ReplaceAll(t.Real, string([]rune{original}), string([]rune{new}))
}

// Считает количество символов в тексте.
func (t *QText) Count(R rune) int {
	return strings.Count(t.Real, string([]rune{R}))
}

// Считает количество всех русских букв в тексте.
func (t *QText) CountAll() int {
	var n int
	for _, R := range t.Real {
		if 'А' <= R && R <= 'я' {
			n++
		}
	}
	return n
}

// Восстанавливает текст до исходного
func (t *QText) Reset() {
	t.Real = t.Original
}

// Делает все русские буквы в тексте заглавными.
func (t *QText) Up() {
	t.Real = strings.ToUpper(t.Real)
}

// Делает все русские буквы в тексте строчными.
func (t *QText) Down() {
	t.Real = strings.ToLower(t.Real)
}

// Возвращает строку текста.
func (t *QText) Print() string {
	return string(t.Real)
}

// Производит простую замену русских букв по карте.
func (t *QText) Crypt(cryptMap map[rune]rune) {
	crypt := []rune(t.Real)
	for i := range crypt {
		if s, ok := cryptMap[crypt[i]]; ok {
			crypt[i] = s
		} else if s, ok := cryptMap[crypt[i]+('а'-'А')]; ok {
			crypt[i] = s + ('А' - 'а')
		}
	}
	t.Real = string(crypt)
	t.Original = t.Real
}

// Сущность содержащая информацию о количестве русских букв в тексте.
// Используется для ускоренного частотного криптоанализа.
type PsevdoText []RuneCount

// Заменяет буквы в псевдотексте
func (pt PsevdoText) Replace(original, new rune) {
	for i := range pt {
		if pt[i].R == original {
			pt[i].R = new
		}
	}
}

// Считает количество указанных русских букв в псевдотексте.
func (pt PsevdoText) Count(R rune) int {
	var n int
	for i := range pt {
		if pt[i].R == R {
			n += pt[i].Count
		}
	}
	return n
}

// Считает количество всех русских букв в псевдотексте.
func (pt PsevdoText) CountAll() int {
	var n int
	for i := range pt {
		n += pt[i].Count
	}
	return n
}

// Для корректной работы функции текст изначально должен быть получен с
// помощью функции NewPsevdoText. Восстанавливает текст до состояния на
// момент получения псевдотекста из реального текста (до всех замен)
func (pt PsevdoText) Reset() {
	for i, R := range ABC() {
		offset := 'А' - 'а'
		pt[i-int(offset)].R = R

		R = R + offset
		pt[i].R = R
	}
}

// Заменяет все строчные буквы псевдотекста на заглавные.
func (pt PsevdoText) Up() {
	for i, rc := range pt {
		if rc.R > 'Я' { // если строчная
			pt[i].R = rc.R + 'А' - 'а'
		}
	}
}

// Заменяет все заглавные буквы псевдотекста строчными.
func (pt PsevdoText) Down() {
	for i, rc := range pt {
		if rc.R < 'а' { // если Заглавная
			pt[i].R = rc.R - 'А' + 'а'
		}
	}
}

// Нельзя восстановить содержимое реального родительского текста из
// псевдотекста. Попытка показать вернет пустую строку.
func (pt PsevdoText) Print() string {
	return ""
}

// Преобразует псевдотекст в новый исходно зашифрованный простой заменой PsevdoText.
// ВАЖНО!!! Отменяет преобразования регистра.
func (pt PsevdoText) Crypt(fullCryptMap map[rune]rune) {
	newPt := make([]RuneCount, len(pt))
	copy(newPt, []RuneCount(pt))
	for j, rc := range pt {
		for _, nrc := range newPt {
			if fullCryptMap[nrc.R] == rc.R || fullCryptMap[nrc.R+'a'-'A'] == rc.R+'a'-'A' {
				pt[j].Count = nrc.Count
			}
		}
	}
}

// Структура для хранения количества букв в тексте, вместо хранения
// самого текста. Используется для ускоренного подсчета букв в тексте.
type RuneCount struct {
	R     rune
	Count int
}

// Получает псевдотекст из реального для ускорения дальнейшего
// частотно-символьного анализа.
func NewPsevdoText(fileName string) PsevdoText {
	text := strFromFile(fileName)
	RuneInText := make(map[rune]int)
	for _, s := range string(text) {
		for _, sAbc := range ABC() {
			if s == sAbc {
				RuneInText[s]++
			} else if sAbc += ('А' - 'а'); s == sAbc {
				RuneInText[s]++
			}
		}
	}
	half := len(ABC())
	masRuneCount := make([]RuneCount, 2*half)
	for i, R := range ABC() {
		masRuneCount[i+half].R = R
		Count := RuneInText[R]
		masRuneCount[i+half].Count = Count

		R = R + 'А' - 'а'
		masRuneCount[i].R = R
		Count = RuneInText[R]
		masRuneCount[i].Count = Count
	}
	return PsevdoText(masRuneCount)
}

func NewQText(text string) *QText {
	return &QText{text, text}
}

// Возвращает карту дешифровки в соответствии с частотнотью русских
// букв.
func DecryptMap(text QwestText) map[rune]rune {
	l := len(ABC())
	textRuneCount := make([]RuneCount, l)
	for i, R := range ABC() {
		textRuneCount[i].R = R
		textRuneCount[i].Count = text.Count(R)
		R += 'А' - 'а'
		textRuneCount[i].Count += text.Count(R)
	}

	sort.Slice(textRuneCount, func(i, j int) bool { return textRuneCount[i].Count < textRuneCount[j].Count })
	Decrypt := make(map[rune]rune)
	for i, sR := range sortedRusRuneCount {
		Decrypt[textRuneCount[i].R] = sR.R
	}
	return Decrypt
}

// Отсортированный по встречаемости массив букв русского алфавита.
var sortedRusRuneCount []RuneCount = func() []RuneCount {
	// Встречаемость букв русского языка (на материале НКРЯ)
	rusRuneCount := []RuneCount{
		{'а', 40487008},
		{'б', 8051767},
		{'в', 22930719},
		{'г', 8564640},
		{'д', 15052118},
		{'е', 42876141},
		{'ж', 4746916},
		{'з', 8329904},
		{'и', 37153142},
		{'й', 6106262},
		{'к', 17653469},
		{'л', 22230174},
		{'м', 16203060},
		{'н', 33838881},
		{'о', 55414481},
		{'п', 14201572},
		{'р', 23916825},
		{'с', 27627040},
		{'т', 31620970},
		{'у', 13245712},
		{'ф', 1335747},
		{'х', 4904176},
		{'ц', 2438807},
		{'ч', 7300193},
		{'ш', 3678738},
		{'щ', 1822476},
		{'ъ', 185452},
		{'ы', 9595941},
		{'ь', 8784613},
		{'э', 1610107},
		{'ю', 3220715},
		{'я', 10139085},
	}
	sort.Slice(rusRuneCount, func(i, j int) bool { return rusRuneCount[i].Count < rusRuneCount[j].Count })
	return rusRuneCount
}()
