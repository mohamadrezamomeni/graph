package cotnact

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	contactProxy "github.com/mohamadrezamomeni/graph/dto/proxy/contact"
	appError "github.com/mohamadrezamomeni/graph/pkg/error"
)

func (c *ContactProxy) Update(id string, updateDto *contactProxy.Update) error {
	scope := "proxy.contact.update"

	uri := fmt.Sprintf("%s/contacts/%s", c.address, id)

	jsonBody, err := json.Marshal(updateDto)
	if err != nil {
		return appError.Wrap(err).Scope(scope).Input(id, updateDto).BadRequest().Errorf("error to send data")
	}

	req, err := http.NewRequest("PUT", uri, bytes.NewReader(jsonBody))
	if err != nil {
		return appError.Wrap(err).Scope(scope).Input(id, updateDto).BadRequest().Errorf("the paramter you have sent is wrong")
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return appError.Wrap(err).Scope(scope).Errorf("error to fetch data")
	}
	defer resp.Body.Close()

	if err := c.validateUpdateContact(resp); err != nil {
		return err
	}

	return nil
}

func (c *ContactProxy) validateUpdateContact(resp *http.Response) error {
	scope := "proxy.contact.validateUpdateContact"

	if resp.StatusCode != http.StatusConflict {
		return appError.Scope(scope).Input(resp).Duplicate().Errorf("One or more of the phone numbers you provided already exist for another contact")
	}
	if resp.StatusCode == http.StatusBadRequest {
		return appError.Scope(scope).Input(resp).Errorf("the input you sent was wrong.")
	}
	if resp.StatusCode != http.StatusNoContent {
		return appError.Scope(scope).Input(resp).Errorf("the request went wrong")
	}
	return nil
}
