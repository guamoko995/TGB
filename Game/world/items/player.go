package items

import (
	"TelegramGameBot/Game/base"
	"TelegramGameBot/Game/engine"
	"fmt"
	"strings"
)

type Box interface {
	base.Positioner
	base.Conteiner
}

type Inventory struct {
	*base.StNamer
	*base.StPositioner
	*base.StLimitedConteiner
}

func (b *Inventory) New() *Inventory {
	b = &Inventory{
		StNamer:            (*base.StNamer).New(&base.StNamer{}),
		StPositioner:       &base.StPositioner{},
		StLimitedConteiner: (*base.StLimitedConteiner).New(&base.StLimitedConteiner{}),
	}
	b.Recapacity(1000)
	return b
}

// По умолчанию строковым представлением коробки является перечисление
// содержимого.
func (b *Inventory) String() string {
	return engine.StrConteiner(b)
}

type Player struct {
	*base.StPositioner
	*engine.TreeHandlers
	Inv   Box
	Focus base.Conteiner
}

// Конструктор.
func (pl *Player) New() *Player {
	pl = &Player{
		StPositioner: &base.StPositioner{},
		TreeHandlers: (*engine.TreeHandlers).New(&engine.TreeHandlers{}),
	}
	pl.Stat = func() string {
		if pl.Focus == pl.Position() {
			return "[ ]"
		}
		return "[осмотр " + pl.Focus.Name("Р") + "]"
	}
	oldOpt := pl.Opt
	pl.Opt = func() [][]string {
		options := oldOpt()[0]
		if pl.Focus == pl.Position() {
			for i, option := range options {
				if option == "х" {
					l := len(options) - 1
					options[i] = options[l]
					options = options[:l]
					break
				}
			}
		}
		return [][]string{options}
	}
	pl.invInit()
	pl.CommandInit()
	return pl
}

// Проброс base.Conteiner метода инвентаря.
func (pl *Player) Content() []base.Positioner {
	return pl.Inv.Content()
}

// Проброс base.Conteiner метода инвентаря.
/*func (pl *Player) Capacity() int {
	return pl.Inv.Capacity()
}*/

// Проброс base.Conteiner метода инвентаря.
/*func (pl *Player) ReCapacity(cap int) {
	pl.Inv.ReCapacity(cap)
}*/

// Проброс base.Conteiner метода инвентаря.
/*func (pl *Player) Vacancy() int {
	return pl.Inv.Vacancy()
}*/

// Проброс base.Conteiner метода инвентаря.
/*func (pl *Player) Taking(obj base.Positioner) error {
	return pl.Inv.Taking(obj)
}

// Проброс base.Conteiner метода инвентаря.
func (pl *Player) Giving(obj base.Positioner) error {
	return pl.Inv.Giving(obj)
}*/

// Проброс base.Conteiner метода инвентаря.
func (pl *Player) Remove(obj ...base.Positioner) {
	pl.Inv.Remove(obj...)
}

// Проброс base.Conteiner метода инвентаря.
func (pl *Player) Add(obj ...base.Positioner) {
	pl.Inv.Add(obj...)
}

// Инициализация инвентаря.
func (pl *Player) invInit() {
	pl.Inv = (*Inventory).New(&Inventory{})
	pl.Inv.AddName([]base.Title{
		{Form: "", Value: "инвентарь"},
		{Form: "Р", Value: "инвентаря"},
		{Form: "Д", Value: "инвентарю"},
		{Form: "В", Value: "инвентарь"},
		{Form: "Т", Value: "инвентарем"},
		{Form: "П", Value: "инвентаре"},
		{Form: "где", Value: "в инвентаре"},
		{Form: "откуда", Value: "из инвентаря"},
		{Form: "куда", Value: "в инвентарь"},
		{Form: "род", Value: "м"},
	}...)
	pl.Inv.AddInfo(base.Title{Form: "род", Value: "м"})
	pl.Inv.Reposition(pl)
}

