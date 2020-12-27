package swdocs

type SwDoc struct {
	id          int
	name        string
	description string
}

type SwDocLink struct {
	modelName   string
	link        string
	header      string
	description string
}
