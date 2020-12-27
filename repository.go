package swdocs

import "database/sql"

func CreateSwDoc(db *sql.DB, swdoc *SwDoc) error {
	query := "INSERT INTO swdocs (name, description) VALUES (?, ?)"
	statement, err := db.Prepare(query)
	if err != nil {
		return err
	}

	res, err := statement.Exec(swdoc.Name, swdoc.Description)
	if err != nil {
		return err
	}

	// Update the Id in the structure.
	lid, err := res.LastInsertId()
	swdoc.Id = lid

	return nil
}

func GetMostRecentSwDocs() *[]SwDoc {
	docs := make([]SwDoc, 0)
	return &docs
}
