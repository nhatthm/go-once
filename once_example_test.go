package once_test

import (
	"fmt"

	"go.nhat.io/once"
)

func ExampleLazyValueMap() {
	type Person struct {
		ID   string
		Name string
	}

	people := once.LazyValueMap[string, *Person]{
		New: func(key string) *Person {
			return &Person{ID: key}
		},
	}

	instance1 := people.Get("1")
	instance2 := people.Get("1")

	fmt.Println(instance2.Name)

	instance1.Name = "John Doe"

	fmt.Println(instance2.Name)

	// Output:
	//
	// John Doe
}
