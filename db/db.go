package db

const (
	NAME = "movie-ticket-booking"
	URI  = "mongodb://localhost:27017"
)

type Store struct {
	User   UserStore
	Cinema CinemaStore
	Movie  MovieStore
	Hall   HallStore
}
