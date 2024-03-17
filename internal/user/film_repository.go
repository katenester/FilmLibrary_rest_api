package user

import "time"

type Movie struct {
	MovieID          int       `json:"movieid"`
	MovieName        string    `json:"moviename"`
	MovieDescription string    `json:"moviedescription"`
	MovieReleaseDate time.Time `json:"release_date"`
	MovieRating      uint8     `json:"movierating"`
}
