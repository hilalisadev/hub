package commands

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type Reading struct {
	Data struct {
		Repository struct {
			Object struct {
				Entries []struct {
					Oid  string `json:"oid"`
					Name string `json:"name"`
					Type string `json:"type"`
				} `json:"entries"`
			} `json:"object"`
		} `json:"repository"`
	} `json:"data"`
}

func Display(data []byte) error {

	var reading Reading
	err := json.Unmarshal(data, &reading)
	for k, _ := range reading.Data.Repository.Object.Entries {
		display(reading.Data.Repository.Object.Entries[k])
	}
	return err
}

// display will display the details of the provided value.
func display(v interface{}) {

	// Inspect the concrete type value that is passed in.
	rv := reflect.ValueOf(v)

	// Was the value a pointer value.
	if rv.Kind() == reflect.Ptr {

		// Get the value that the pointer points to.
		rv = rv.Elem()
	}

	// Based on the Kind of value customize the display.
	switch rv.Kind() {

	case reflect.Struct:
		displayStruct(rv)
	}
}

// displayStruct will display details about a struct type.
func displayStruct(rv reflect.Value) {

	// Show each field and the field information.
	for i := 0; i < rv.NumField(); i++ {

		// Get field information for this field.
		fld := rv.Type().Field(i)
		fmt.Printf("Fields: %s ", fld.Name)

		// Display the value of this field.
		fmt.Printf("\tValue: ")
		displayValue(rv.Field(i))

		// Add an extra line feed for the display.
		fmt.Println()
	}
}

// displayValue extracts the native value from the reflect value that is
// passed in and properly displays it.
func displayValue(rv reflect.Value) {

	// Display each value based on its Kind.
	switch rv.Type().Kind() {

	case reflect.String:
		fmt.Printf("%s", rv.String())

	case reflect.Int:
		fmt.Printf("%v", rv.Int())

	case reflect.Float32:
		fmt.Printf("%v", rv.Float())

	case reflect.Bool:
		fmt.Printf("%v", rv.Bool())

	case reflect.Slice:
		for i := 0; i < rv.Len(); i++ {
			displayValue(rv.Index(i))
			fmt.Printf(" ")
		}
	}
}
