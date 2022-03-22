package items

import (
	"TelegramGameBot/Game/base"
	"TelegramGameBot/Game/engine"
	"TelegramGameBot/Game/world/buildTools"
	inLaptop "TelegramGameBot/Game/world/items/inLaptop"
	worldGame "TelegramGameBot/Game/world/items/inLaptop/wordGame"
	"TelegramGameBot/Game/world/texts"
	"fmt"
)

var Murakami = worldGame.NewPsevdoText("Game/mediaFiles/Murakami.txt")

// Предмет: ноутбук.
type Laptop struct {
	*base.StPositioner
	*base.StSizer
	*base.StLimitedConteiner
	*engine.TreeHandlers
	te     *inLaptop.TextEditor
	sm     *inLaptop.SlotMakhine
	useMan bool
}

func (b *Laptop) Text() string {
	return b.te.Text.Print()
}

func (b *Laptop) New() *Laptop {
	b = &Laptop{
		StPositioner:       &base.StPositioner{},
		StSizer:            &base.StSizer{},
		StLimitedConteiner: (*base.StLimitedConteiner).New(&base.StLimitedConteiner{}),
		TreeHandlers:       (*engine.TreeHandlers).New(&engine.TreeHandlers{}),
		te:                 (*inLaptop.TextEditor).New(&inLaptop.TextEditor{}),
		sm: &inLaptop.SlotMakhine{
			SlotMakhine: &worldGame.SlotMakhine{},
		},
	}

	b.Stat = func() string {
		return "[использование ноутбука]"
	}

	// Текст состоит из двух текстов. Первый - не большой реально
	// отображаемый, второй - не отображаемый псевдотекст (карта,
	// ключи которой составляют русский алфавит, а значения
	// соответствуют количеству символов, соответствующих ключу в
	// тексте).
	copyText := worldGame.PsevdoText(make([]worldGame.RuneCount, len([]worldGame.RuneCount(Murakami))))
	copy([]worldGame.RuneCount(copyText), []worldGame.RuneCount(Murakami))
	text := worldGame.MQT([]worldGame.QwestText{
		worldGame.QwestText(worldGame.NewQText(texts.GameText("шифр"))),
		worldGame.QwestText(copyText),
	})

	// Текст зашифрован случайной заменой.
	crMap := worldGame.GenCryptMap()
	text.Crypt(crMap)

	b.te.Text = text

	buildTools.SetName(b, "ноутбук")
	b.Resize(3)
	b.Recapacity(3)

	b.Applications["подключить оборудование"] = engine.PrimalHandlers(func(args string) (engine.Response, string) {
		resp := engine.Response{
			Msg:     "",
			Status:  b.Status(),
			Options: b.Options(),
		}
		W := engine.RootConteiner(b)
		obj, endStr, err := W.ConsumePositionerFoundByName(args, base.FindPosition{Where: W.Pl.(*Player).Inv, Deep: 0, IncludWhere: false})
		if err != nil {
			if err == engine.ErrEmptyStr {
				//err = fmt.Errorf("Не указано что именно вы хотите подключить")
				return W.NewActiveHandler(&engine.ObjComplementer{
					W:               W,
					LastCommand:     "подключить оборудование",
					NextImplementer: b,
					Exceptions:      []string{b.Name("В")},
					Fp: []base.FindPosition{
						{
							Where:       W.Pl.(*Player).Inv,
							Deep:        0,
							IncludWhere: false,
						},
					},
					OptionForm: "В",
				}), ""
			}
			resp.Msg = err.Error()
			return resp, ""
		}
		err = base.Place(obj, b)
		if err != nil {
			resp.Msg = err.Error()
			return resp, ""
		}
		o := engine.Postfix(obj)
		resp.Msg = obj.Name() + " подключен" + o + " к " + b.Name("Д")
		return resp, endStr
	})
	b.Applications["отключить оборудование"] = engine.PrimalHandlers(func(args string) (engine.Response, string) {
		resp := engine.Response{
			Msg:     "",
			Status:  b.Status(),
			Options: b.Options(),
		}
		W := engine.RootConteiner(b)
		obj, endStr, err := W.ConsumePositionerFoundByName(args, base.FindPosition{Where: b, Deep: 0, IncludWhere: false})
		if err != nil {
			if err == engine.ErrEmptyStr {
				//err = fmt.Errorf("Не указано что именно вы хотите отключить")
				return W.NewActiveHandler(&engine.ObjComplementer{
					W:               W,
					LastCommand:     "отключить оборудование",
					NextImplementer: b,
					Exceptions:      []string{},
					Fp: []base.FindPosition{
						{
							Where:       b,
							Deep:        0,
							IncludWhere: false,
						},
					},
					OptionForm: "В",
				}), ""
			}
			resp.Msg = err.Error()
			return resp, ""
		}
		err = base.Place(obj, W.Pl.(*Player).Inv)
		if err != nil {
			resp.Msg = err.Error()
			return resp, ""
		}
		o := engine.Postfix(obj)
		resp.Msg = obj.Name() + " отключен" + o + " от " + b.Name("Р") + " и помещен" + o + " " + W.Pl.(*Player).Inv.Name("куда")
		return resp, endStr
	})
	b.Applications["запустить редактор текcта"] = engine.PrimalHandlers(func(args string) (engine.Response, string) {
		resp := engine.Response{
			Msg:     "",
			Status:  b.Status(),
			Options: b.Options(),
		}
		if len(b.Content()) > 0 {
			W := engine.RootConteiner(b)
			b.sm.W = W
			b.te.W = W

			ok := engine.FindObjByPartName(
				"конспект по криптоанализу",
				base.FindPosition{
					Where:       W.Pl.(*Player),
					Deep:        0,
					IncludWhere: false,
				},
			)
			if ok {
				b.te.Applications["воспользоваться конспектом"] = engine.PrimalHandlers(func(args string) (engine.Response, string) {
					msg := ""
					if !b.useMan {
						msg = texts.GameText("первое использование конспекта")

						// Производит замену букв в тексте в соответствии с
						// распространенностью букв в русском языке
						b.sm.SetText(text)
						b.useMan = true
					} else {
						msg = texts.GameText("второе использование конспекта")
					}
					b.sm.Update()
					resp := W.NewActiveHandler(b.sm)
					resp.Msg = msg
					return resp, args
				})
			} else {
				delete(b.te.Applications, "воспользоваться конспектом")
			}
			resp = W.NewActiveHandler(b.te)
			//}

			resp.Msg = "На карте памяти Вы нашли текстовый файл с именем \"Ш" +
				"ифр\" и открыли его в текстовом редакторе.\n\nВам доступны " +
				"следующие возможности:\n" +
				"- показать - выводит текст файла,\n" +
				"- заменить <х1> <х2> - заменяет символ <х1> на символ <х2>,\n" +
				"- нижний  регистр - переводит весь текст в нижний регистр,\n" +
				"- верхний регистр - переводит весь текст в верхний регистр,\n" +
				"- количество <х> - подсчитывает сколько раз символ <x> встречается в тексте,\n" +
				"- заново - отменяет все совершенные замены."

			if ok {
				resp.Msg = resp.Msg + "\n\nЕще Вы можете воспользоваться конспектом по криптоанализу, который предусмотрительно держите при себе."
				return resp, args
			}
			return resp, args

		} else {
			resp.Msg = "На жестком диске нет подходящих файлов"
			return resp, args
		}
	})
	b.Applications["х"] = engine.PrimalHandlers(func(args string) (engine.Response, string) {
		W := engine.RootConteiner(b)
		return W.NewActiveHandler(W.Pl), args
	})

	b.sm.NextHandler = b
	b.te.NextHandler = b
	return b
}

func (b *Laptop) String() string {
	if engine.RootConteiner(b).Pl.(*Player).Focus == b {
		engine.RootConteiner(b).Pl.(*Player).Focus = b
	}
	s := "На ноутбуке установлен текстовый редактор. На жестком диске больше ничего интересного"
	for _, sdd := range b.Content() {
		s += ". " + sdd.(*Ssd).Inf
	}
	return s
}

// В ноутбук можно помещать только карты памяти.
func (pc *Laptop) Take(obj base.Positioner) error {
	if obj.Name() == "карта памяти" {
		return nil
	}
	return fmt.Errorf("Ноутбук не предусматривает возможность подключения " + obj.Name("Р"))
}
