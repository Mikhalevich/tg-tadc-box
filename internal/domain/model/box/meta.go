package box

type Meta struct {
	BonusBox BonusBox
}

type BonusBox struct {
	IsValid  bool
	Type     Type
	Attempts int
}
