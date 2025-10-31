package contact

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/mohamadrezamomeni/graph/entity"
	"github.com/mohamadrezamomeni/graph/repository/sqlite"

	contactRepoDto "github.com/mohamadrezamomeni/graph/dto/repository/contact"
	appErr "github.com/mohamadrezamomeni/graph/pkg/error"
)

func (c *Contact) Update(id string, updateDto *contactRepoDto.Update) error {
	scope := "repository.contact.update"

	tx, err := c.db.Conn().Begin()
	if err != nil {
		return err
	}

	err = c.deleteContactPhones(tx, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = c.assignPhones(tx, id, updateDto.Phones)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = c.updateContact(tx, id, updateDto)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return appErr.Wrap(err).Scope(scope).Errorf("error to commit")
	}
	return nil
}

func (c *Contact) updateContact(tx *sql.Tx, id string, updateDto *contactRepoDto.Update) error {
	scope := "repository.contact.updateContact"

	_, err := tx.Exec(
		"UPDATE contacts SET first_name = ?, last_name = ? WHERE id = ?",
		updateDto.FirstName, updateDto.LastName, id,
	)
	if err != nil {
		return appErr.Wrap(err).Scope(scope).Input(id, updateDto).ErrorWrite()
	}
	return nil
}

func (c *Contact) deleteContactPhones(tx *sql.Tx, id string) error {
	scope := "repository.contact.deleteExistedPhones"
	_, err := tx.Exec("DELETE FROM Phones WHERE contact_id = $1", id)
	if err != nil {
		return appErr.Wrap(err).Scope(scope).Input(id).ErrorWrite()
	}
	return nil
}

func (c *Contact) Create(createDto *contactRepoDto.Create) (string, error) {
	scope := "repository.contact.create"

	tx, err := c.db.Conn().Begin()
	if err != nil {
		return "", appErr.Wrap(err).Scope(scope).ErrorWrite()
	}

	contactID, err := c.createContact(tx, createDto)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	err = c.assignPhones(tx, contactID, createDto.Phones)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	if err != nil {
		tx.Rollback()
		return "", err
	}
	if err := tx.Commit(); err != nil {
		return "", appErr.Wrap(err).Scope(scope).Errorf("error to commit")
	}
	return contactID, nil
}

func (c *Contact) createContact(tx *sql.Tx, createDto *contactRepoDto.Create) (string, error) {
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

	if err != nil && sqlite.IsDuplicateError(err) {
		return appErr.Wrap(err).Scope(scope).Input(contactID, phones).Duplicate().Errorf("the phone numbers you have sent contain duplicate phone numbers")
	}
	if err != nil {
		return appErr.Wrap(err).Input(contactID, phones).UnExpected().Scope(scope).DebuggingError()
	}

	return nil
}

func (c *Contact) Filter(filterDto *contactRepoDto.Filter) ([]*entity.Contact, error) {
	scope := "repository.contact.FilterContacts"

	query := c.makeFilterContactsQuery(filterDto)
	rows, err := c.db.Conn().Query(query)
	if err != nil {
		return nil, appErr.Wrap(err).Scope(scope).Input(filterDto).Errorf("error to query")
	}
	defer rows.Close()

	return c.scanContacts(rows)
}

func (c *Contact) makeFilterContactsQuery(filterDto *contactRepoDto.Filter) string {
	query := "SELECT * FROM contacts"

	subQueries := make([]string, 0)

	if filterDto.FirstNames != nil && len(filterDto.FirstNames) > 0 {
		subQueries = append(subQueries,
			fmt.Sprintf("first_name IN ('%s')", strings.Join(filterDto.FirstNames, "', '")),
		)
	}

	if filterDto.LastNames != nil && len(filterDto.LastNames) > 0 {
		subQueries = append(subQueries,
			fmt.Sprintf("last_name IN ('%s')", strings.Join(filterDto.LastNames, "', '")),
		)
	}

	if len(subQueries) > 0 {
		query = fmt.Sprintf("%s WHERE %s", query, strings.Join(subQueries, " AND "))
	}

	queryIncludedPhones := fmt.Sprintf(
		"SELECT c.id, c.first_name, c.last_name, GROUP_CONCAT(p.phone, ',') AS phones "+
			"FROM (%s) AS c "+
			"LEFT JOIN phones AS p ON c.id = p.contact_id",
		query,
	)

	if filterDto.Phones != nil && len(filterDto.Phones) > 0 {
		queryIncludedPhones = fmt.Sprintf(
			"SELECT c.id, c.first_name, c.last_name, GROUP_CONCAT(p2.phone, ',') AS phones "+
				"FROM (%s) AS c "+
				"INNER JOIN phones AS p ON c.id = p.contact_id  AND p.phone IN ('%s')"+
				"LEFT JOIN phones AS p2 ON p2.contact_id = c.id",
			query,
			strings.Join(filterDto.Phones, "', '"),
		)
	}

	queryIncludedPhones += " GROUP BY c.id, c.first_name, c.last_name"

	return queryIncludedPhones
}

func (c *Contact) scanContacts(rows *sql.Rows) ([]*entity.Contact, error) {
	scope := "repository.contact.scanContacts"

	contacts := make([]*entity.Contact, 0)
	for rows.Next() {
		var phones sql.NullString
		contact := new(entity.Contact)

		err := rows.Scan(&contact.ID, &contact.FirstName, &contact.LastName, &phones)
		if err != nil {
			return nil, appErr.Wrap(err).Scope(scope).Errorf("error to scan")
		}

		if phones.Valid {
			contact.Phones = strings.Split(phones.String, ",")
		} else {
			contact.Phones = []string{}
		}

		contacts = append(contacts, contact)
	}

	return contacts, nil
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
