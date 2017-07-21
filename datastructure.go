package main

//library member (person)
type Member struct {
	ID    int64
	Name  string
	Email string
}

//category of book
type BookCategory struct {
	ID   int64
	Name string
}

//book (not actual book, but the unique representation with unique ISBN number).
//library can have many copies of the same book.
type Book struct {
	ID   int64
	Name string
	ISBN string
}

//status of library book copies. can be set as rented to someone or available
type BookStatus struct {
	ID             int64
	BookID         int64
	RentedByMember int64 //if not rented then this value is 0
}
