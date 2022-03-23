package items

import (
	"TelegramGameBot/Game/base"
	"TelegramGameBot/Game/engine"
	"TelegramGameBot/Game/world/buildTools"
	"TelegramGameBot/Game/world/texts"
)

// Предмет: кушетка.
type Couch struct {
	*base.StPositioner
	*base.StSizer
	*engine.TreeHandlers
	about string
}

func (b *Couch) New() *Couch {
	b = &Couch{
		StPositioner: &base.StPositioner{},
		StSizer:      &base.StSizer{},
		TreeHandlers: (*engine.TreeHandlers).New(&engine.TreeHandlers{}),
	}

	b.Stat = func() string {
		return "[использование кушетки]"
	}

	buildTools.SetName(b, "кушетка")
	b.Resize(2500)
	b.about = "Хорошее место чтобы немного вздремнуть. Интересно, что " +
		"приснится, если заснуть в воображаемом кабинете"

	apdate := func() {
		b.Applications["лечь спать"] = engine.PrimalHandlers(func(args string) (engine.Response, string) {
			W := engine.RootConteiner(b)
			resp := W.NewActiveHandler(W.Pl)
			resp.Msg = "Вы уже выспались. Слишком много спать вредно"
			return resp, args
		})
		b.about = "Просто кушетка"
	}

	b.Applications["лечь спать"] = engine.PrimalHandlers(func(args string) (engine.Response, string) {
		W := engine.RootConteiner(b)
		defer apdate()
		resp := W.NewActiveHandler(W.Pl)
		resp.Msg = texts.GameText("сон")
		return resp, args
	})
	b.Applications["х"] = engine.PrimalHandlers(func(args string) (engine.Response, string) {
		W := engine.RootConteiner(b)
		resp := W.NewActiveHandler(W.Pl)
		return resp, args
	})
	return b
}

func (b *Couch) String() string {
	return b.about
}

// Предмет: чашка кофе.
type Coffee struct {
	*base.StPositioner
	*base.StSizer
	*engine.TreeHandlers
	about string
}

func (b *Coffee) New() *Coffee {
	b = &Coffee{
		StPositioner: &base.StPositioner{},
		StSizer:      &base.StSizer{},
		TreeHandlers: (*engine.TreeHandlers).New(&engine.TreeHandlers{}),
	}

	b.Stat = func() string {
		return "[использование чашки кофе]"
	}

	buildTools.SetName(b, "чашка кофе")
	b.Resize(20)
	b.about = "Горячий кофе, то что нужно чтобы взбодриться"

	apdate := func() {
		b.Applications = map[string]engine.Handler{}
		b.about = "Здесь остатки кофе. Кофейная гуща предсказывает шокирующие новости и трудный выбор"
		buildTools.SetName(b, "чашка")
	}

	b.Applications["выпить"] = engine.PrimalHandlers(func(args string) (engine.Response, string) {
		W := engine.RootConteiner(b)
		defer apdate()
		resp := W.NewActiveHandler(W.Pl)
		resp.Msg = texts.GameText("кофе")
		return resp, args
	})
	b.Applications["х"] = engine.PrimalHandlers(func(args string) (engine.Response, string) {
		W := engine.RootConteiner(b)
		resp := W.NewActiveHandler(W.Pl)
		return resp, args
	})
	return b
}

func (b *Coffee) String() string {
	return b.about
}
