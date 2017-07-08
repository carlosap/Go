# Constants

* Constants are declared like variables, but with the const keyword.

* Constants can be character, string, boolean, or numeric values.

* Constants cannot be declared using the := syntax.

```javascript
package main
import (
	"fmt"
	"os"
)
const Pi = 3.14
func main() {
	const World = "World"
	fmt.Println("Hello", World)
    fmt.Println("Happy", Pi, "Day")
    
	name := os.Getenv("USERNAME")
	fmt.Println(name)
	for _, env := range os.Environ(){
		fmt.Println(env)
	}
}
```