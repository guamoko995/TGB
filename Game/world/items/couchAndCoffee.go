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

func NewCouchAndCoffee() (*Couch, *Coffee) {
	couch := (*Couch).New(&Couch{})
	coffee := (*Coffee).New(&Coffee{})
	dr := coffee.Applications["выпить кофе"]
	coffee.Applications["выпить кофе"] = engine.PrimalHandlers(func(args string) (engine.Response, string) {
		resp, str := dr.Handle(args)
		couch.about = "Хорошее место для дневного сна."
		couch.Applications["лечь спать"] = engine.PrimalHandlers(func(args string) (engine.Response, string) {
			W := engine.RootConteiner(couch)
			resp := W.NewActiveHandler(W.Pl)
			resp.Msg = "Вы представляли, что пьете кофе, а теперь у Вас не получается представить, что вы спите. " +
				"Наблюдается ярко выраженное психосоматическое явление."
			return resp, args
		})
		return resp, str
	})

	sl := couch.Applications["лечь спать"]
	couch.Applications["лечь спать"] = engine.PrimalHandlers(func(args string) (engine.Response, string) {
		resp, str := sl.Handle(args)
		coffee.about = "Холодный кофе - то что нужно чтобы взбодриться"
		coffee.Applications["выпить кофе"] = engine.PrimalHandlers(func(args string) (engine.Response, string) {
			resp, str := dr.Handle(args)
			resp.Msg = texts.GameText("холодный кофе")
			return resp, str
		})
		return resp, str
	})
	return couch, coffee
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
	b.about = "Горячий кофе - то что нужно чтобы взбодриться"

	apdate := func() {
		b.Applications = map[string]engine.Handler{}
		b.about = "Здесь остатки кофе. Кофейная гуща предсказывает шокирующие новости и трудный выбор"
		buildTools.SetName(b, "чашка")
	}

	b.Applications["выпить кофе"] = engine.PrimalHandlers(func(args string) (engine.Response, string) {
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
