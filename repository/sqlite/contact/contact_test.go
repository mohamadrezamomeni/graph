package contact

import (
	"os"
	"testing"

	contactRepoDto "github.com/mohamadrezamomeni/graph/dto/repository/contact"
	"github.com/mohamadrezamomeni/graph/repository/migrate"
	"github.com/mohamadrezamomeni/graph/repository/sqlite"
)

var contact *Contact

func TestMain(m *testing.M) {
	config := &sqlite.DBConfig{
		Path: "contact-test.db",
	}

	migrate := migrate.New(config)
	migrate.UP()

	db := sqlite.New(config)

	contact = New(db)

	code := m.Run()

	migrate.DOWN()

	os.Exit(code)
}

func TestCreateContact(t *testing.T) {
	defer contact.deleteAll()

	c1 := &contactRepoDto.Create{
		FirstName: "ali",
		LastName:  "pirzadeh",
		Phones:    []string{"+989123456789", "+989113456789"},
	}

	err := contact.Create(c1)
	if err != nil {
		t.Fatalf("somehting went wrong that was %v", err)
	}

	c2 := &contactRepoDto.Create{
		FirstName: "ali",
		LastName:  "pirzadeh",
		Phones:    []string{"+989123456789", "+989383456789"},
	}

	err = contact.Create(c2)
	if err == nil {
		t.Fatal("we expected an error but we got nothing")
	}
}
