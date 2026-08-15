package button

type Operation string

const (
	OperationNoOperation           Operation = "NoOperation"
	OperationOpenBox               Operation = "OpenBox"
	OperationCardPage              Operation = "CardPage"
	OperationCardTotal             Operation = "CardTotal"
	OperationAbstractDuplicatesAll Operation = "AbstractDuplicatesAll"
	OperationAbstractCard          Operation = "AbstractCard"
	OperationBuyBox                Operation = "BuyBox"
	OperationShop                  Operation = "Shop"
	OperationGetCommonBox          Operation = "GetCommonBox"
	OperationLike                  Operation = "Like"
)

func (o Operation) String() string {
	return string(o)
}

func (o Operation) IsNoOperation() bool {
	return o == OperationNoOperation
}
