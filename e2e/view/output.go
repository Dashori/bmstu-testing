package view

import (
	models "consoleApp/models"
	"fmt"
	"os"
	"text/tabwriter"
)

func PrintPets(pets models.Pets) {
	fmt.Printf("\nВаши питомцы:\n")
	w := new(tabwriter.Writer)

	// minwidth, tabwidth, padding, padchar, flags
	w.Init(os.Stdout, 9, 8, 0, '\t', 0)

	fmt.Fprintf(w, "\n %s\t%s\t%s\t%s\t%s\t%s\n",
		"№", "Id питомца", "Кличка", "Тип", "Возраст", "Уровень здоровья")

	for i, p := range pets.Pets {
		fmt.Fprintf(w, "\n %d\t%d\t%s\t%s\t%d\t%d\n",
			i+1, p.PetId, p.Name, p.Type, p.Age, p.Health)
	}
	w.Flush()

	fmt.Printf("\nКонец!\n\n")
}

func PrintClientInfo(client models.Client) {
	fmt.Printf("\nВаш логин: %s\nВаш Id: %d\n\n", client.Login, client.ClientId)
}
