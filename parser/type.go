package parser

type CarInfo struct {
	FullName   string
	Milage     string
	Year       string
	Price      float64
	Power      int
	EngineSize int
	Drive      string
	FuelType   string
	SpecID     string
	CarId      string
}


// For parsing car specs
type SpecResponse struct {
	ReturnCode int        `json:"returncode"`
	Message    string     `json:"message"`
	Result     SpecResult `json:"result"`
}

type SpecResult struct {
	SpecID         int             `json:"specid"`
	ParamTypeItems []ParamTypeItem `json:"paramtypeitems"`
}

type ParamTypeItem struct {
	Name       string      `json:"name"`
	ParamItems []ParamItem `json:"paramitems"`
}

type ParamItem struct {
	Name  string `json:"name"`
	ID    int    `json:"id"`
	Value string `json:"value"`
}
