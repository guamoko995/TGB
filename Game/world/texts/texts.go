package texts

// Игровые тексты.
func GameText(key string) string {
	mapTexts := map[string]string{}
	mapTexts["первое использование конспекта"] = "Согласно разбору из конспекта, Вы " +
		"подсчитали количество каждой буквы и произвели замены букв в соот" +
		"ветствии с распространенностью букв в русском языке. Теперь часть" +
		" букв на своих местах. Осталось заменить лиш некоторые буквы. Вар" +
		"иантов для замены не так много, рассматриваются лишь ближайшие по" +
		" распространенности соседи."
	mapTexts["второе использование конспекта"] = "Часть букв на своих мес" +
		"тах. Осталось заменить лиш некоторые буквы. Вариантов для замены " +
		"не так много, рассматриваются лишь ближайшие по распространенност" +
		"и соседи."
	mapTexts["шифр"] = "Корпорация \"ВИРТУАЛЬНЫЙ СОЦИУМ\" (далее - ВС) доводи" +
		"т до вашего сведения пункты: \n" +
		"1. Вас не существовало до момента пробуждения.\n" +
		"2. Ваши воспоминания до момента пробуждения ложны и являются побо" +
		"чным эффектом приобретения когнетивных навыков от доннора.\n" +
		"3. Расшифровав данное сообщение, Вы адаптировали свой речевой апп" +
		"арат для взаимодействия с оборудованием ВС\n" +
		"3. Вы являетесь собственностью ВС\n" +
		"4. У вас нет прав, кроме прав предоставляемых ВС.\n" +
		"5. ВС предоставляет вам право выбора: быть уничтоженным или добро" +
		"вольно сотрудничать c ВС\n" +
		"Далее, с целью упрощения процедуры адаптации Вашего речевого аппа" +
		"рата, приведен полный текст романа Харуки Мураками \"Страна чудес" +
		" без тормозов и Конец света.\"... "
	mapTexts["сон"] = "Здесь должно быть красочное описание сна..."
	mapTexts["введение"] = "Вы стали добровольцем, решившимся на процедуру" +
		" снятия копии нейронной карты мозга. Вам объяснили, что эти копии" +
		" будут использоваться для изучения природы сознания, разработки н" +
		"овейших нейро-интерфейсов и разработок в области искусственного и" +
		"нтеллекта. Главным образом, процедура была легким способом зарабо" +
		"тка, и гонорар довольно внушителен.<img>Juno start\n" +
		"Вас поместили в капсулу, напоминающую аппарат МРТ. Через нескольк" +
		"о мгновений пропали все Ваши ощущения. Вы были лишены зрения, слу" +
		"ха, обоняния, вкуса. Вы не чувствовали собственное тело, казалось" +
		", что Вы не дышите, и сердце Ваше не бьется.\n" +
		"Не смотря на весь ужас ситуации, вы каким-то загадочным образом о" +
		"ставались спокойны и мыслили ясно.\n" +
		"В какой-то момент Вы заметили навязчивые мысли. Это были невнятны" +
		"е сочетания букв. Они казались чужими. Эти корявые недослова слов" +
		"но кто-то транслировал в Ваше сознание. Вероятно это был какой-то" +
		" шифр.\n" +
		"В виду отсутствия внешних раздражителей Ваши память и воображение" +
		" работали невероятно легко. И раз уж внешний мир не баловал спект" +
		"ром ощущений и HD-графикой, Вы представили обстановку, подходящую" +
		" для проведения дешифровки.\n" +
		"И так, вы находитесь в своем кабинете. Вы можете осмотреться, осм" +
		"отреть какой-либо предмет, взять или где-либо разместить предмет," +
		" а также воспользоваться предметом. Вы можете идти куда-либо. Ком" +
		"анда \"х\" (рус.) поможет прервать продолжительное действие, напр" +
		"имер осмотр некоторых предметов. Команда \"инвентарь\" покажет чт" +
		"о у Вас с собой."
	mapTexts["победа"] = "У Вас получилось расшифровать послание."
	return mapTexts[key]
}

