package relateiq

import (
	"fmt"
	"os"
)

func Example() {
	key, secret := os.Getenv("RELATEIQ_API_KEY"), os.Getenv("RELATEIQ_API_SECRET")
	if key == "" || secret == "" {
		fmt.Println("Example requires RelateIQ API credentials (which you can obtain from your organization's integration settings screen) to be set in the RELATEIQ_API_KEY and RELATEIQ_API_SECRET environment variables.")
		return
	}

	c := NewClient(nil, Credentials{APIKey: key, APISecret: secret})

	accounts, _, err := c.Accounts.List(AccountsListOptions{})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, a := range accounts {
		fmt.Printf("%-25s  %-24s  %s\n", a.Name, a.ID, a.ModifiedDate)
	}

	// output is dependent on your RelateIQ data
}