// Формирует стандартный ответ из сообщения.
func (pl *Player) StResp(Msg string) engine.Response {
	return engine.Response{
		Msg:     Msg,
		Status:  pl.Status(),
		Options: pl.Options(),
	}
}

// Инициализация игровых команд.
func (pl *Player) CommandInit() {
	pl.Applications = map[string]engine.Handler{

		// Выходит из режима осмотра, приравнивая фокус игрока к его
		// позиции.
		"х": engine.PrimalHandlers(func(args string) (engine.Response, string) {
			if pl.Focus == pl.Position() {
				return pl.StResp(""), args
			}
			if FocusPos, ok := pl.Focus.(base.Positioner); ok {
				pl.Focus = FocusPos.Position()
				if pl.Focus == nil {
					pl.Focus = pl.Position()
				}
				if pl.Focus == pl.Inv {
					pl.Focus = pl.Position()
				}
			} else {
				pl.Focus = pl.Position()
			}
			return pl.StResp(""), args
		}),

		// Перечисляет содержимое контейнера в фокусе игрока.
		/*"содержимое": func(args string) (string, string) {
			cont := pl.Focus.Content()
			if len(cont) > 0 {
				objNames := make([]string, 0)
				for _, obj := range cont {
					if obj != base.Positioner(pl) {
						objNames = append(objNames, obj.Name())
					}
				}
				if len(objNames) > 0 {
					return pl.Focus.Name("где") + " есть " + List(objNames...), args
				}
			}
			o := engine.Postfix(pl.Focus)
			return pl.Focus.Name() + " пуст" + o, args
		},*/

		// Перемещает игрока в указанную локацию при наличии прохода.
		"идти": engine.PrimalHandlers(func(where string) (engine.Response, string) {
			W := engine.RootConteiner(pl)

			distanation, endStr, err := pl.ConsumeBridgeFoundByName(where)
			if err != nil {
				if err == engine.ErrEmptyStr {
					//err = fmt.Errorf("Не указано куда именно вы хотите идти")
					return W.NewActiveHandler(&engine.ObjComplementer{
						W:               W,
						LastCommand:     "идти",
						NextImplementer: pl,
						Exceptions:      []string{},
						Fp: []base.FindPosition{
							{
								Where:       (pl.Position()).(*engine.Location).Bridge,
								Deep:        0,
								IncludWhere: false,
							},
						},
						OptionForm: "куда",
					}), ""
				}
				return pl.StResp(err.Error()), ""
			}
			err = base.Place(pl, distanation)
			if err != nil {
				return pl.StResp(err.Error()), ""
			}
			return pl.StResp("Вы пришли " + pl.Position().Name("куда")), endStr
		}),

		// Собщает местоположение игрока, перечисляет содержимое локации
		// и проходы.
		"осмотреться": engine.PrimalHandlers(func(args string) (engine.Response, string) {
			pl.Focus = pl.Position()
			s := pl.Focus.String()
			ms := strings.SplitN(s, ". ", 2)
			ms[0] = "Вы находитесь " + pl.Position().Name("где")
			return pl.StResp(ms[0] + ". " + ms[1]), args
		}),

		// Описывает предмет. В случае если предмет является контейнером
		// "помещает" его в фокус игрока.
		"осмотреть": engine.PrimalHandlers(func(args string) (engine.Response, string) {
			W := engine.RootConteiner(pl)
			obj, endStr, err := W.ConsumePositionerFoundByName(
				args,
				base.FindPosition{Where: pl.Focus, Deep: 0, IncludWhere: false},
				base.FindPosition{Where: pl.Inv, Deep: 0, IncludWhere: false},
			)
			if err != nil {
				if err == engine.ErrEmptyStr {
					//err = fmt.Errorf("Не указано что именно вы хотите осмотреть")
					return W.NewActiveHandler(&engine.ObjComplementer{
						W:               W,
						LastCommand:     "осмотреть",
						NextImplementer: pl,
						Exceptions:      []string{pl.Name("В")},
						Fp: []base.FindPosition{
							{
								Where:       pl.Focus,
								Deep:        0,
								IncludWhere: false,
							},
							{
								Where:       pl.Inv,
								Deep:        0,
								IncludWhere: false,
							},
						},
						OptionForm: "В",
					}), ""
				}
				return pl.StResp(err.Error()), ""
			}
			if FocusCont, ok := obj.(base.Conteiner); ok {
				if len(FocusCont.Content()) > 0 {
					pl.Focus = FocusCont
				}
			}
			return pl.StResp(obj.String()), endStr
		}),

		// Список всех локаций
		/*"карта": func(args string) (string, string) {
			W := engine.RootConteiner(pl)
			return W.String(), args
		},*/

		// Перечисляет содержимое инвентаря.
		"инвентарь": engine.PrimalHandlers(func(args string) (engine.Response, string) {
			return pl.StResp(pl.Inv.String()), args
		}),

		// Перемещает указанный предмет в инвентарь.
		"взять": engine.PrimalHandlers(func(args string) (engine.Response, string) {
			W := engine.RootConteiner(pl)
			obj, endStr, err := W.ConsumePositionerFoundByName(
				args,
				base.FindPosition{Where: pl.Focus, Deep: 0, IncludWhere: false},
			)
			if err != nil {
				if err == engine.ErrEmptyStr {
					//err = fmt.Errorf("Не указано что именно вы хотите взять")
					return W.NewActiveHandler(&engine.ObjComplementer{
						W:               W,
						LastCommand:     "взять",
						NextImplementer: pl,
						Exceptions:      []string{pl.Name("В")},
						Fp: []base.FindPosition{
							{
								Where:       pl.Focus,
								Deep:        0,
								IncludWhere: false,
							},
						},
						OptionForm: "В",
					}), ""
				}
				return pl.StResp(err.Error()), ""
			}
			return pl.StResp(engine.Place(obj, pl.Inv)), endStr
		}),

		// Перемещает указанный предмет из инвентаря в указаный
		// контейнер. Если контейнер не указан, то перемещает в
		// контейнер, находящийся в фокусе игрока.
		"разместить": engine.PrimalHandlers(func(args string) (engine.Response, string) {
			W := engine.RootConteiner(pl)
			obj, endStr, err := W.ConsumePositionerFoundByName(
				args,
				base.FindPosition{Where: pl.Inv, Deep: 0, IncludWhere: false})
			if err != nil {
				if err == engine.ErrEmptyStr {
					//err = fmt.Errorf("Не указано что именно вы хотите разместить")
					return W.NewActiveHandler(&engine.ObjComplementer{
						W:               W,
						LastCommand:     "разместить",
						NextImplementer: pl,
						Exceptions:      []string{pl.Name("В")},
						Fp: []base.FindPosition{
							{
								Where:       pl.Inv,
								Deep:        0,
								IncludWhere: false,
							},
						},
						OptionForm: "В",
					}), ""
				}
				return pl.StResp(err.Error()), ""
			}
			c, endStr, err := W.ConsumePositionerFoundByName(
				endStr,
				base.FindPosition{Where: pl.Focus, Deep: 0, IncludWhere: true},
				base.FindPosition{Where: pl.Inv, Deep: 0, IncludWhere: false},
			)
			if err != nil {
				if err == engine.ErrEmptyStr {
					//err = fmt.Errorf("Не указано где именно вы хотите разместить")
					return W.NewActiveHandler(&engine.ObjComplementer{
						W:               W,
						LastCommand:     "разместить " + obj.Name("В"),
						NextImplementer: pl,
						Exceptions:      []string{pl.Name("где"), obj.Name("где")},
						Fp: []base.FindPosition{
							{
								Where:       pl.Focus,
								Deep:        0,
								IncludWhere: true,
							},
							{
								Where:       pl.Inv,
								Deep:        0,
								IncludWhere: false,
							},
						},
						OptionForm: "где",
					}), ""
				}
				return pl.StResp(err.Error()), ""
			}
			cont, ok := c.(base.Conteiner)
			if !ok {
				return pl.StResp(c.Name() + " не может содержать предметы"), endStr
			}
			return pl.StResp(engine.Place(obj, cont)), endStr
		}),

		// Возвращает список возможностей предмета.
		"воспользоваться": engine.PrimalHandlers(func(args string) (engine.Response, string) {
			W := engine.RootConteiner(pl)
			obj, endStr, err := W.ConsumePositionerFoundByName(
				args,
				base.FindPosition{Where: pl.Focus, Deep: 0, IncludWhere: false},
				base.FindPosition{Where: pl.Inv, Deep: 0, IncludWhere: false},
			)
			if err != nil {
				if err == engine.ErrEmptyStr {
					//err = fmt.Errorf("Не указано чем именно вы хотите воспользоваться")
					return W.NewActiveHandler(&engine.ObjComplementer{
						W:               W,
						LastCommand:     "воспользоваться",
						NextImplementer: pl,
						Exceptions:      []string{pl.Name("Т")},
						Fp: []base.FindPosition{
							{
								Where:       pl.Focus,
								Deep:        0,
								IncludWhere: false,
							},
							{
								Where:       pl.Inv,
								Deep:        0,
								IncludWhere: false,
							},
						},
						OptionForm: "Т",
					}), ""
				}
				return pl.StResp(err.Error()), ""
			}
			us, ok := obj.(engine.Handler)
			if !ok {
				return pl.StResp(obj.Name("Т") + " нельзя воспользоваться"), ""
			}
			options := us.Options()[0]
			if len(options) == 0 {
				return pl.StResp(obj.Name("Т") + " нельзя воспользоваться"), ""
			}
			for i, option := range options {
				if option == "х" {
					l := len(options) - 1
					options[i] = options[l]
					options = options[:l]
					break
				}
			}
			resp := W.NewActiveHandler(us)
			resp.Msg = "вы можете " + engine.List(options...)
			return resp, endStr
		}),
	}
}

