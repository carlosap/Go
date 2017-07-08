# Structs


* defines custom data types
* OO programmer  Go way- Modular re-useble components
* In GO- No Objects, No Classes, No Inheritance PUT ASIDE C###### -:)


```javascript
package main
import (
	"fmt"
)

func main(){

	type PersonMeta struct {
		Name string
		LastName string
		Age float64
	}

	//define 
	var person PersonMeta
	person := new(PersonMeta)  //stars new values with empmty. this gives us a pointer

	//adding second person. literal way
	secondPerson := personMeta{
		Name: "Carlos",
		LastName: "Perez",
		Age: 19
	}

	fmt.Println(secondPerson)
	fmt.Println(secondPerson.Name)

}

```