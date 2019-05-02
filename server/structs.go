package main

//User structure
type user struct {
	ID       int
	Username string
	Password string
	PubKey   string
	Hash     []byte
	Salt     []byte
	Data     string
}

//Data field structure for notes
type notesData struct {
	Date string
	Text string
}

//Database field structure for passwords, notes and cards
type field struct {
	Data    string
	UserKey string
}

//Server response structure
type resp struct {
	Code int
	Msg  string
}
