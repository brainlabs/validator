## Simple Struct Validator

### example :
```go
package main

import (
	"fmt"

	"github.com/brainlabs/validator"
)

type Person struct {
	Name       string   `json:"name" valid:"required|min:3|max:10|alpha_space"`
	Age        int      `json:"age" valid:"required|min:4|max:100"`
	Status     string   `json:"status" valid:"in:success,failed"`
	Address    *Address `json:"address" valid:"required"`
	Phone      string   `json:"phone" valid:"required|id_phone"`
	Validation func()
}

type Address struct {
	AddressName string `json:"address_name" valid:"required|min:10|max:100"`
}



func main() {
	
	ps := &Person{
		Name:   "Jhon Doe",
		Age:    10,
		Status: "success",
		Address: &Address{
			//AddressName: "Street Walker Petir Jakarta No.20",
		},
		Phone: "62821020102010201",
	}

	vl := validator.New()
	result := vl.ValidateStruct(ps)

	fmt.Println(validator.DumpToString(result))

}
```




### Author