package parser

import "fmt"

type carInfo struct {
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

func (c carInfo) String() string {
	return fmt.Sprintf(
		"Автомобиль: %s\n"+
			"Год: %s | Пробег: %v\n"+
			"Цена: %v\n"+
			"Двигатель: %d см³ (%d kW.)\n"+
			"Привод: %s | Топливо: %s",
		c.FullName,
		c.Year,
		c.Milage,
		c.Price,
		c.EngineSize,
		c.Power,
		c.Drive,
		c.FuelType,
	)
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
