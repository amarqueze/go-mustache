package main

import (
	"errors"
	"fmt"
	"strings"
)

type Person struct {
	Name     string
	LastName string
	Home     Address
}

type Address struct {
	Street string
	City   string
	State  string
	Zip    string
}

func (p Person) ToString() string {
	return p.Name + " " + p.LastName + " lives at " + p.Home.Street + ", " + p.Home.City + ", " + p.Home.State + " " + p.Home.Zip
}

func main() {
	person := Person{Name: "John", LastName: "Doe", Home: Address{Street: "123 Main St", City: "Anytown", State: "CA", Zip: "12345"}}
	person2 := person
	person2.Name = "Jane"
	person2.Home.City = "Othertown"

	fmt.Println("Hello, " + strings.ToUpper(person2.ToString()) + "!")

	if person.Name == person2.Name {
		fmt.Println("Both persons have the same name.")
	} else {
		fmt.Println("The persons have different names.")
	}

	array := []string{"Option 1", "Option 2", "Option 3", "Option 4", "Option 5"}
	array = append(array, "Option 6")
	for _, value := range array {
		fmt.Print(value + " | ")
	}

	mapData := map[string]int{
		"Option 1": 1,
		"Option 2": 2,
	}
	mapData["Option 3"] = 3

	for key, value := range mapData {
		fmt.Printf("\n%s: %d", key, value)
	}

	switchValue := "Option 2"
	switch switchValue {
	case "Option 1":
		fmt.Println("You selected value: " + fmt.Sprint(mapData["Option 1"]))
	case "Option 2":
		fmt.Println("You selected value: " + fmt.Sprint(mapData["Option 2"]))
	default:
		fmt.Println("You selected an unknown option")
	}

	_, exists := mapData["Option 3"]
	fmt.Println("\nOption 3 exists:", exists)

	str, num := getTupla()
	fmt.Printf("\nTuple: %s, %d\n", str, num)

	result, err := isGreaterThanTen(-15)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Is greater than ten:", result)
	}
}

func getTupla() (string, int) {
	return "Hello", 42
}

func isGreaterThanTen(num int) (bool, error) {
	if num < 0 {
		return false, errors.New("number cannot be negative")
	}

	if num > 10 {
		return true, nil
	}

	return false, nil
}
