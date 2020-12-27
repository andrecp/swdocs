package swdocs

func NewSwDoc(name string, title string, description string) *SwDoc {
	swdoc := SwDoc{name: name, title: title, description: description}
	return &swdoc
}
