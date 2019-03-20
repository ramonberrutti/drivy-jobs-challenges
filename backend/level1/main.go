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

type input struct {
	Cars    []car    `json:"cars"`
	Rentals []rental `json:"rentals"`
}

type rentalsOut struct {
	ID    uint `json:"id"`
	Price uint `json:"price"`
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

		var price uint
		price += rental.Distance * car.PricePerKm
		price += (1 + uint(rental.EndDate.Time.Sub(rental.Startdate.Time)/(time.Hour*24))) * car.PricePerDay

		o.Rentals = append(o.Rentals, rentalsOut{
			ID:    rental.ID,
			Price: price,
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

func (t *myTime) UnmarshalJSON(buf []byte) error {
	tt, err := time.Parse("2006-01-2", strings.Trim(string(buf), `"`))
	if err != nil {
		return err
	}
	t.Time = tt
	return nil
}
