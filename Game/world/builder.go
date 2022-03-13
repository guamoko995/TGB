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
	buildTools.SetName(w.Pl, "игрок")

	// Создание повествователя.
	w.Nr.Texts = strings.Split(texts.GameText("введение"), "\n")

	// После окончания повествования, повествователь передает роль
	// исполнителя команд игроку.
	w.Nr.NextImplementer = w.Pl

	// Назначение повествователя исполнителем команд.
	w.NewActiveHandler(w.Nr)

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

	// Создание прохода из никуда в кабинет.
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

	// Создание ноутбука.
	Obj = (*items.Laptop).New(&items.Laptop{})

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
	Obj = (*items.Couch).New(&items.Couch{})

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

	// Создание инструкции
	Obj = (*items.Manual).New(&items.Manual{})

	// Размемещениее инструкции в сохраненную позицию (книжный шкаф)
	base.Place(Obj, Pos)

	// Возвращает вновь созданный мир.
	return w
}
