package button

type Operation string

const (
	OperationOpenBox               Operation = "OpenBox"
	OperationCardPage              Operation = "CardPage"
	OperationCardTotal             Operation = "CardTotal"
	OperationAbstractDuplicatesAll Operation = "AbstractDuplicatesAll"
	OperationBuyBox                Operation = "BuyBox"
	OperationShop                  Operation = "Shop"
	OperationGetCommonBox          Operation = "GetCommonBox"
	OperationLike                  Operation = "Like"
)

func (o Operation) String() string {
	return string(o)
}
