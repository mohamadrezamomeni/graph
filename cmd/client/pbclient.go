package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	contactProxyDto "github.com/mohamadrezamomeni/graph/dto/proxy/contact"
	appLogger "github.com/mohamadrezamomeni/graph/pkg/log"
	contactProxy "github.com/mohamadrezamomeni/graph/proxy/cotnact"
)

var serverAddress string = "http://localhost:1234/api/v1"

func main() {
	appLogger.DiscardLogging()
	if len(os.Args) < 2 || os.Args[1] == "help" || os.Args[1] == "-h" || os.Args[1] == "--help" {
		printGeneralHelp()
		return
	}

	switch os.Args[1] {

	case "create":
		create()
	case "update":
		update()
	case "filter":
		filter()
	default:
		fmt.Println("Unknown subcommand. Use 'help' to see available commands.")
		os.Exit(1)
	}
}

func printGeneralHelp() {
	fmt.Println("Usage: pbclient <subcommand> [flags]")
	fmt.Println("\nSubcommands:")
	fmt.Println("  create    Create a new contact")
	fmt.Println("  update    Update an existing contact")
	fmt.Println("  filter    Filter contacts")
	fmt.Println("\nUse 'pbclient <subcommand> --help' for more information about a subcommand.")
}

func printCreateHelp() {
	fmt.Println("Usage: pbclient create --firstname <first> --lastname <last> --phones <phone1,phone2,...>")
	fmt.Println("\nFlags:")
	fmt.Println("  --firstname   First name of the contact (required)")
	fmt.Println("  --lastname    Last name of the contact (required)")
	fmt.Println("  --phones      Comma-separated list of phones (required)")
}

func printUpdateHelp() {
	fmt.Println("Usage: pbclient update --id <id> [--firstname <first>] [--lastname <last>] [--phones <phone1,phone2,...>]")
	fmt.Println("\nFlags:")
	fmt.Println("  --id          ID of the contact (required)")
	fmt.Println("  --firstname   First name (optional)")
	fmt.Println("  --lastname    Last name (optional)")
	fmt.Println("  --phones      Comma-separated list of phones (optional)")
}

func printFilterHelp() {
	fmt.Println("Usage: pbclient filter [--firstnames <first1,first2,...>] [--lastnames <last1,last2,...>] [--phones <phone1,phone2,...>]")
	fmt.Println("\nFlags:")
	fmt.Println("  --firstnames  Comma-separated first names to filter (optional)")
	fmt.Println("  --lastnames   Comma-separated last names to filter (optional)")
	fmt.Println("  --phones      Comma-separated phones to filter (optional)")
}

func create() {
	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	firstName := createCmd.String("firstname", "", "First name (required)")
	lastName := createCmd.String("lastname", "", "Last name (required)")
	phones := createCmd.String("phones", "", "Comma-separated list of phones (required)")

	if len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help") {
		printCreateHelp()
		return
	}

	createCmd.Parse(os.Args[2:])

	if *firstName == "" || *lastName == "" || *phones == "" {
		fmt.Println("Error: All flags --firstname, --lastname, --phones are required")
		printCreateHelp()
		os.Exit(1)
	}

	phoneList := strings.Split(*phones, ",")

	cp := contactProxy.New(serverAddress)

	err := cp.Create(&contactProxyDto.Create{
		FirstName: *firstName,
		LastName:  *lastName,
		Phones:    phoneList,
	})
	if err != nil {
		fmt.Println("the request you have sent was wrong")
		return
	}

	fmt.Println("Your request has been successful.")
}

func update() {
	updateCmd := flag.NewFlagSet("update", flag.ExitOnError)
	id := updateCmd.String("id", "", "ID of the contact (required)")
	firstName := updateCmd.String("firstname", "", "First name (optional)")
	lastName := updateCmd.String("lastname", "", "Last name (optional)")
	phones := updateCmd.String("phones", "", "Comma-separated list of phones (optional)")

	if len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help") {
		printUpdateHelp()
		return
	}

	updateCmd.Parse(os.Args[2:])

	if *id == "" || *firstName == "" || *lastName == "" || *phones == "" {
		fmt.Println("Error: All flags --id --firstname, --lastname, --phones are required for update")
		printUpdateHelp()
		os.Exit(1)
	}

	phoneList := []string{}
	if *phones != "" {
		phoneList = strings.Split(*phones, ",")
	}

	cp := contactProxy.New(serverAddress)
	err := cp.Update(*id, &contactProxyDto.Update{
		FirstName: *firstName,
		LastName:  *lastName,
		Phones:    phoneList,
	})
	if err != nil {
		fmt.Println("the request you have sent was wrong")
		return
	}

	fmt.Println("Your request has been successful.")
}

func filter() {
	filterCmd := flag.NewFlagSet("filter", flag.ExitOnError)
	firstNames := filterCmd.String("firstnames", "", "Comma-separated first names (optional)")
	lastNames := filterCmd.String("lastnames", "", "Comma-separated last names (optional)")
	phones := filterCmd.String("phones", "", "Comma-separated phones (optional)")

	if len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help") {
		printFilterHelp()
		return
	}

	filterCmd.Parse(os.Args[2:])

	firstNameList := []string{}
	lastNameList := []string{}
	phoneList := []string{}

	if *firstNames != "" {
		firstNameList = strings.Split(*firstNames, ",")
	}
	if *lastNames != "" {
		lastNameList = strings.Split(*lastNames, ",")
	}
	if *phones != "" {
		phoneList = strings.Split(*phones, ",")
	}

	cp := contactProxy.New(serverAddress)
	contacts, err := cp.FilterContacts(&contactProxyDto.Filter{
		FirstNames: firstNameList,
		LastNames:  lastNameList,
		Phones:     phoneList,
	})
	if err != nil {
		fmt.Println("the request you have sent was wrong")
		return
	}

	fmt.Println("ID\tFirst Name\tLast Name\tPhones")
	fmt.Println("--------------------------------------------------")
	for _, c := range contacts {
		fmt.Printf("%s\t%s\t%s\t%s\n", c.ID, c.FirstName, c.LastName, strings.Join(c.Phones, ","))
	}
}
