package cotnact

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	contactProxy "github.com/mohamadrezamomeni/graph/dto/proxy/contact"
	"github.com/mohamadrezamomeni/graph/entity"
	appError "github.com/mohamadrezamomeni/graph/pkg/error"
)

type ItemContactResponse struct {
	ID        string   `json:"id"`
	LastName  string   `json:"last_name"`
	FirstName string   `json:"first_name"`
	Phones    []string `json:"phones"`
}

type FilterContactsResponse struct {
	Items []*ItemContactResponse `json:"items"`
}

func (c *ContactProxy) FilterContacts(filterDto *contactProxy.Filter) ([]*entity.Contact, error) {
	scope := "proxy.contact.filter"

	uri, err := c.makeUriContacts(filterDto)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, appError.Wrap(err).Scope(scope).Input(filterDto).BadRequest().Errorf("the paramter you have sent is wrong")
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, appError.Wrap(err).Scope(scope).Errorf("error to fetch data")
	}
	defer resp.Body.Close()

	err = c.validateFilterResponse(resp)
	if err != nil {
		return nil, err
	}

	return c.getFilterContactsResponse(resp)
}

func (c *ContactProxy) makeUriContacts(filterDto *contactProxy.Filter) (string, error) {
	scope := "proxy.contact.filter"

	uri := fmt.Sprintf("%s/contacts", c.address)

	u, err := url.Parse(uri)
	if err != nil {
		return "", appError.Wrap(err).Scope(scope).Input(filterDto, uri).Errorf("error to parse uri")
	}

	q := u.Query()

	if filterDto.FirstNames != nil && len(filterDto.FirstNames) > 0 {
		q.Set("first_names", strings.Join(filterDto.FirstNames, ","))
	}

	if filterDto.Phones != nil && len(filterDto.Phones) > 0 {
		q.Set("phones", strings.Join(filterDto.Phones, ","))
	}

	if filterDto.LastNames != nil && len(filterDto.LastNames) > 0 {
		q.Set("last_names", strings.Join(filterDto.LastNames, ","))
	}

	u.RawQuery = q.Encode()
	return u.String(), nil
}

func (c *ContactProxy) validateFilterResponse(resp *http.Response) error {
	scope := "proxy.contact.validateFilterResponse"

	if resp.StatusCode == http.StatusBadRequest {
		return appError.Scope(scope).Input(resp).Errorf("the input you sent was wrong.")
	}
	if resp.StatusCode != http.StatusOK {
		return appError.Scope(scope).Input(resp).Errorf("the request went wrong")
	}
	return nil
}

func (c *ContactProxy) getFilterContactsResponse(resp *http.Response) ([]*entity.Contact, error) {
	scope := "proxy.contact.filter"

	var contactsRes FilterContactsResponse
	if err := json.NewDecoder(resp.Body).Decode(&contactsRes); err != nil {
		return nil, appError.Wrap(err).Scope(scope).Errorf("error to parse data")
	}

	contacts := make([]*entity.Contact, 0)
	for _, contactResponse := range contactsRes.Items {
		contacts = append(contacts, &entity.Contact{
			ID:        contactResponse.ID,
			FirstName: contactResponse.FirstName,
			LastName:  contactResponse.LastName,
			Phones:    contactResponse.Phones,
		})
	}

	return contacts, nil
}
