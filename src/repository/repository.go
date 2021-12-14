package repository

type Repository interface {
	Transaction
	Question
	Admin
	Option
}
