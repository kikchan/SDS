package main

//User password structure
type passwordsData struct {
	Username string
	Password string
	Modified string
	Site     string
}

//User card structure
type cardsData struct {
	Pan    string
	Ccv    string
	Expiry string
	Owner  string
}

//User note structure
type notesData struct {
	Date string
	Text string
}

//User's personal data structure
type userData struct {
	Name       string
	Surname    string
	Email      string
	PrivateKey string
}

//User structure
type user struct {
	Username string
	Password string
	Hash     []byte
	Salt     []byte
	Data     map[string]string
}

//Server's response
type resp struct {
	Ok  bool
	Msg string
}
