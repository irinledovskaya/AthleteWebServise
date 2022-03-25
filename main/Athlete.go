package main

import (
	"fmt"
	"time"
)

type Athlete struct {
	Id        int
	Birth     time.Time
	SportClub string
	Name      string
	Surname   string
	Weight    float64
	Attempt   Attempts
}

func (a *Athlete) toString() string {
	return fmt.Sprintf("%v %v %v %v %v %v\n", a.Id, a.Birth.Format("2006-01-02"),
		a.SportClub, a.Name, a.Surname, a.Weight)
}

type AthleteTable struct {
	Caption string
	Table   []Athlete
}

type Attempts struct {
	sn1, sn2, sn3, cj1, cj2, cj3 int
}

type successButton struct {
	Ref string
	Cap string
}
