package db

const (
	NAME = "movie-ticket-booking"
	URI  = "mongodb://localhost:27017"
)

type Map map[string]any

type Pagination struct {
	Limit int64
	Page  int64
}

type Store struct {
	User    UserStore
	Cinema  CinemaStore
	Movie   MovieStore
	Hall    HallStore
	Booking BookingStore
}
