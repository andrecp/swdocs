package swdocs

type SwDoc struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type SwDocLink struct {
	modelName   string
	link        string
	header      string
	description string
}
