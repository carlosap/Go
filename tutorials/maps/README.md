# Maps


* github.com/docker/docker/blob/master/graph/graph.go
* maps are unordered list (results are ramdom)
* maps are key>:value pairs - Dictionaries/hashtable
* maps are dynamically resizable
* maps are references
* make(map[string]int, size) can help to improve performance


```javascript
package main
import (
	"fmt"
)

func main(){
	titles := make(map[string]int)
	titles["windows"] = 6
	titles["linux"] = 2
	recentTitles := map[string] int{
		"windows": 5,
		"linux": 0
	}

	//pritn single value- returns 0
	fmt.Println(recentTitles["linux"])
	//update the map
	recentTitles["linux"] = 100
	fmt.Println(recentTitles["linux"])

	//adding a new one
	//if Mac already exist it just updates them else it adds new
	recentTitles["Mac"] = 19

	//delete
	delete(recentTitles, "Mac")
    
	//shows unordered list - 
	for key, value := range recentTitles{
		fmt.Printf("\nKey is: %v Value is: %v", key, value)
	}

}

```