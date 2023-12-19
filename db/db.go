package db

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

const MONGO_DB_ENV_NAME = "MONGO_DB_NAME"
