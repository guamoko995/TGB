package engine

import (
	mediafiles "TelegramGameBot/Game/mediaFiles"
	"fmt"
	"strings"
)

// Повествователь - игровой исполнитель, осуществляющий повествование
type Narrator struct {
	Texts           []string
	NumberText      int
	NextImplementer Handler
	W               *World
}

// Статус строка повествователя - номер текущего сообщения и общее количество сообщений.
func (n *Narrator) Status() string {
	return fmt.Sprintf("[ %v/%v ]", n.NumberText, len(n.Texts))
}

// Повествователь принимает одну команду - показать следующее сообщение "->"
func (*Narrator) Options() [][]string {
	return [][]string{{"->"}}
}

// Повествователь показывает следующее сообщение и передает статус исполниетеля
// следующему исполнителю сразу после показа последнего сообщения.
func (n *Narrator) Handle(str string) (Response, string) {
	l := len(n.Texts) - 1
	mStr := strings.SplitN(n.Texts[n.NumberText], "<img>", 2)
	t := mStr[0]
	img := ""
	if len(mStr) > 1 {
		img = mediafiles.Image[mStr[1]]
	}
	if n.NumberText < l {
		n.NumberText++
		return Response{
			Img:     img,
			Msg:     t,
			Status:  n.Status(),
			Options: n.Options(),
		}, ""
	} else {
		n.W.ActiveHandler = n.NextImplementer
		return Response{
			Img:     img,
			Msg:     t,
			Status:  n.W.ActiveHandler.Status(),
			Options: n.W.ActiveHandler.Options(),
		}, ""
	}
}
