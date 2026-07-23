package button

type Operation string

const (
	OperationOpenBox  Operation = "OpenBox"
	OperationCardPage Operation = "CardPage"
)

func (o Operation) String() string {
	return string(o)
}
