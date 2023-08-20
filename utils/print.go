package utils

import (
	"encoding/json"
	"fmt"
)

func PrettyPrint(m interface{}) error {
	if b, err := SPrettyPrint(m); err != nil {
		return err
	} else {
		fmt.Print(b)
		fmt.Println()
		return nil
	}
}

func SPrettyPrint(m interface{}) (string, error) {
	if b, err := json.MarshalIndent(&m, "", "  "); err != nil {
		return "", err
	} else {
		return string(b), nil
	}
}

func SPrettyPrintFlatten(m interface{}) (string, error) {
	if b, err := json.Marshal(&m); err != nil {
		return "", err
	} else {
		return string(b), nil
	}
}
