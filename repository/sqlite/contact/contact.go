package contact

import (
	"database/sql"
	"fmt"

	contactRepoDto "github.com/mohamadrezamomeni/graph/dto/repository/contact"
	appErr "github.com/mohamadrezamomeni/graph/pkg/error"
)

func (c *Contact) Create(createDto *contactRepoDto.CreateContact) error {
	scope := "repository.contact.create"

	tx, err := c.db.Conn().Begin()
	if err != nil {
		return appErr.Wrap(err).Scope(scope).ErrorWrite()
	}

	contactID, err := c.createContact(tx, createDto)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = c.assignPhones(tx, contactID, createDto.Phones)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return appErr.Wrap(err).Scope(scope).Errorf("error to commit")
	}
	return nil
}

func (c *Contact) createContact(tx *sql.Tx, createDto *contactRepoDto.CreateContact) (string, error) {
	scope := "repository.repository.createContact"

	var id string
	err := tx.QueryRow(`
	             INSERT INTO contacts (first_name, last_name)
	             VALUES ($1, $2)
	             RETURNING id
	`,
		createDto.FirstName,
		createDto.LastName,
	).Scan(&id)
	if err != nil {
		return "", appErr.Wrap(err).Scope(scope).Input(createDto).DebuggingError()
	}
	return id, nil
}

func (c *Contact) assignPhones(tx *sql.Tx, contactID string, phones []string) error {
	scope := "repository.contact.assignPhones"

	if len(phones) == 0 {
		return nil
	}
	query := `INSERT INTO phones (phone, contact_id) VALUES `
	args := []interface{}{}

	for i, phone := range phones {
		if i > 0 {
			query += ", "
		}
		rowIndex := i * 2
		query += fmt.Sprintf("($%d, $%d)", rowIndex+1, rowIndex+2)
		args = append(args, phone, contactID)
	}

	_, err := tx.Exec(query, args...)
	if err != nil {
		return appErr.Wrap(err).Input(contactID, phones).UnExpected().Scope(scope).DebuggingError()
	}

	return nil
}

func (c *Contact) deleteAll() error {
	scope := "repository.contacts.deleteAll"

	sql := "DELETE FROM contacts"
	res, err := c.db.Conn().Exec(sql)
	if err != nil {
		return appErr.Wrap(err).Scope(scope).ErrorWrite()
	}

	_, err = res.RowsAffected()
	if err != nil {
		return appErr.Wrap(err).Scope(scope).ErrorWrite()
	}

	return nil
}
