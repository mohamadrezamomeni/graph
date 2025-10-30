package contact

type Filter struct {
	LastNames  string `query:"last_names"`
	FirstNames string `query:"first_names"`
	Phones     string `query:"phones"`
}
