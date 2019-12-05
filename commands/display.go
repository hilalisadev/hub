package commands

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/hilalisadev/hub/ui"
	"reflect"
)

func init() {
	addFlags()
}
func addFlags() {
	flag.String("", "0", "... Affiche des informations")
	flag.String("--sav", "1", "Sauvegarde des informations dans la base")
	flag.String("--toto", "3", "Sauvegarde des informations dans la base")
}

var out = ui.Stdout

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

func Display(data []byte, option ...string) error {

	var reading Reading
	err := json.Unmarshal(data, &reading)
	if err != nil {
		return err
	}
	if len(option) > 0 {
		isFlag := make(map[string]string)
		for k, v := range option {
			if IsSet(option[k]) {
				isFlag[v] = flag.Lookup(option[k]).DefValue
			} else if option[k] == "--help" {
				help()
				return nil
			}
		}
		for k := range option {
			if len(isFlag) == len(option) {
				y := flag.Lookup(option[k])
				if y != nil {
					_, _ = fmt.Fprintf(out, "%s %s\r\n", y.Usage, y.Name)
					for k := range reading.Data.Repository.Object.Entries {
						display(reading.Data.Repository.Object.Entries[k], isFlag)
					}
				}
			}

		}
	} else {
		help()
	}
	return nil
}

var (
	setFlags map[string]bool
)

// Determines if the flag was actually set exists
func IsSet(name string) bool {
	if setFlags == nil {
		setFlags = make(map[string]bool)
	}
	if flag.Lookup(name) != nil {
		setFlags[name] = true
	}
	return setFlags[name] == true
}

func help() (int, error) {
	return fmt.Fprintf(out, "%s \r\n", `Usage of Display:
--sav 	"... Sauvegarde des informations dans base"
--help	"... Affiche l'aide"`)
	/*	return fmt.Println(`Usage of Display:
		--sav 	"... Sauvegarde des informations dans base"
		--help	"... Affiche l'aide"`)*/
}

// display will display the details of the provided value.
func display(v interface{}, f map[string]string) {
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
		if f != nil {
			for k, _ := range f {
				switch k {
				case "--sav":
					displayStruct(rv)
				case "--toto":

					fmt.Println("toto")
				default:
				}
			}
		}
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

// displayStruct will display details about a struct type.
func displayStructAll(rv reflect.Value) {

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
