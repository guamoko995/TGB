package world

import (
	"TelegramGameBot/Game/engine"
	mediafiles "TelegramGameBot/Game/mediaFiles"
	"fmt"
	"strings"
)

// Повествователь - игровой исполнитель, осуществляющий повествование
type Narrator struct {
	Texts           []string
	NumberText      int
	NextImplementer engine.Handler
	W               *engine.World
	StatPrefix      string
}

// Статус строка повествователя - номер текущего сообщения и общее количество сообщений.
func (n *Narrator) Status() string {
	if n.NumberText == 0 {
		return "[повествование]"
	} else {
		return fmt.Sprintf("[%v %v/%v ]", n.StatPrefix, n.NumberText, len(n.Texts))
	}
}

// Повествователь принимает одну команду - показать следующее сообщение "->"
func (*Narrator) Options() [][]string {
	return [][]string{{"->"}}
}

// Повествователь показывает следующее сообщение и передает статус исполниетеля
// следующему исполнителю сразу после показа последнего сообщения.
func (n *Narrator) Handle(str string) (engine.Response, string) {
	l := len(n.Texts) - 1
	mStr := strings.SplitN(n.Texts[n.NumberText], "<img>", 2)
	t := mStr[0]
	img := ""
	if len(mStr) > 1 {
		img = mediafiles.Image[mStr[1]]
	}
	mStr = strings.SplitN(t, "<doc>", 2)
	t = mStr[0]
	doc := ""
	if len(mStr) > 1 {
		doc = mediafiles.Doc[mStr[1]]
	}
	if n.NumberText < l {
		n.NumberText++
		return engine.Response{
			Doc:     doc,
			Img:     img,
			Msg:     t,
			Status:  n.Status(),
			Options: n.Options(),
		}, ""
	} else {
		n.W.ActiveHandler = n.NextImplementer
		return engine.Response{
			Doc:     doc,
			Img:     img,
			Msg:     t,
			Status:  n.W.ActiveHandler.Status(),
			Options: n.W.ActiveHandler.Options(),
		}, ""
	}
}

type gameEnder struct {
	status string
	msg    string
}

func (*gameEnder) Options() [][]string {
	return [][]string{{"сотрудничать", "быть уничтоженным"}}
}

func (ge *gameEnder) Status() string {
	return ge.status
}

func (ge *gameEnder) Handle(request string) (engine.Response, string) {
	ge.status = "[конец ознакомительного фрагмента]"
	msg := ge.msg
	ge.msg = "Вы можете начать игру заново по команде /start"
	return engine.Response{
		Img:     mediafiles.Image["Juno finish"],
		Msg:     msg,
		Status:  ge.status,
		Options: [][]string{},
	}, ""
}

func (ge *gameEnder) New() *gameEnder {
	ge.status = "[принятие решения]"
	ge.msg = "Ваш ответ поступил в очередь на обработку"
	return ge
}

type stEvent struct {
	check  func() bool
	handle func() engine.Response
}

func (ev *stEvent) Check() bool {
	return ev.check()
}
func (ev *stEvent) Handle() engine.Response {
	return ev.handle()
}
