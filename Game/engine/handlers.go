package engine

import (
	"TelegramGameBot/Game/base"
	"strings"
)

// Обработчик.
type Handler interface {
	Handle(request string) (Response, string)
	Status() string
	Options() [][]string
}

// Ответ.
type Response struct {
	Msg     string     // Основное сообщение - реакциия на запрос.
	Status  string     // Доп сообщение для отображения режима игры.
	Options [][]string // Кнопки клавиатуры в tlg.
	Img     string     // Изображение
}

type PrimalHandlers func(string) (Response, string)

func (h PrimalHandlers) Handle(arg string) (Response, string) {
	f := (func(arg string) (Response, string))(h)
	return f(arg)
}

func (PrimalHandlers) Options() [][]string {
	return [][]string{}
}

func (PrimalHandlers) Status() string {
	return ""
}

type TreeHandlers struct {
	*base.StNamer
	InputFormat  func(string) string
	OutputFormat func(string) string
	Applications map[string]Handler
}

func (*TreeHandlers) New() *TreeHandlers {
	h := TreeHandlers{
		StNamer:      (*base.StNamer).New(&base.StNamer{}),
		Applications: make(map[string]Handler),
		InputFormat:  InputFormat,
		OutputFormat: OutputFormat,
	}
	return &h
}

func (h *TreeHandlers) Options() [][]string {
	options := make([]string, 0)
	for key := range h.Applications {
		options = append(options, key)
	}
	return [][]string{options}
}

func (h *TreeHandlers) Status() string {
	return ""
}

func (h *TreeHandlers) Handle(request string) (Response, string) {
	request = h.InputFormat(request)

	app, remainder := ConsumeAppUs(h, request)
	if app == nil {
		words := strings.Split(request, " ")
		var s string
		for i, w := range words {
			s += w
			if !ChekPartAppUs(h, s) {
				text := make([]string, 2)
				text[0] = "у " + h.Name("Р") + " нет возможности \""
				if i == 0 {
					text[1] = "a"
				}
				return Response{
					Msg:     text[0] + s + "\" или начинающейся со слов" + text[1] + " \"" + s + "\"",
					Status:  h.Status(),
					Options: h.Options(),
				}, ""
			}
			s += " "
		}
		panic("Не найдена существующая возможность.")
	}
	resp, remainder := app.Handle(remainder)
	resp.Msg = h.OutputFormat(resp.Msg)
	return resp, remainder
}

type Complementer struct {
	W               *World
	LastCommand     string
	NextImplementer Handler
}

func (c *Complementer) Status() string {
	return "[...]"
}

func (c *Complementer) Options() [][]string {
	return [][]string{}
}

// Обрабатывает пользовательский вариант.
func (c *Complementer) Handle(args string) (Response, string) {
	c.W.ActiveHandler = c.NextImplementer
	/*if args=="х"{
		return Response{
			Msg: "",
			Status: c.W.ActiveHandler.Status(),
			Options: c.W.ActiveHandler.Options(),
		},""
	}*/
	return c.W.ActiveHandler.Handle(c.LastCommand + " " + args)
}

type ObjComplementer struct {
	W               *World
	LastCommand     string
	NextImplementer Handler
	Exceptions      []string
	Fp              []base.FindPosition
	OptionForm      string
}

func (c *ObjComplementer) Status() string {
	return "[...]"
}

func (c *ObjComplementer) Options() [][]string {
	masObj := make([]base.Positioner, 0)
	if len(c.Fp) > 0 {
		masObj = c.Fp[0].Where.Content()

		// При необходимости ввключает место поиска в варианты.
		if c.Fp[0].IncludWhere {
			if obj, ok := c.Fp[0].Where.(base.Positioner); ok {
				masObj = append(masObj, obj)
			}
		}
	}

	options := make([]string, 0)
	// Включает содержимое места поиска в варианты, игнорируя исключения.
	for _, br := range masObj {
		name := br.Name(c.OptionForm)

		for _, exc := range c.Exceptions {
			if name == exc {
				name = ""
				break
			}
		}
		if name != "" {
			options = append(options, name)
		}
	}

	// Включает вложенный вариант для переключения на следующее место
	// поиска.
	if len(c.Fp) > 1 {
		options = append(options, c.Fp[1].Where.Name("откуда")+"...")
	}

	// Ввключает прерывание в варианты.
	options = append(options, "х")
	return [][]string{options}
}

// Обрабатывает пользовательский вариант.
func (c *ObjComplementer) Handle(args string) (Response, string) {
	if len(c.Fp) > 1 {
		if args == c.Fp[1].Where.Name("откуда")+"..." {
			c.Fp = c.Fp[1:]
			return Response{
				Msg:     "",
				Status:  c.W.ActiveHandler.Status(),
				Options: c.W.ActiveHandler.Options(),
			}, ""
		}
	}
	c.W.ActiveHandler = c.NextImplementer
	if args == "х" {
		return Response{
			Msg:     "",
			Status:  c.W.ActiveHandler.Status(),
			Options: c.W.ActiveHandler.Options(),
		}, ""
	}
	return c.W.ActiveHandler.Handle(c.LastCommand + " " + args)
}
