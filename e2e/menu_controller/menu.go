package menu_controller

import (
	"net/http"
)

func RunMenu(client *http.Client) error {

	err := clientMenu(client)
	if err != nil {
		return err
	}

	return nil
}
