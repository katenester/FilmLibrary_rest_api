package user

import "time"

type Actors struct {
	ActorsID          int       `json:"actorid"`
	ActorsName        string    `json:"actorname"`
	ActorsGender      string    `json:"actorgender"`
	ActorsDateOfBirth time.Time `json:"actordateofbirth"`
}
