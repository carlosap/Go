# Functions

Two parameter and returns two strings

```javascript
func converter(s1, s2 string) (str1, str2 string){
	s1 = strings.Title(s1)
	s2 = strings.ToUpper(s2)
	return s1, s2
}
```
When the number of input paramters are unknown (variadic function) use the spread operator

 * Add '&' to reference a pointer
```javascript
func bestLeagueFinishes(finishes ...int) int{
	best := finishes[0]
	for _, i := range finishes {
		if i < best {
			best = i
		}
	}

	return best
}
```