package items

import (
	"TelegramGameBot/Game/base"
	"TelegramGameBot/Game/engine"
	"TelegramGameBot/Game/world/buildTools"
	worldGame "TelegramGameBot/Game/world/items/wordGame"
	"TelegramGameBot/Game/world/texts"
	"fmt"
	"strings"
)

// Предмет: ноутбук.
type Laptop struct {
	*base.StPositioner
	*base.StSizer
	*base.StLimitedConteiner
	*engine.TreeHandlers
	te       *textEditor
	wq       *SlotMakhine
	freqSort bool
}

func (b *Laptop) Status() string {
	return "[использование ноутбука]"
}

func (b *Laptop) New() *Laptop {
	b = &Laptop{
		StPositioner:       &base.StPositioner{},
		StSizer:            &base.StSizer{},
		StLimitedConteiner: (*base.StLimitedConteiner).New(&base.StLimitedConteiner{}),
		TreeHandlers:       (*engine.TreeHandlers).New(&engine.TreeHandlers{}),
		te:                 (*textEditor).New(&textEditor{}),
		wq: &SlotMakhine{
			SlotMakhine: &worldGame.SlotMakhine{},
		},
	}

	b.wq.SetText(worldGame.MQT([]worldGame.QwestText{
		worldGame.QwestText(worldGame.NewQText(texts.GameText("шифр"))),
		worldGame.QwestText(murakami()),
	}))

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
	b.Applications["запустить текстовый редактор"] = engine.PrimalHandlers(func(args string) (engine.Response, string) {
		resp := engine.Response{
			Msg:     "",
			Status:  b.Status(),
			Options: b.Options(),
		}
		if len(b.Content()) > 0 {
			W := engine.RootConteiner(b)
			b.wq.W = W
			if b.freqSort {
				resp = W.NewActiveHandler(b.wq)
			} else {
				resp = W.NewActiveHandler(b.te)
			}

			resp.Msg = "На карте памяти Вы нашли текстовый файл с именем \"Ш" +
				"ифр\" и открыли его в текстовом редакторе. Вам доступны" +
				"следующие возможности (чувствительны к регистру):\n" +
				"показать,\n" +
				"заменить <символ> на <символ>,\n" +
				"нижний  регистр,\n" +
				"верхний регистр,\n" +
				"количество <символ>,\n" +
				"заново"
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

	b.wq.NextHandler = b

	b.te.InputFormat = func(s string) string { return s }
	b.te.OutputFormat = func(s string) string { return s }
	b.te.Applications["показать"] = engine.PrimalHandlers(func(args string) (engine.Response, string) {
		return b.te.StResp(b.te.Text.Print()), args
	})
	b.te.Applications["заменить"] = engine.PrimalHandlers(func(args string) (engine.Response, string) {
		W := engine.RootConteiner(b)
		if args == "" {
			resp := W.NewActiveHandler(&engine.Complementer{
				W:               W,
				LastCommand:     "заменить",
				NextImplementer: b.te,
			})
			return resp, ""
		}
		ar := strings.SplitN(args, " ", 3)
		if len(ar) < 2 {
			return b.te.StResp("Команда \"заменить\" принимает два символа, разделенные пробелом: первый - заменяемый, второй - заменяющий"), ""
		}
		endStr := ""
		if len([]rune(ar[0])) != 1 || len([]rune(ar[1])) != 1 {
			return b.te.StResp("Команда \"заменить\" принимает два символа, разделенные пробелом: первый - заменяемый, второй - заменяющий"), ""
		}
		if len(ar) > 2 {
			endStr = ar[3]
		}
		b.te.Text.Replace([]rune(ar[0])[0], []rune(ar[1])[0])
		return b.te.StResp("Символ '" + ar[0] + "' заменен на символ '" + ar[1] + "'"), endStr
	})
	b.te.Applications["нижний регистр"] = engine.PrimalHandlers(func(args string) (engine.Response, string) {
		b.te.Text.Down()
		return b.te.StResp("текст переведен в нижжний регистр"), args
	})
	b.te.Applications["верхний регистр"] = engine.PrimalHandlers(func(args string) (engine.Response, string) {
		b.te.Text.Up()
		return b.te.StResp("текст переведен в верхний регистр"), args
	})
	b.te.Applications["посчитать"] = engine.PrimalHandlers(func(args string) (engine.Response, string) {
		W := engine.RootConteiner(b)
		if args == "" {
			resp := W.NewActiveHandler(&engine.Complementer{
				W:               W,
				LastCommand:     "посчитать",
				NextImplementer: b.te,
			})
			return resp, ""
		}
		ar := strings.SplitN(args, " ", 2)
		endStr := ""
		if len([]rune(ar[0])) != 1 {
			return b.te.StResp("Команда \"посчитать\" принимает один символ"), ""
		}
		if len(ar) > 2 {
			endStr = ar[2]
		}
		s := []rune(ar[0])[0]
		n := b.te.Text.Count(s)
		all := b.te.Text.CountAll()
		p := float32(100*n) / float32(all)
		return b.te.StResp(fmt.Sprintf("Символ '%c' встречается в тексте %v раз, что составляет %.2f%% от общего количества символов в тексте", s, n, p)), endStr
	})
	b.te.Applications["заново"] = engine.PrimalHandlers(func(args string) (engine.Response, string) {
		b.te.Text.Reset()
		return b.te.StResp("Все измененияя отменены"), args
	})
	b.te.Applications["х"] = engine.PrimalHandlers(func(args string) (engine.Response, string) {
		W := engine.RootConteiner(b)
		return W.NewActiveHandler(b), args
	})
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

// Псевдопредмет: текстовый редактор
type textEditor struct {
	*engine.TreeHandlers
	Text worldGame.QwestText
}

var varMurakami = worldGame.NewPsevdoText("Game/mediaFiles/Murakami.txt")

func murakami() worldGame.PsevdoText {
	copyText := worldGame.PsevdoText(make([]worldGame.RuneCount, len([]worldGame.RuneCount(varMurakami))))
	copy([]worldGame.RuneCount(copyText), []worldGame.RuneCount(varMurakami))
	return varMurakami
}

func (te *textEditor) New() *textEditor {
	te = &textEditor{
		TreeHandlers: (*engine.TreeHandlers).New(&engine.TreeHandlers{}),
	}
	buildTools.SetName(te, "текстовый редактор")

	// Текст состоит из двух текстов. Первый - не большой реально
	// отображаемый, второй - не отображаемый псевдотекст (карта,
	// ключи которой составляют русский алфавит, а значения
	// соответствуют количеству символов, соответствующих ключу в
	// тексте.
	te.Text = worldGame.MQT([]worldGame.QwestText{
		worldGame.QwestText(worldGame.NewQText(texts.GameText("шифр"))),
		worldGame.QwestText(murakami()),
	})

	crMap := worldGame.GenCryptMap()

	te.Text.Crypt(crMap)
	fmt.Printf("%c\n", crMap)
	decrMap := worldGame.DecryptMap(te.Text)
	te.Text.Crypt(decrMap)
	fmt.Printf("%c\n", decrMap)
	return te
}

func (te *textEditor) StResp(Msg string) engine.Response {
	return engine.Response{
		Msg:     Msg,
		Status:  te.Status(),
		Options: te.Options(),
	}
}

func (te *textEditor) Status() string {
	return "[использование текстового редактора]"
}

// Машина слотов.
type SlotMakhine struct {
	*worldGame.SlotMakhine
	NextHandler engine.Handler
	W           *engine.World
}

// Возвращает варианты взаимодействия с машиной слотов.
func (sm *SlotMakhine) Options() [][]string {
	options := [][]string{
		{},
		{"<-", "зафиксировать слово", "->"},
		{"заново", "x"},
	}
	for _, mr := range sm.Str {
		R := mr.Name()
		sR := string([]rune{R})
		options[0] = append(options[0], sR)
	}
	return options
}

func (sm *SlotMakhine) Status() string {
	return "[дешифровка]"
}

func (sm *SlotMakhine) Handle(str string) (engine.Response, string) {
	l := len(sm.Words)
	switch str {
	case "<-":
		sm.Pos--
		if sm.Pos < 0 {
			sm.Pos += l
		}
		sm.Update()
	case "->":
		sm.Pos++
		if sm.Pos == l {
			sm.Pos = 0
		}
		sm.Update()
	case "заново":
		sm.Text.Reset()
		sm.Text.Down()
		sm.Update()
	case "зафиксировать слово":
		for key, R := range sm.ReplaceMap() {
			sm.Text.Replace(key, R)
		}
		sm.Update()
	case "x":
		return sm.W.NewActiveHandler(sm.NextHandler), ""
	default:
		R := []rune(str)[0]
		sm.SmartClick(R)
	}
	msg := make([]rune, 0)
	for _, pos := range sm.Str {
		msg = append(msg, pos.Name())
	}
	return engine.Response{
		Msg:     string(msg),
		Status:  "[клавиатура]",
		Options: sm.Options(),
	}, ""
}
