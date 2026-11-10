package main

import (
	"errors"
	"fmt"
	"gomustache/internal/client"
	"gomustache/internal/handler"
	"gomustache/internal/router"
	"gomustache/internal/service"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type shoppingCart struct {
	Items []string
}

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
	defer finally()

	person := Person{Name: "John", LastName: "Doe", Home: Address{Street: "123 Main St", City: "Anytown", State: "CA", Zip: "12345"}}
	person2 := person
	person2.Name = "Jane"
	person2.Home.City = "Othertown"

	fmt.Println("Hello, " + strings.ToUpper(person2.ToString()) + "!")

	multiplicar := func(a int, b int) int {
		return a * b
	}
	fmt.Println(multiplicar(3, 4))

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
		fmt.Printf("%s: %d\n", key, value)
	}

	switchValue := "Option 2"
	switch switchValue {
	case "Option 1":
		fmt.Println("You selected value: " + strconv.Itoa(mapData["Option 1"]))
	case "Option 2":
		fmt.Println("You selected value: " + strconv.Itoa(mapData["Option 2"]))
	default:
		fmt.Println("You selected an unknown option")
	}

	_, exists := mapData["Option 3"]
	fmt.Println("Option 3 exists:", exists)

	str, num := getTupla()
	fmt.Printf("Tuple: %s, %d\n", str, num)

	result, err := isGreaterThanTen(-15)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Is greater than ten:", result)
	}

	printPersonInfo(person, person2)

	shoppingCart := &shoppingCart{}
	addItem := AddItem(shoppingCart)
	fmt.Println("Items in cart:", addItem("Apple"))
	fmt.Println("Items in cart:", addItem("Banana"))
	fmt.Println("Items in cart:", addItem("Orange"))
	fmt.Println("Total: ", len(shoppingCart.Items))

	addAddress := addAdress()
	addAddress(&Address{Street: "456 Elm St", City: "Sometown", State: "TX", Zip: "67890"})
	addAddress(&Address{Street: "789 Oak St", City: "Newtown", State: "FL", Zip: "54321"})
	addressSoFar := addAddress(&Address{Street: "101 Pine St", City: "Oldtown", State: "WA", Zip: "98765"})
	fmt.Println("Addresses so far:", addressSoFar)

	//canal := make(chan string)
	//go backgroundTask(canal)
	//processMessages(canal)

	runServer()
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

func printPersonInfo(p ...Person) {
	for _, person := range p {
		fmt.Println(person.ToString())
	}
}

func AddItem(cart *shoppingCart) func(item string) int {
	return func(item string) int {
		cart.Items = append(cart.Items, item)
		return len(cart.Items)
	}
}

func addAdress() func(address *Address) *[]Address {
	AllAddresses := []Address{}
	return func(address *Address) *[]Address {
		AllAddresses = append(AllAddresses, *address)
		return &AllAddresses
	}
}

func backgroundTask(canal chan<- string) {
	fmt.Println("Background task is running...")
	for i := 0; i < 5; i++ {
		time.Sleep(time.Duration(i) * time.Second)
		canal <- fmt.Sprintf("Message %d from background task", i)
	}
	close(canal)
}

func processMessages(canal <-chan string) {
	for msg := range canal {
		fmt.Println("Received:", msg)
	}
}

func finally() {
	fmt.Println("Finally block executed.")
}

func runServer() {
	httpClient := &http.Client{Timeout: 5 * time.Second}
	jpClient := client.NewJSONPlaceholderClient("https://jsonplaceholder.typicode.com", httpClient)

	userService := service.NewUserService(jpClient, 2*time.Second)

	r := router.New(router.Handlers{
		User: handler.NewUserHandler(userService),
	})

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
