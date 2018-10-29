package main

import "fmt"

type Person struct {
	Id   int
	Name string
}

// изменяет оригинальную структуру
func (p *Person) SetName(name string) {
	p.Name = name
}

type Account struct {
	Id   int
	Name string
	Person
}

//func (p *Account) SetName(name string) {
//	p.Name = name
//}

func main() {
	var acc Account = Account{
		Id:   1,
		Name: "rvasily",
		Person: Person{
			Id:   2,
			Name: "Vasily Romanov",
		},
	}

	acc.SetName("romanov.vasily")
	//acc.Person.SetName("Test")

	fmt.Printf("updated account: %v\n", acc)

	

	// fmt.Printf("%#v \n", acc)

}
