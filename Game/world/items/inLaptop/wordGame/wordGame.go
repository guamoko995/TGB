package worldGame

// Слот (столбец, колесо, спин) с буквами
type position struct {
	column []rune // Список вариантов
	p      int    // Хранит индекс текущего варианта
	defP   int    // Хранит индекс варианта по-умолчанию
	mark   bool   // Отметка о клике (вращении)
}

// Обновляет столбец в соответствии с переданной буквой R.
func (pos *position) update(R rune) {
	pos.column = []rune{}
	for i, rc := range sortedRusRuneCount {
		if rc.R == R {
			l := len(sortedRusRuneCount)
			lb := i - 2
			ub := i + 3
			if lb < 0 {
				lb = 0
			}
			if ub > l {
				ub = l
			}
			pos.defP = i - lb
			pos.p = pos.defP
			for i := lb; i < ub; i++ {
				pos.column = append(pos.column, sortedRusRuneCount[i].R)
			}
			return
		}
	}
	pos.column = append(pos.column, R)
}

// "Перещелкивает" столбец (колесо) на одну позицию.
func (pos *position) click() {
	l := len(pos.column)
	pos.p++
	if pos.p == l {
		pos.p = 0
	}
}

// Имя столбца соответствует текущей букве.
func (pos *position) Name() rune {
	return pos.column[pos.p]
}

// Машина слотов.
type SlotMakhine struct {
	Str   []*position
	Words [][]rune
	Text  QwestText
	Pos   int
}

// Контрольная сумма. Уникальна для каждого положения спинов.
func (sm *SlotMakhine) kSum() int {
	sum := 0
	for i, pos := range sm.Str {
		sum += pos.p*2 ^ (3 * i)
	}
	return sum
}

// Обнвляет машину слотов.
func (sm *SlotMakhine) Update() {

	text := sm.Text.Print()
	sm.Words = make([][]rune, 0)
	word := []rune{}
	for _, R := range text {
		if 'А' <= R && R <= 'я' {
			word = append(word, R)
		} else if len(word) > 0 {
			sm.Words = append(sm.Words, word)
			word = []rune{}
		}
	}

	sm.Str = []*position{}

	for _, R := range sm.Words[sm.Pos] {
		p := &position{}
		p.update(R)
		sm.Str = append(sm.Str, p)
	}

	for _, R := range sm.Words[sm.Pos] {
		if sm.Text.Count(R+'А'-'а') != 0 {
			sm.SmartClick(R)
		}
	}
}

// Удаляет пометки со всех слотов.
func (sm *SlotMakhine) markFalseAll() {
	for _, pos := range sm.Str {
		pos.mark = false
	}
}

// Перещелкивает помеченные (mark=true) или не помеченные (mark=false)
// слоты с выбранной буквой R на одну позицию.
func (sm *SlotMakhine) next(R rune, mark bool) (nextR rune) {
	nextR = R
	for _, p := range sm.Str {
		if R == p.Name() && p.mark == mark {
			p.click()
			p.mark = true
			nextR = p.Name()
		} else {
			p.mark = false
		}
	}
	return
}

// Перещелкивает помеченные (mark=true) или не помеченные (mark=false)
// слоты с выбранной буквой R с пропуском уже восстановленных букв.
func (sm *SlotMakhine) click(R rune, mark bool) rune {
	ok := func() bool {
		// Если буква уже восстановлена
		return sm.Text.Count(R+'А'-'а') == 0
	}
	bufS := sm.kSum()
	R = sm.next(R, mark)
	for !ok() {
		R = sm.next(R, true)
		if bufS == sm.kSum() {
			break
		}
	}
	return R
}

// Перещелкивает помеченные слоты с выбранной буквой R с пропуском уже
// восстановленных букв и с отслеживанием повторений.
func (sm *SlotMakhine) SmartClickOnePos(R rune) (nextR rune) {
	// Запомним начальное состояние с помощью контрольной суммы.
	bufS := sm.kSum()
	// Кликнем по не задействованым в предыдущий клик позициям R.
	R = sm.click(R, false)
	// Если нет таких позиций то ничего не изменится, дальше делать
	// нечего.
	if bufS == sm.kSum() {
		return
	}
	// Задаем критерий допустимости позиции
	ok := func() bool {
		for _, pos := range sm.Str {
			if pos.Name() == R && !pos.mark {
				return false
			}
		}
		return true
	}

	for !ok() {
		R = sm.click(R, true)
		if bufS == sm.kSum() {
			R = sm.click(R, true)
			break
		}
	}
	return R
}

// Перещелкивает помеченные слоты с выбранной буквой R с пропуском уже
// восстановленных букв и с отслеживанием повторений. При отсутствии
// новых вариантов перещелкивает "мешающие" слоты.
func (sm *SlotMakhine) SmartClick(R rune) {
	sm.markFalseAll()
	bufS := sm.kSum()

	// Задаем критерий допустимости позиции
	ok := func() bool {
		for _, pos := range sm.Str {
			if pos.Name() == R && !pos.mark {
				return false
			}
		}
		return true
	}

	for !ok() {
		R = sm.SmartClickOnePos(R)
		if bufS == sm.kSum() {
			break
		}
	}
}

// Генерирует карту замен по текущему состоянию слот машины.
func (sm *SlotMakhine) ReplaceMap() map[rune]rune {
	mp := make(map[rune]rune)
	for _, pos := range sm.Str {
		if 'а' <= pos.column[pos.defP] && pos.column[pos.defP] <= 'я' {
			mp[pos.column[pos.defP]] = pos.column[pos.p] - 'а' + 'А'
		}
	}
	return mp
}

// Устанавливает значение текста.
func (sm *SlotMakhine) SetText(qt QwestText) {
	qt.Down()
	cm := DecryptMap(qt)
	qt.Crypt(cm)

	sm.Text = qt
	sm.Update()
}
