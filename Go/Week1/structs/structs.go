package mai22n

import "fmt"

type person struct {
	Id      int
	name    string
	address string
}

type account struct {
	Id int
	// Name    string
	Cleaner func(string) string
	owner   person
	person
}

func main() {
	// полное объявление структуры
	var acc account = account{
		Id: 1,
		// Name: "rvasily",
		person: person{
			name:    "Василий",
			address: "Москва",
		},
	}
	fmt.Printf("%#v\n", acc)

	// короткое объявление структуры
	acc.owner = person{2, "Romanov Vasily", "Moscow"}

	fmt.Printf("%#v\n", acc)

	fmt.Println(acc.name)
	fmt.Println(acc.person.name)
}