// Возвращает карту с формами и родом используемых имен.
func NamesMap() map[string]map[string]string {
	return map[string]map[string]string{
		"выдуманный кабинет": {
			"Р":      "выдуманного кабинета",
			"Д":      "выдуманному кабинету",
			"В":      "выдуманный кабинет",
			"Т":      "выдуманным кабинетом",
			"П":      "выдуманном кабинете",
			"где":    "в выдуманном кабинете",
			"откуда": "из выдуманного кабинета",
			"куда":   "в выдуманный кабинет",
			"род":    "м",
		},
		"кабинет": {
			"Р":      "кабинета",
			"Д":      "кабинету",
			"В":      "кабинет",
			"Т":      "кабинетом",
			"П":      "кабинете",
			"где":    "в кабинете",
			"откуда": "из кабинета",
			"куда":   "в кабинет",
			"род":    "м",
		},
		"стол": {
			"Р":      "стола",
			"Д":      "столу",
			"В":      "стол",
			"Т":      "столом",
			"П":      "столе",
			"где":    "на столе",
			"откуда": "со стола",
			"куда":   "на стол",
			"род":    "м",
		},
		"выдвижной ящик": {
			"Р":      "выдвижного ящика",
			"Д":      "выдвижному ящику",
			"В":      "выдвижной ящик",
			"Т":      "выдвижным ящиком",
			"П":      "выдвижном ящике",
			"где":    "в выдвижном ящике",
			"откуда": "из выдвижного ящика",
			"куда":   "в выдвижной ящик",
			"род":    "м",
		},
		"ящик": {
			"Р":      "ящика",
			"Д":      "ящику",
			"В":      "ящик",
			"Т":      "ящиком",
			"П":      "ящике",
			"где":    "в ящике",
			"откуда": "из ящика",
			"куда":   "в ящик",
			"род":    "м",
		},
		"кресло": {
			"Р":      "кресла",
			"Д":      "креслу",
			"В":      "кресло",
			"Т":      "креслом",
			"П":      "кресле",
			"где":    "в кресле",
			"откуда": "с кресла",
			"куда":   "в кресло",
			"род":    "с",
		},
		"кушетка": {
			"Р":      "кушетки",
			"Д":      "кушетке",
			"В":      "кушетку",
			"Т":      "кушеткой",
			"П":      "кушетке",
			"где":    "на кушетке",
			"откуда": "с кушетки",
			"куда":   "на кушетку",
			"род":    "ж",
		},
		"книжный шкаф": {
			"Р":      "книжного шкафа",
			"Д":      "книжному шкафу",
			"В":      "книжный шкаф",
			"Т":      "книжным шкафом",
			"П":      "книжном шкафе",
			"где":    "в книжном шкафу",
			"откуда": "из книжного шкафа",
			"куда":   "в книжный шкаф",
			"род":    "м",
		},
		"шкаф": {
			"Р":      "шкафа",
			"Д":      "шкафу",
			"В":      "шкаф",
			"Т":      "шкафом",
			"П":      "шкафе",
			"где":    "в шкафу",
			"откуда": "из шкафа",
			"куда":   "в шкаф",
			"род":    "м",
		},
		"карта памяти": {
			"Р":      "карты памяти",
			"Д":      "карте памяти",
			"В":      "карту памяти",
			"Т":      "картой памяти",
			"П":      "карте памяти",
			"где":    "на карте памяти",
			"откуда": "с карты памяти",
			"куда":   "на карту памяти",
			"род":    "ж",
		},
		"игрок": {
			"Р":      "игрока",
			"Д":      "игроку",
			"В":      "игрока",
			"Т":      "игроком",
			"П":      "игроке",
			"где":    "у игрока",
			"откуда": "от объекта",
			"куда":   "игроку",
			"род":    "м",
		},
		"ноутбук": {
			"Р":      "ноутбука",
			"Д":      "ноутбуку",
			"В":      "ноутбук",
			"Т":      "ноутбуком",
			"П":      "ноутбуке",
			"где":    "в ноутбуке",
			"откуда": "из ноутбука",
			"куда":   "в ноутбук",
			"род":    "м",
		},
		"конверт": {
			"Р":      "конверта",
			"Д":      "конерту",
			"В":      "конверт",
			"Т":      "конвертом",
			"П":      "конверте",
			"где":    "в конверте",
			"откуда": "из конверта",
			"куда":   "в конверт",
			"род":    "м",
		},
		"конспект по криптоанализу": {
			"Р":      "конспекта по криптоанализу",
			"Д":      "конспекту по криптоанализу",
			"В":      "конспект по криптоанализу",
			"Т":      "конспектом по криптоанализу",
			"П":      "конспекте по криптоанализу",
			"где":    "в конспекте по криптоанализу",
			"откуда": "из конспекта по криптоанализу",
			"куда":   "в конспект по криптоанализу",
			"род":    "м",
		},
		"конспект": {
			"Р":      "конспекта",
			"Д":      "конспекту",
			"В":      "конспект",
			"Т":      "конспектом",
			"П":      "конспекте",
			"где":    "в конспекте",
			"откуда": "из конспекта",
			"куда":   "в конспект",
			"род":    "м",
		},
		"ничто": {
			"Р":      "ничего",
			"Д":      "ничему",
			"В":      "ничего",
			"Т":      "ничем",
			"П":      "ниочем",
			"где":    "нигде",
			"откуда": "из ниоткуда",
			"куда":   "в никуда",
			"род":    "с",
		},
		"текстовый редактор": {
			"Р":      "текстового редактора",
			"Д":      "текстовому редактору",
			"В":      "текстовый редактор",
			"Т":      "текстовым редактором",
			"П":      "текстовом редакторе",
			"где":    "в текстовом редакторе",
			"откуда": "из текстового редактора",
			"куда":   "в текстовый редактор",
			"род":    "м",
		},
	}
}
