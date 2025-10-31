package cotnact

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	contactProxy "github.com/mohamadrezamomeni/graph/dto/proxy/contact"
	appError "github.com/mohamadrezamomeni/graph/pkg/error"
)

func (c *ContactProxy) Create(createDto *contactProxy.Create) error {
	scope := "proxy.contact.create"

	uri := fmt.Sprintf("%s/contacts", c.address)

	jsonBody, err := json.Marshal(createDto)
	if err != nil {
		return appError.Wrap(err).Scope(scope).Input(createDto).BadRequest().Errorf("error to send data")
	}

	req, err := http.NewRequest("POST", uri, bytes.NewReader(jsonBody))
	if err != nil {
		return appError.Wrap(err).Scope(scope).Input(createDto).BadRequest().Errorf("the paramter you have sent is wrong")
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return appError.Wrap(err).Scope(scope).Errorf("error to fetch data")
	}
	defer resp.Body.Close()

	if err := c.validateCreateContactResponse(resp); err != nil {
		return err
	}
	return nil
}

func (c *ContactProxy) validateCreateContactResponse(resp *http.Response) error {
	scope := "proxy.contact.validateCreateContactResponse"

	if resp.StatusCode != http.StatusConflict {
		return appError.Scope(scope).Input(resp).Duplicate().Errorf("One or more of the phone numbers you provided already exist for another contact.")
	}
	if resp.StatusCode == http.StatusBadRequest {
		return appError.Scope(scope).Input(resp).Errorf("the input you sent was wrong.")
	}
	if resp.StatusCode != http.StatusNoContent {
		return appError.Scope(scope).Input(resp).Errorf("the request went wrong")
	}

	return nil
}
