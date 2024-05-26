package testutil

import (
	"encoding/json"
	"fmt"
)

func PrettyPrint(v interface{}) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", b)

	return nil
}
