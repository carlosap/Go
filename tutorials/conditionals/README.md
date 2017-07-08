# Conditionals

* if must be boolean expressions

* boolean comaparison 

```javascript
package main
import (
	"fmt"
)

func main(){
	//if must be boolean expressions
	//if <boolean expresssion>{....}
	//boolean comaparison 
	// == equal
	// != not equal
	// < Less than
	// <= Less than or equal
	// > Greater than
	// >= Greater than or equal
	// && AND
	// || OR

	firstrank := "39"
	secondrank := "614"
	if firstrank < secondrank{
		fmt.Println("\nRank 1 is Less Rank 2")
	}else if firstrank > secondrank{
		fmt.Println("\nRank 1 is Greater Rank 2")
	}else{
		fmt.Println("\nBoth Ranks are the same")
	}



}
```

### Simple Initialization Statements


```javascript
package main
import (
	"fmt"
)

func main(){

	if firstrank, secondrank :=39,600; firstrank < secondrank{
		fmt.Println("\nRank 1 is Less Rank 2")
	}else if firstrank > secondrank{
		fmt.Println("\nRank 1 is Greater Rank 2")
	}else{
		fmt.Println("\nBoth Ranks are the same")
	}
}
```

### Switch Statements

* fallthrough keyword forces the next case 
* no need to implicit  "break"

```javascript
package main
import (
	"fmt"
	"math/rand"
	"time"
)

func main(){
	switch "docker" {
		case "linux":
			fmt.Println("course is linux")
		case "docker":
			fmt.Println("course is docker")
			fallthrough
		case "windows":
			fmt.Println("course is windows")
		default:
			fmt.Println("sorry no OS found")

	}

	switch tmpNum := random(); tmpNum{
		case 0,2,4,6,8:
			fmt.Println("We got Even Numbers: ", tmpNum)
		case 1,3,5,7,9:
			fmt.Println("We got odd Numbers: ", tmpNum)
	}
}

func random() int{
	rand.Seed(time.Now().Unix())
	return rand.Intn(10)

}
```

### Errors Statements
* Errors idiomatic to return an "error" as the last
* nil is used to indicate success
* idiomatic to always check for the value of returned errors

```javascript
package main
import (
	"fmt"
	"os"
)

func main(){
	_,err := os.Open("c:\\temp\\servicelog1.txt")
	if err != nil {
		fmt.Println("Error: ", err)
	}else{
		//fmt.Println("Results:", _)
	}
}
```
