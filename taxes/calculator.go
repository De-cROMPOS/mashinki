package taxes

import (
	"fmt"
	"mashinki/parser"
	"strconv"
	"strings"
)

const (
	CNYRate     = 11     // Yuan to rubles
	EURRate     = 100    // Euro to rubles
	BaseUtilFee = 20_000 // recycling base fee
)

type fullCarInfo struct {
	CI           *parser.CarInfo
	customsDuty  float64 // таможенная пошлина
	customsFee   float64 // таможенная сбор
	recyclingFee float64 // Утиль сбор
}

// func that counts all taxes
func NewFullCarInfo(ci *parser.CarInfo) *fullCarInfo {
	fci := &fullCarInfo{
		CI: ci,
	}
	fci.Calculate()
	return fci
}

func (c *fullCarInfo) Calculate() {
	c.calculateCustomsFee()
	c.calculateCustomsDuty()
	c.calculateRecyclingFee()
}

// getting car's age
func (c *fullCarInfo) getCarAge() int {
	parts := strings.Split(c.CI.Year, "-")
	if len(parts) != 2 {
		return 0
	}

	carYear, _ := strconv.Atoi(parts[0])
	currentYear := 2025

	return currentYear - carYear
}

// Утиль сбор
func (c *fullCarInfo) calculateRecyclingFee() {
	age := c.getCarAge()
	isNew := age < 3
	engineSize := c.CI.EngineSize

	var coef float64
	switch {
	case engineSize <= 1000:
		coef = 0.17
		if !isNew {
			coef = 0.26
		}
	case engineSize <= 2000:
		coef = 0.17
		if !isNew {
			coef = 0.26
		}
	case engineSize <= 3000:
		coef = 0.17
		if !isNew {
			coef = 0.26
		}
	case engineSize <= 3500:
		if isNew {
			coef = 107.67
		} else {
			coef = 165.84
		}
	default:
		if isNew {
			coef = 137.11
		} else {
			coef = 180.24
		}
	}

	c.recyclingFee = BaseUtilFee * coef
}

// таможка
func (c *fullCarInfo) calculateCustomsFee() {
	priceRub := c.CI.Price * CNYRate

	switch {
	case priceRub <= 200_000:
		c.customsFee = 1_067
	case priceRub <= 450_000:
		c.customsFee = 2_134
	case priceRub <= 1_200_000:
		c.customsFee = 4_269
	case priceRub <= 2_700_000:
		c.customsFee = 11_746
	case priceRub <= 4_200_000:
		c.customsFee = 16_524
	case priceRub <= 5_500_000:
		c.customsFee = 21_344
	case priceRub <= 7_000_000:
		c.customsFee = 27_540
	default:
		c.customsFee = 30_000
	}
}

// пошлина для < 3 лет
func (c *fullCarInfo) calculateCustomsDutyUnder3Years() float64 {
	priceEUR := c.CI.Price * CNYRate / EURRate
	engineSize := float64(c.CI.EngineSize)

	var rate, minPerCC float64
	switch {
	case priceEUR <= 8500:
		rate = 0.54
		minPerCC = 2.5
	case priceEUR <= 16700:
		rate = 0.48
		minPerCC = 3.5
	case priceEUR <= 42300:
		rate = 0.48
		minPerCC = 5.5
	case priceEUR <= 84500:
		rate = 0.48
		minPerCC = 7.5
	case priceEUR <= 169000:
		rate = 0.48
		minPerCC = 15.0
	default:
		rate = 0.48
		minPerCC = 20.0
	}

	percentDuty := priceEUR * rate * EURRate
	minDuty := engineSize * minPerCC * EURRate

	if minDuty > percentDuty {
		return minDuty
	}
	return percentDuty
}

//  для (3-5) и (5+) лет
func (c *fullCarInfo) calculateCustomsDutyOver3Years() float64 {
	engineSize := float64(c.CI.EngineSize)
	age := c.getCarAge()
	isOver5 := age > 5

	var ratePerCC float64
	switch {
	case engineSize <= 1000:
		ratePerCC = 3.0
		if !isOver5 {
			ratePerCC = 1.5
		}
	case engineSize <= 1500:
		ratePerCC = 3.2
		if !isOver5 {
			ratePerCC = 1.7
		}
	case engineSize <= 1800:
		ratePerCC = 3.5
		if !isOver5 {
			ratePerCC = 2.5
		}
	case engineSize <= 2300:
		ratePerCC = 4.8
		if !isOver5 {
			ratePerCC = 2.7
		}
	case engineSize <= 3000:
		ratePerCC = 5.0
		if !isOver5 {
			ratePerCC = 3.0
		}
	default:
		ratePerCC = 5.7
		if !isOver5 {
			ratePerCC = 3.6
		}
	}

	return engineSize * ratePerCC * EURRate
}

// Пошлина для всех возрастов
func (c *fullCarInfo) calculateCustomsDuty() {
	age := c.getCarAge()

	if age < 3 {
		c.customsDuty = c.calculateCustomsDutyUnder3Years()
	} else {
		c.customsDuty = c.calculateCustomsDutyOver3Years()
	}
}

func (c fullCarInfo) String() string {
	return fmt.Sprintf(
		"🚗 *%s*\n\n"+
			"📅 Год выпуска: %s\n"+
			"📊 Пробег: %v\n"+
			"💰 Цена: %.2f тугриков\n\n"+
			"🔧 Характеристики:\n"+
			"   • Двигатель: %d см³\n"+
			"   • Мощность: %d kW\n"+
			"   • Привод: %s\n"+
			"   • Топливо: %s\n\n"+
			"💳 Таможенные платежи:\n"+
			"   • Пошлина: %.2f ₽\n"+
			"   • Сбор: %.2f ₽\n"+
			"   • Утилизационный сбор: %.2f ₽\n\n"+
			"💵 Итого к оплате: %.2f ₽",
		c.CI.FullName,
		c.CI.Year,
		c.CI.Milage,
		c.CI.Price,
		c.CI.EngineSize,
		c.CI.Power,
		c.CI.Drive,
		c.CI.FuelType,
		c.customsDuty,
		c.customsFee,
		c.recyclingFee,
		c.customsDuty+c.customsFee+c.recyclingFee+c.CI.Price*CNYRate,
	)
}
