# Loops

Go has only one looping construct, the for loop.

The basic for loop has three components separated by semicolons:

* the init statement: executed before the first iteration
* the condition expression: evaluated before every iteration
* the post statement: executed at the end of every iteration


The init statement will often be a short variable declaration, and the variables declared there are visible only in the scope of the for statement.

The loop will stop iterating once the boolean condition evaluates to false.

```javascript
package main
import (
	"fmt"
	"time"
)

func main(){
	for timer := 10; timer >=0; timer-- {
		if timer == 0 {
			fmt.Println("Boom")
			break
		}
		fmt.Println(timer)
		time.Sleep(1 * time.Second)
	}
}
```

### Range

range iterates over elements in a variety of data structures.

```javascript
package main
import (
	"fmt"
)

func main(){
	courses := []string{"docker","linux","windows"}
	completed := []string{"docker","windows"}
	for _, item := range courses{
		fmt.Println(item)
		for _, completedItem := range completed{
			if item == completedItem {
				fmt.Println("\n Found duplicate items with: ", item)
			}
		}
	}	
}
```

### Continue and Breaks

```javascript
package main
import (
	"fmt"
	"time"
)

func main(){
	for timer := 10; timer >= 0; timer-- {
		if timer % 2 == 0 {
			continue
		}
		fmt.Println(timer)
		time.Sleep(1 * time.Second)

	}
}
```


