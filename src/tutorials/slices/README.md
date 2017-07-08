# Slices


* numbered lists of single type (think of List)
* Can be resized
* Slices are built on top of arrays
* they are references
* changes values on slices changes values on array
* flexible length


```javascript
package main
import (
	"fmt"
)

func main(){

	mySlice := []int{1,2,3,4,5}
	mySlice[1] = 0;
	fmt.Println(mySlice)

	sliceOfSlice := mySlice[2:5]
	fmt.Println(sliceOfSlice)

	//growth slice
	courses := make([]string, 5,10)
	mycourses := []string {"docker","linux","windows"}
	fmt.Printf("Length is: %d. \nCapacity is: %d", len(courses),cap(courses))
	fmt.Printf("Length is: %d. \nCapacity is: %d", len(mycourses),cap(mycourses))

	for i := 1; i < 17; i++ {
		mySlice = append(mySlice, i)
		fmt.Printf("\nCapacity is: %d", cap(mySlice))
	}

	for _, i := range mySlice{
		fmt.Println("for range loop", i)
	}

	//appending to existing Slice
	newSlice := []int{10,20,30}
	mySlice = append(mySlice, newSlice...)
	fmt.Println(newSlice)
}

```