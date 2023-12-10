package db

const (
	DBNAME     = "movie-ticket-booking"
	TestDBNAME = "movie-ticket-booking-test"
	DBURI      = "mongodb://localhost:27017"
)

type Store struct {
	User   UserStore
	Cinema CinemaStore
	Movie  MovieStore
	Hall   HallStore
}
