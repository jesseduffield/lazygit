// code examples for godoc

package toml_test

import (
	"fmt"
	"log"

	toml "github.com/pelletier/go-toml"
)

func Example_tree() {
	config, err := toml.LoadFile("config.toml")

	if err != nil {
		fmt.Println("Error ", err.Error())
	} else {
		// retrieve data directly
		user := config.Get("postgres.user").(string)
		password := config.Get("postgres.password").(string)

		// or using an intermediate object
		configTree := config.Get("postgres").(*toml.Tree)
		user = configTree.Get("user").(string)
		password = configTree.Get("password").(string)
		fmt.Println("User is", user, " and password is", password)

		// show where elements are in the file
		fmt.Printf("User position: %v\n", configTree.GetPosition("user"))
		fmt.Printf("Password position: %v\n", configTree.GetPosition("password"))
	}
}

func Example_unmarshal() {
	type Employer struct {
		Name  string
		Phone string
	}
	type Person struct {
		Name     string
		Age      int64
		Employer Employer
	}

	document := []byte(`
	name = "John"
	age = 30
	[employer]
		name = "Company Inc."
		phone = "+1 234 567 89012"
	`)

	person := Person{}
	toml.Unmarshal(document, &person)
	fmt.Println(person.Name, "is", person.Age, "and works at", person.Employer.Name)
	// Output:
	// John is 30 and works at Company Inc.
}

func ExampleMarshal() {
	type Postgres struct {
		User     string `toml:"user"`
		Password string `toml:"password"`
		Database string `toml:"db" commented:"true" comment:"not used anymore"`
	}
	type Config struct {
		Postgres Postgres `toml:"postgres" comment:"Postgres configuration"`
	}

	config := Config{Postgres{User: "pelletier", Password: "mypassword", Database: "old_database"}}
	b, err := toml.Marshal(config)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
	// Output:
	// # Postgres configuration
	// [postgres]
	//
	//   # not used anymore
	//   # db = "old_database"
	//   password = "mypassword"
	//   user = "pelletier"
}

func ExampleUnmarshal() {
	type Postgres struct {
		User     string
		Password string
	}
	type Config struct {
		Postgres Postgres
	}

	doc := []byte(`
	[postgres]
	user = "pelletier"
	password = "mypassword"`)

	config := Config{}
	toml.Unmarshal(doc, &config)
	fmt.Println("user=", config.Postgres.User)
	// Output:
	// user= pelletier
}
