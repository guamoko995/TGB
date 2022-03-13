package base

import (
	"strings"
)

// Возвращает ошибку если newPos вложено в obj на любую глубину.
func chekLoop(obj Positioner, newPosition Conteiner) error {
	if newPosObj, ok := newPosition.(Positioner); ok {
		if obj == newPosObj {
			// Если obj является newPosition.
			return ErrStorageLoop{
				obj:         obj,
				newPosition: newPosition.(Positioner),
			}
		}
	}

	box, ok := newPosition.(Positioner)
	if !ok {
		return nil
	}
	nextCont := box.Position()
	if nextCont == nil {
		// Если newPosition нигде не хранится.
		return nil
	}
	// Рекурсивный вызов для менее глубокой позиции.
	if err := chekLoop(obj, nextCont); err != nil {
		esl := err.(ErrStorageLoop)
		return ErrStorageLoop{
			obj:         obj,
			newPosition: newPosition.(Positioner),
			nestedErr:   &esl,
		}
	}
	// Если рекурсивный вызов вернул nil.
	return nil
}

// Размещает obj в newPos.
func Place(obj Positioner, newPos Conteiner) error {
	// Проверки допустимости перемещения.
	err := chekLoop(obj, newPos)
	if relocatbleObj, ok := obj.(Relocater); ok {
		if err := relocatbleObj.Relocate(newPos); err != nil {
			// Объект не допускает перемещения.
			return err
		}
	}
	if newTake, ok := newPos.(Taker); ok {
		if err := newTake.Take(obj); err != nil {
			// Новый контейнер не допускает принятия объекта на хранение.
			return err
		}
	}
	oldPos := obj.Position()
	if oldPos != nil {
		if oldGive, ok := oldPos.(Giver); ok {
			if err := oldGive.Give(obj, newPos); err != nil {
				// Текущий контейнер не допускает передачи объекта.
				return err
			}
		}
	}
	// Проверка вместимости.
	if newLimitedPos, ok := newPos.(Limiter); ok {
		if SizedObj, ok := obj.(Sizer); ok {
			// Объект не вмещается.
			if SizedObj.Size() > newLimitedPos.Capacity() {
				return ErrObjExCapacity{obj, newPos}
			}
			// Недостаточно свободного места.
			if SizedObj.Size() > newLimitedPos.Vacancy() {
				return ErrObjExVacancy{obj, newPos}
			}
		}
	}
	if err != nil {
		return err // Попытка петлевого хранения.
	}
	if oldPos != nil {
		oldPos.Remove(obj) // Изъятие с хранения.
	}
	obj.Reposition(newPos) // Перемещение.
	newPos.Add(obj)        //Принятие на хранение.
	return nil
}

// Ошибка петлевого хранения.
type ErrStorageLoop struct {
	obj         Positioner
	newPosition Positioner
	nestedErr   *ErrStorageLoop
}

// Возвращает цепочку вложенных контейнеров.
func (err ErrStorageLoop) list() []string {
	if err.nestedErr == nil {
		return []string{err.newPosition.Name()}
	} else {
		return append(err.nestedErr.list(), err.newPosition.Name())
	}
}

func (est ErrStorageLoop) Error() string {
	str := "попытка петлевого хранения:"
	list := append(est.list(), est.obj.Name())
	str += " " + strings.Join(list, ">")
	return str
}

// Ошибка объект превышает вместимость.
type ErrObjExCapacity struct {
	obj  Positioner
	cont Conteiner
}

func (err ErrObjExCapacity) Error() string {
	return err.obj.Name() + " не помещается " + err.cont.Name("куда")
}

// Ошибка объект превышает свободное место.
type ErrObjExVacancy struct {
	obj  Positioner
	cont Conteiner
}

func (err ErrObjExVacancy) Error() string {
	return err.cont.Name("где") + " не достаточно свободного места для " + err.obj.Name("Р")
}
