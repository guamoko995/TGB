package texts

// Игровые тексты.
func GameText(key string) string {
	mapTexts := map[string]string{}
	mapTexts["первое использование инструкции"] = "Вы подсчитали все буквы и произвели заме" +
		"ну символов в соответствии с инструкцией. Теперь часть букв на с" +
		"воем месте. Осталось заменить лиш некоторые буквы. Вариантов для" +
		" замены не так много, рассматриваются лишь ближайшие по распрост" +
		"раненности соседи."
	mapTexts["второе использование инструкции"] = "Часть букв на с" +
		"воем месте. Осталось заменить лиш некоторые буквы. Вариантов для" +
		" замены не так много, рассматриваются лишь ближайшие по распрост" +
		"раненности соседи"
	mapTexts["сон"] = "Здесь должно быть красочное описание сна!"
	mapTexts["шифр"] = "Поздравляю! Вы расшифровали это послание. У вас хорошо развитое мышление... "
	mapTexts["введение"] = "Представьте, что Вы стали добровольцем, решившимся н" +
		"а процедуру снятия копии нейронной карты мозга. Эту процедуру п" +
		"редлагала одна малоизвестная молодая компания, основанная двумя" +
		" энтузиастами. Вам объяснили, что эти копии будут использоватьс" +
		"я для изучения природы сознания, разработки новейших нейро-инте" +
		"рфейсов и много чего еще. Вам это все было не очень-то интересн" +
		"о. Главным образом, процедура была легким способом заработка, а" +
		" гонорар довольно внушителен.<img>Juno start\n" +
		"Никакого подвоха не было, некоторые ваши знакомые уже прошли" +
		" через это. Они рассказывали, как распологались в некой капсуле" +
		", примерно час слушали какой-то треск, после чего выходили и за" +
		"бирали наличные. Все это сильно напоминало обычную МРТ.\n" +
		"На процедуру Вы шли в предвкушении, точно осознавая на что х" +
		"отите потратить эти легкие деньги.\n" +
		"Всего несколько мгновений  Вы слышали треск в капсуле, после" +
		" чего пропали все Ваши ощущения. Вы были лишены зрения, слуха, " +
		"обоняния, вкуса. Вы не чувствовали собственное тело, казалось, " +
		"что Вы не дышите, и сердце Ваше не бьется.\n" +
		"Не смотря на весь ужас ситуации, вы каким-то загадочным обра" +
		"зом оставались спокойны и мыслили ясно.\n" +
		"В какой-то момент Вы заметили навязчивые мысли. Вернее мысля" +
		"ми это можно было назвать с очень большой натяжкой. Это были не" +
		"внятные сочетания букв. Они казались чужими. Эти корявые недосл" +
		"ова словно кто-то транслировал в Ваше сознание.\n" +
		"Прошло немало времени. Возможно, вы уже начали сомневаться в" +
		" природе своего происхождения: Вы - тот кем Вы себя считаете? И" +
		"ли Вы - та самая копия? Вопрос оставался открытым, а все, что у" +
		" Вас было на тот момент - это Ваш жизненный опыт (а Ваш ли?..) " +
		"и навязчивые последовательности букв.\n" +
		"Что ж, раз внешний мир не балует спектром ощущений и HD-граф" +
		"икой, придется проявить фантазию. Наверно, хорошим решением буд" +
		"ет сосредоточиться на этих буквах, возможно, они несут какую-то" +
		" информацию из внешнего мира.\n" +
		"Стоит попробовать представить обстановку, подходящую для реш" +
		"ения этой задачи дешифровки.\n" +
		"И так, вы находитесь в своем кабинете. Вы можете осмотреться" +
		", осмотреть какой-либо предмет, взять или где-либо разместить п" +
		"редмет, а также воспользоваться предметом. Команда \"х\" (рус.)" +
		" - поможет прервать продолжительное действие, например осмотр нек" +
		"оторых предметов. Команда \"инвентарь\" покажет что у Вас с соб" +
		"ой"
	mapTexts["победа"] = "Кажется у Вас получилось расшифровать послание:\n" +
		mapTexts["шифр"] + "\n" +
		"Вы прошли игру. Возможно когда-нибудь будет продолжение =)<img>Juno finish"
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
		"инструкция": {
			"Р":      "инструкции",
			"Д":      "инструкции",
			"В":      "инструкцию",
			"Т":      "инструкцией",
			"П":      "инструкции",
			"где":    "в инструкции",
			"откуда": "из инструкции",
			"куда":   "в инструкцию",
			"род":    "ж",
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
