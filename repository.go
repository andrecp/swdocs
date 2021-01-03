package swdocs

import "database/sql"

const (
	createSwDocSQL    = "INSERT INTO swdocs (name, description, sections) VALUES (?, ?, ?)"
	getSwDocSQL       = "SELECT name, description, sections, updated FROM swdocs WHERE name=?"
	getRecentSwDocSQL = "SELECT name, description, created FROM swdocs ORDER BY ID DESC LIMIT 15"
)

func CreateSwDoc(db *sql.DB, swdoc *SwDoc) error {
	statement, err := db.Prepare(createSwDocSQL)
	if err != nil {
		return err
	}

	res, err := statement.Exec(swdoc.Name, swdoc.Description, swdoc.Sections)
	if err != nil {
		return err
	}

	// Update the Id in the structure.
	lid, err := res.LastInsertId()
	swdoc.Id = lid

	return nil
}

func GetMostRecentSwDocs(db *sql.DB) ([]SwDoc, error) {
	rows, err := db.Query(getRecentSwDocSQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var docs []SwDoc

	for rows.Next() {
		var s SwDoc
		if err := rows.Scan(&s.Name, &s.Description, &s.Created); err != nil {
			return nil, err
		}
		docs = append(docs, s)
	}

	return docs, nil
}

func GetSwDocByName(db *sql.DB, name string) (SwDoc, error) {
	var s SwDoc
	statement, err := db.Prepare(getSwDocSQL)
	if err != nil {
		return s, err
	}

	rows, err := statement.Query(name)
	if err != nil {
		return s, err
	}

	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&s.Name, &s.Description, &s.Sections, &s.Updated); err != nil {
			return s, err
		}
	}

	return s, nil

}
