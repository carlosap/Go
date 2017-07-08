# Variables

Basic useful details about Go:

 * Strongtly Type Language
 * Smart to know types base on variables (name, module = "carlos", 3.2)
 * Declaring 'var' at the package level are Global
 * Variable initializer ":=" only works inside functions

Example of declaring variables with types:

```javascript
var (
	name string
    module float64
)
```

Basic Pointers

 * Add '&' to reference a pointer
```javascript
ptr := &module
```
 * De-referencing a pointer means getting the value content use '*' from memory

```javascript
package main
import (
	"fmt"
	"reflect" 
)

var (
	name = "carlos"
	module = 3.2
)

func main(){
	a := 10.0000000
	b := 3
	c := int(a) + b
	fmt.Println("The sum of C is: ", c)
	fmt.Println("Name is a type of", reflect.TypeOf(name))
	fmt.Println("Module is type of", reflect.TypeOf(module))
	fmt.Println("Name is set to ", name)
	fmt.Println("Module is set to", module)
	ptr := &module
	fmt.Println("memory address of *module* variable is: ", ptr)
	fmt.Println("and the value of *module* is ", *ptr)

}
```



### Re-Assigning Variables:

Go uses the "=" to reasign the varible. The code sample below re-assigns the course value.

```javascript
package main
import (
	"fmt"
)

func main(){
	name := "carlos"
	course := "Go Deep Dive"
	fmt.Println("\nHi",name, "You are now watching course: ",course)
	changeCourse(course)
	fmt.Println("\nYou are now watching course", course)
}

func changeCourse(course string) string{
	course = "First Look: Native Go Services"
	fmt.Println("\nTrying to change course to", course)
	return course
}
```
### Using '&' ref by and '*' dereference by:


```javascript
package main
import (
	"fmt"
)

func main() {
	name := "carlos"
	course := "Deep Dive Go"

	fmt.Println("\nHi", name, "You are watching", course)
	changeCourse(&course)
	fmt.Println("\nYou are now watching course", course)

}

func changeCourse(course *string) string{
	*course = "Native C# Clustering"
	fmt.Println("\nTrying to change your course to ", *course)
	return *course
}
```
