package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"
)

type myTime struct {
	time.Time
}

type car struct {
	ID          uint `json:"id"`
	PricePerDay uint `json:"price_per_day"`
	PricePerKm  uint `json:"price_per_km"`
}

type rental struct {
	ID        uint   `json:"id"`
	CarID     uint   `json:"car_id"`
	Startdate myTime `json:"start_date"`
	EndDate   myTime `json:"end_date"`
	Distance  uint   `json:"distance"`
}

type option struct {
	ID       uint   `json:"id"`
	RentalID uint   `json:"rental_id"`
	Type     string `json:"type"`
}

type input struct {
	Cars    []car    `json:"cars"`
	Rentals []rental `json:"rentals"`
	Options []option `json:"options"`
}

type rentalsOut struct {
	ID      uint     `json:"id"`
	Options []string `json:"options"`
	Actions []action `json:"actions"`
}

type action struct {
	Who    string `json:"who"`
	Type   string `json:"type"`
	Amount uint   `json:"amount"`
}

type output struct {
	Rentals []rentalsOut `json:"rentals"`
}

func main() {
	f, err := os.Open("./data/input.json")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var i input
	err = json.NewDecoder(f).Decode(&i)
	if err != nil {
		log.Fatal(err)
	}

	var o output
	for _, rental := range i.Rentals {
		car := i.getCarByID(rental.CarID)
		options := i.getOptionByRental(rental.ID)

		var price uint
		days := (1 + uint(rental.EndDate.Time.Sub(rental.Startdate.Time)/(time.Hour*24)))

		price = car.PricePerDay
		if days > 10 {
			price += uint(float32((days-10)*car.PricePerDay)*0.5) + uint(float32(6*car.PricePerDay)*0.7) + uint(float32(3*car.PricePerDay)*0.9)
		} else if days > 4 {
			price += uint(float32((days-4)*car.PricePerDay)*0.7) + uint(float32(3*car.PricePerDay)*0.9)
		} else if days > 1 {
			price += uint(float32((days-1)*car.PricePerDay) * 0.9)
		}

		price += rental.Distance * car.PricePerKm
		commissionFee := uint(float32(price) * 0.3)

		var optionsFee, additionalInsuranceFee uint
		for _, option := range options {
			switch option {
			case "gps":
				optionsFee += 500 * days
			case "baby_seat":
				optionsFee += 200 * days
			case "additional_insurance":
				additionalInsuranceFee += 1000 * days
			}
		}

		actions := make([]action, 0)

		actions = append(actions, action{
			Who:    "driver",
			Type:   "debit",
			Amount: price + optionsFee + additionalInsuranceFee,
		})

		actions = append(actions, action{
			Who:    "owner",
			Type:   "credit",
			Amount: price - commissionFee + optionsFee,
		})

		insuranceFee := commissionFee / 2
		assistanceFee := days * 100
		drivyFee := commissionFee - insuranceFee - assistanceFee

		actions = append(actions, action{
			Who:    "insurance",
			Type:   "credit",
			Amount: insuranceFee,
		})

		actions = append(actions, action{
			Who:    "assistance",
			Type:   "credit",
			Amount: assistanceFee,
		})

		actions = append(actions, action{
			Who:    "drivy",
			Type:   "credit",
			Amount: drivyFee + additionalInsuranceFee,
		})

		o.Rentals = append(o.Rentals, rentalsOut{
			ID:      rental.ID,
			Options: options,
			Actions: actions,
		})
	}

	fo, err := os.Create("./data/output.json")
	if err != nil {
		log.Fatal(err)
	}

	enc := json.NewEncoder(fo)
	enc.SetIndent("", "  ")
	err = enc.Encode(o)
	if err != nil {
		log.Fatal(err)
	}
}

// We can set cars in a map for get in a better way
func (i *input) getCarByID(id uint) car {
	for _, car := range i.Cars {
		if car.ID == id {
			return car
		}
	}

	return car{}
}

func (i *input) getOptionByRental(rentalid uint) []string {
	out := make([]string, 0)

	for _, option := range i.Options {
		if option.RentalID == rentalid {
			out = append(out, option.Type)
		}
	}
	return out
}

func (t *myTime) UnmarshalJSON(buf []byte) error {
	tt, err := time.Parse("2006-01-2", strings.Trim(string(buf), `"`))
	if err != nil {
		return err
	}
	t.Time = tt
	return nil
}
