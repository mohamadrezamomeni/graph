package contact

import (
	"os"
	"testing"

	contactRepoDto "github.com/mohamadrezamomeni/graph/dto/repository/contact"
	"github.com/mohamadrezamomeni/graph/entity"
	appLogger "github.com/mohamadrezamomeni/graph/pkg/log"
	"github.com/mohamadrezamomeni/graph/pkg/utils"
	"github.com/mohamadrezamomeni/graph/repository/migrate"
	"github.com/mohamadrezamomeni/graph/repository/sqlite"
)

var contactRepo *Contact

func TestMain(m *testing.M) {
	config := &sqlite.DBConfig{
		Path: "contact-test.db",
	}

	migrate := migrate.New(config)
	appLogger.DiscardLogging()

	migrate.UP()

	db := sqlite.New(config)

	contactRepo = New(db)

	code := m.Run()

	migrate.DOWN()

	os.Exit(code)
}

func TestUpdatingContact(t *testing.T) {
	defer contactRepo.deleteAll()

	id, err := contactRepo.Create(&contactRepoDto.Create{
		FirstName: "ali",
		LastName:  "pirzadeh",
		Phones:    []string{"09123456789", "09113456789"},
	})
	if err != nil {
		t.Fatalf("error to create data")
	}

	err = contactRepo.Update(id, &contactRepoDto.Update{
		FirstName: "ali",
		LastName:  "pirzadeh",
		Phones:    []string{"09127853850", "09127853851"},
	})
	if err != nil {
		t.Fatalf("error to update data")
	}
}

func TestCreateContact(t *testing.T) {
	defer contactRepo.deleteAll()

	c1 := &contactRepoDto.Create{
		FirstName: "ali",
		LastName:  "pirzadeh",
		Phones:    []string{"09123456789", "09113456789"},
	}

	_, err := contactRepo.Create(c1)
	if err != nil {
		t.Fatalf("somehting went wrong that was %v", err)
	}

	c2 := &contactRepoDto.Create{
		FirstName: "ali",
		LastName:  "pirzadeh",
		Phones:    []string{"09123456789", "09383456789"},
	}

	_, err = contactRepo.Create(c2)
	if err == nil {
		t.Fatal("we expected an error but we got nothing")
	}
}

func TestFilterContactsByPhonesGetingComplete(t *testing.T) {
	defer contactRepo.deleteAll()

	contactRepo.Create(&contactRepoDto.Create{
		FirstName: "ali",
		LastName:  "Pirzadeh",
		Phones:    []string{"09123456789", "09113456789"},
	})

	contactRepo.Create(&contactRepoDto.Create{
		FirstName: "ali",
		LastName:  "Pirzadeh",
		Phones:    []string{"09123455789", "09163456789"},
	})

	contacts, err := contactRepo.Filter(&contactRepoDto.Filter{
		Phones: []string{"09123455789"},
	})
	if err != nil {
		t.Fatalf("something went wrong that was %v", err)
	}

	if len(contacts) != 1 {
		t.Fatalf("error to compare data we expected the lengh of result be 1 but we got %d", len(contacts))
	}

	if len(contacts[0].Phones) != 2 {
		t.Fatalf("we expected we got complete phones but we got the phones we query")
	}
}

func TestFilterContacts(t *testing.T) {
	defer contactRepo.deleteAll()

	contactRepo.Create(&contactRepoDto.Create{
		FirstName: "ali",
		LastName:  "Pirzadeh",
		Phones:    []string{"09123456789", "09113456789"},
	})
	contactRepo.Create(&contactRepoDto.Create{
		FirstName: "ali",
		LastName:  "Pirzadeh",
		Phones:    []string{"09123455789", "09163456789"},
	})
	contactRepo.Create(&contactRepoDto.Create{
		FirstName: "yasin",
		LastName:  "ahmadi",
		Phones:    []string{"09173455781", "09663456734"},
	})

	contactRepo.Create(&contactRepoDto.Create{
		FirstName: "babak",
		LastName:  "ahmadi",
		Phones:    []string{"09173455782", "09663456783"},
	})

	contactRepo.Create(&contactRepoDto.Create{
		FirstName: "babak",
		LastName:  "alvandi",
		Phones:    []string{"09173455789", "09663456785"},
	})

	for i, testCase := range []struct {
		input contactRepoDto.Filter
		count int
	}{
		{
			input: contactRepoDto.Filter{
				FirstNames: []string{"ali", "yasin"},
			},
			count: 3,
		},
		{
			input: contactRepoDto.Filter{
				LastNames: []string{"alvandi", "Pirzadeh"},
			},
			count: 3,
		},
		{
			input: contactRepoDto.Filter{
				FirstNames: []string{"babak"},
				LastNames:  []string{"alvandi", "Pirzadeh"},
			},
			count: 1,
		},
		{
			input: contactRepoDto.Filter{
				FirstNames: []string{"babak"},
				LastNames:  []string{"alvandi", "Pirzadeh"},
				Phones:     []string{"09121111111"},
			},
			count: 0,
		},
		{
			input: contactRepoDto.Filter{
				FirstNames: []string{"babak"},
				LastNames:  []string{"alvandi", "Pirzadeh"},
				Phones:     []string{"09173455789"},
			},
			count: 1,
		},
	} {
		contacts, err := contactRepo.Filter(&testCase.input)
		if err != nil {
			t.Errorf("something went wrong at index %d the problem was %v", i, err)
		} else if len(contacts) != testCase.count && isFilterContactsResponseValid(contacts, &testCase.input) {
			t.Errorf("error to compare data at %d", i)
		}
	}
}

func isFilterContactsResponseValid(
	contacts []*entity.Contact,
	filterDto *contactRepoDto.Filter,
) bool {
	for _, contact := range contacts {
		isCondistionSatisFiy := true
		if filterDto.FirstNames != nil && len(filterDto.FirstNames) > 0 {
			isCondistionSatisFiy = utils.IsContain(contact.FirstName, filterDto.FirstNames)
		}
		if filterDto.LastNames != nil && len(filterDto.LastNames) > 0 {
			isCondistionSatisFiy = utils.IsContain(contact.LastName, filterDto.LastNames)
		}
		if filterDto.Phones != nil && len(filterDto.Phones) > 0 {
			isCondistionSatisFiy = hasCommonPhone(filterDto.Phones, contact.Phones)
		}
		if isCondistionSatisFiy != false {
			return false
		}
	}
	return true
}

func hasCommonPhone(filterPhones []string, contactPhones []string) bool {
	for _, phone := range contactPhones {
		if dd := utils.IsContain(phone, filterPhones); dd {
			return true
		}
	}
	return false
}
