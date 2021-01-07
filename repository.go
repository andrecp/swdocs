package swdocs

import "database/sql"

const (
	dbSchema = `
    CREATE TABLE IF NOT EXISTS swdocs (
		id INTEGER PRIMARY KEY,
		name TEXT UNIQUE,
		user TEXT,
		created NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated NOT NULL DEFAULT CURRENT_TIMESTAMP,
		description TEXT,
		sections TEXT)
	`
	createOrUpdateSwDocSQL = `INSERT INTO swdocs (name, description, user, sections) VALUES (?, ?, ?, ?)
								ON CONFLICT (name) DO UPDATE SET
									sections=excluded.sections,
									description=excluded.description,
									user=excluded.user,
									updated=CURRENT_TIMESTAMP`
	getSwDocSQL              = "SELECT name, description, sections, user, updated FROM swdocs WHERE name=?"
	getRecentCreatedSwDocSQL = "SELECT name, description, user, created, updated FROM swdocs ORDER BY ID DESC LIMIT 15"
	getRecentUpdatedSwDocSQL = "SELECT name, description, user, created, updated FROM swdocs ORDER BY updated DESC LIMIT 15"
	searchSwDocSQL           = "SELECT name, user, updated FROM swdocs WHERE name like ?"
	deleteSwDocSQL           = "DELETE FROM swdocs WHERE name=?"
)

func CreateOrUpdateSwDoc(db *sql.DB, swdoc *SwDoc) error {
	statement, err := db.Prepare(createOrUpdateSwDocSQL)
	if err != nil {
		return err
	}

	res, err := statement.Exec(swdoc.Name, swdoc.Description, swdoc.User, swdoc.Sections)
	if err != nil {
		return err
	}

	// Update the Id in the structure.
	lid, err := res.LastInsertId()
	swdoc.Id = lid

	return nil
}

func GetMostRecentCreatedSwDocs(db *sql.DB) ([]SwDoc, error) {
	rows, err := db.Query(getRecentCreatedSwDocSQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var docs []SwDoc

	for rows.Next() {
		var s SwDoc
		if err := rows.Scan(&s.Name, &s.Description, &s.User, &s.Created, &s.Updated); err != nil {
			return nil, err
		}
		docs = append(docs, s)
	}

	return docs, nil
}

func GetMostRecentUpdatedSwDocs(db *sql.DB) ([]SwDoc, error) {
	rows, err := db.Query(getRecentUpdatedSwDocSQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var docs []SwDoc

	for rows.Next() {
		var s SwDoc
		if err := rows.Scan(&s.Name, &s.Description, &s.User, &s.Created, &s.Updated); err != nil {
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
		if err := rows.Scan(&s.Name, &s.Description, &s.Sections, &s.User, &s.Updated); err != nil {
			return s, err
		}
	}

	return s, nil

}

func SearchSwDocsByName(db *sql.DB, name string) ([]SwDoc, error) {
	var docs []SwDoc

	statement, err := db.Prepare(searchSwDocSQL)
	if err != nil {
		return docs, err
	}

	rows, err := statement.Query(name)
	if err != nil {
		return docs, err
	}

	defer rows.Close()

	for rows.Next() {
		var s SwDoc
		if err := rows.Scan(&s.Name, &s.User, &s.Updated); err != nil {
			return nil, err
		}
		docs = append(docs, s)
	}

	return docs, nil

}

func DeleteSwDoc(db *sql.DB, name string) error {
	statement, err := db.Prepare(deleteSwDocSQL)
	if err != nil {
		return err
	}

	_, err = statement.Exec(name)
	if err != nil {
		return err
	}

	return nil
}
