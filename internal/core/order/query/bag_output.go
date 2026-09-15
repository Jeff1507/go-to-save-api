package query

type BagItemOutput struct {
	ID       string
	Name     string
	ImageURL *string
	Quantity int32
}

type BagSummaryOutput struct {
	ID        string
	Title     string
	Quantity  int32
	UnitPrice int64
	Items     []BagItemOutput
}
