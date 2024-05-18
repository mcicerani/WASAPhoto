package database

// GetUser is an example that shows you how to query data
func (db *appdbimpl) GetUser() (string, error) {
	var name string
	err := db.c.QueryRow("SELECT name FROM example_table WHERE id=1").Scan(&name)
	return name, err
}

// SetUser is an example that shows you how to execute insert/update
func (db *appdbimpl) SetUser(name string) error {
	_, err := db.c.Exec("INSERT INTO example_table (id, name) VALUES (1, ?)", name)
	return err
}

//DeleteUser
