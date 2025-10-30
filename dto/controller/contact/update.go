package contact

type Update struct {
	IdentifyContact
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Phones    []string `json:"phones"`
}