// Обертка base.Positioner метода для привязки фокуса игрока к его позиции при ее
// изменении.
func (pl *Player) Reposition(c base.Conteiner) {
	pl.StPositioner.Reposition(c)
	pl.Focus = pl.Position()
}

// Ищет доступный проход по началу строки, вторым аргументом возвращает
// конец строки (без считанного прохода).
func (pl *Player) ConsumeBridgeFoundByName(name string) (*engine.Location, string, error) {
	W := engine.RootConteiner(pl)
	if name == "" {
		return nil, "", engine.ErrEmptyStr
	}
	obj, endStr := engine.ConsumePositionerFoundByName(
		name,
		base.FindPosition{Where: pl.Position().(*engine.Location).Bridge, Deep: 0, IncludWhere: false},
	)
	if obj != nil {
		return obj.(*engine.Location), endStr, nil
	}
	obj, _ = engine.ConsumePositionerFoundByName(
		name,
		base.FindPosition{Where: W, Deep: 0, IncludWhere: false},
	)
	if obj != nil {
		return nil, "", fmt.Errorf(pl.Position().Name("откуда") + " нельзя пройти " + obj.Name("куда"))
	}
	words := strings.Split(name, " ")
	var s string
	for i, w := range words {
		s += w
		if !engine.FindObjByPartName(s, base.FindPosition{Where: W, Deep: 0, IncludWhere: false}) {
			var postfix string
			if i == 0 {
				postfix = "a"
			}
			return nil, "", fmt.Errorf("нельзя идти \"" + s + "\" или в место начинающееся со слов" + postfix + " \"" + s + "\"")
		}
	}
	return nil, "", fmt.Errorf("вероятно Вы что-то недоговариваете")
}
