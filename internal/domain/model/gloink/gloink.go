package gloink

import "github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"

type Amount int

func (a Amount) Int() int {
	return int(a)
}

func (a Amount) IsFree() bool {
	return a == 0
}

func AmountFromInt(amount int) Amount {
	return Amount(amount)
}

func (a Amount) Multiply(count int) Amount {
	return Amount(a.Int() * count)
}

type BoxCost struct {
	Type   box.Type
	Amount Amount
}
