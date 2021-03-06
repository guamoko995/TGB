package world

import (
	"TelegramGameBot/Game/base"
	"TelegramGameBot/Game/engine"
	"TelegramGameBot/Game/world/buildTools"
	"TelegramGameBot/Game/world/items"
	"TelegramGameBot/Game/world/texts"
	"strings"
)

// Создает новый игровой мир при запуске игры.
func Constructor() *engine.World {
	// Создание экземпляра игрового мира.
	w := (*engine.World).New(&engine.World{})

	// Создание игрока.
	w.Pl = (*items.Player).New(&items.Player{})
	buildTools.SetName(w.Pl, "Вы")

	// Создание повествователя.
	Nr := &Narrator{
		Texts:           strings.Split(texts.GameText("введение"), "\n"),
		W:               w,
		NextImplementer: w.Pl,
	}

	// Назначение повествователя исполнителем команд.
	w.NewActiveHandler(Nr)

	var Obj base.Positioner
	var Pos base.Conteiner

	// Создание локации выдуманный кабинет.
	Obj = w.AddLocation()
	buildTools.SetName(Obj, "выдуманный кабинет")
	buildTools.SetAdditName(Obj, "кабинет", "сокр")

	// Размещение игрока в выдуманный кабинет.
	base.Place(w.Pl, w.Locations("кабинет"))

	// Создание локации ничто.
	Obj = w.AddLocation()
	buildTools.SetName(Obj, "ничто")

	// Создание прохода из ниоткуда в кабинет.
	Obj.(*engine.Location).Bridge.Add(w.Locations("кабинет"))

	// Создание прохода из кабинета в никуда.
	loc := (base.Positioner(w.Locations("кабинет"))).(*engine.Location)
	loc.Bridge.Add(w.Locations("ничто"))

	// Создание стола со встроенным ящиком.
	Obj = (*items.TableWithDrawer).New(&items.TableWithDrawer{})

	// Размещение стола в кабинете.
	base.Place(Obj, w.Locations("кабинет"))

	// Сохранение стола в качестве позиции.
	Pos = Obj.(base.Conteiner)

	couch, coffee := items.NewCouchAndCoffee()

	// Создание чашки кофе.
	Obj = coffee

	// Размещение чашки кофе в сохраненную позицию (на стол).
	base.Place(Obj, Pos)

	// Создание ноутбука.
	lp := (*items.Laptop).New(&items.Laptop{})
	Obj = lp

	// Создание события <победа>
	win := stEvent{}

	// Отслеживание победы
	win.check = func() bool {
		t1 := engine.InputFormat(lp.Text())
		t2 := engine.InputFormat(texts.GameText("шифр"))
		return t1 == t2
	}

	// Что делать при наступлении победы
	win.handle = func() engine.Response {
		// Предотвращает повторную обработку события (Победить можно один раз).
		win.check = func() bool {
			return false
		}

		// Переконфигурирование повествователя.
		Nr.Texts = []string{texts.GameText("победа"), texts.GameText("шифр") + "... " + texts.GameText("Мураками")}
		Nr.NumberText = 0

		// Создание заглушки конца игры после того как повествователь отработал.
		Nr.NextImplementer = (*gameEnder).New(&gameEnder{})

		// Назначение повествователя исполнителем команд.
		w.NewActiveHandler(Nr)

		// Старт повествователя.
		resp, _ := Nr.Handle("->")

		// Увеличение счетчика прохождений в базе данных.
		engine.DB.Up(w.ID, engine.DB.UpPasseds)

		return resp
	}

	// Добавление события победы в список отслеживаемых событий.
	w.Ev = append(w.Ev, &win)

	// Размещение ноутбука в сохраненную позицию (на стол).
	base.Place(Obj, Pos)

	// Сохранение встроенного ящика стола в качестве позиции.
	Obj, _ = engine.ConsumePositionerFoundByName("ящик", base.FindPosition{Where: Pos, Deep: 0, IncludWhere: false})
	Pos = Obj.(base.Conteiner)

	// Создание конверта.
	Obj = (*items.Envelope).New(&items.Envelope{})

	// Размещение конверта в сохранённую позицию (встроенный ящик стола).
	base.Place(Obj, Pos)

	// Сохранение конверта в качестве позиции.
	Pos = Obj.(base.Conteiner)

	// Создание карты памяти
	Obj = (*items.Ssd).New(&items.Ssd{})

	// Размемещениее карты памяти в сохраненную позицию (конверт)
	base.Place(Obj, Pos)

	// Создание кушетки.
	Obj = couch

	// Размещение кушетки в кабинете
	base.Place(Obj, w.Locations("кабинет"))

	// Создание кресла.
	Obj = (*items.Armchair).New(&items.Armchair{})

	// Размещение кресла в кабинете.
	base.Place(Obj, w.Locations("кабинет"))

	// Создание книжного шкафа.
	Obj = (*items.Bookcase).New(&items.Bookcase{})

	// Размещение книжного шкафа в кабинете.
	base.Place(Obj, w.Locations("кабинет"))

	// Сохранение книжного шкафа в качестве позиции.
	Pos = Obj.(base.Conteiner)

	// Создание конспекта
	Obj = (*items.Manual).New(&items.Manual{})

	// Размемещениее конспекта в сохраненную позицию (книжный шкаф)
	base.Place(Obj, Pos)

	// Возвращает вновь созданный мир.
	return w
}
