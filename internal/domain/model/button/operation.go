package button

type Operation string

const (
	OperationGetBox                Operation = "GetBox"
	OperationOpenBox               Operation = "OpenBox"
	OperationCardPage              Operation = "CardPage"
	OperationCardTotal             Operation = "CardTotal"
	OperationAbstractDuplicatesAll Operation = "AbstractDuplicatesAll"
	OperationBuyBox                Operation = "BuyBox"
	OperationShop                  Operation = "Shop"
)

func (o Operation) String() string {
	return string(o)
}
