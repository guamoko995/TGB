package items

import (
	"TelegramGameBot/Game/base"
	"TelegramGameBot/Game/engine"
	"TelegramGameBot/Game/world/buildTools"
	"fmt"
)

// Предмет: стол с выдвижным ящиком.
type TableWithDrawer struct {
	*base.StNamer
	*base.StPositioner
	*base.StSizer
	*base.StLimitedConteiner
}

func (b *TableWithDrawer) New() *TableWithDrawer {
	b = &TableWithDrawer{
		StNamer:            (*base.StNamer).New(&base.StNamer{}),
		StLimitedConteiner: (*base.StLimitedConteiner).New(&base.StLimitedConteiner{}),
		StPositioner:       &base.StPositioner{},
		StSizer:            &base.StSizer{},
	}
	buildTools.SetName(b, "стол")
	b.Resize(2500)
	b.Recapacity(2000)

	//var Obj base.Positioner

	Obj := (*Drawer).New(&Drawer{})

	base.Place(Obj, b)
	Obj.NoRePlace = fmt.Errorf("выдвижной ящик является частью стола и не может быть перемещен")
	return b
}

// Обособленное представление выдвижного ящика как части стола.
func (c *TableWithDrawer) String() string {
	s := "в столе есть выдвижной ящик"
	ms := make([]string, 0)
	for _, obj := range c.Content() {
		if ok, _ := engine.ConsumeNameObj(obj, "выдвижной ящик"); !ok {
			ms = append(ms, obj.Name())
		}
	}
	if len(ms) > 0 {
		return s + ". " + c.Name("где") + " " + engine.List(ms...)
	}
	return s
}
