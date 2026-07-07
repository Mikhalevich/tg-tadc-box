package button

type Operation string

const (
	OperationOpenBox Operation = "OpenBox"
)

func (o Operation) String() string {
	return string(o)
}
