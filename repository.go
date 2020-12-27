package swdocs

func NewSwDoc(name string, description string) *SwDoc {
	swdoc := SwDoc{name: name, description: description}
	return &swdoc
}
