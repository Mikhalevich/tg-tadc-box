package gloink

import (
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

type Amount int

func (a Amount) Int() int {
	return int(a)
}

func (a Amount) Multiply(count int) Amount {
	return Amount(a.Int() * count)
}

type BoxCosts map[reward.RewardType]Amount
