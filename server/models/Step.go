package models

//Order determines which order the steps should apear in.
type Order int
//OrderSteps sorts Steps by order
type OrderSteps struct {
	Steps
}

//Step represents a step in the data pipeline process
type Step struct {
	ID    string `json:"id"`
	Order int    `json:"order"`
}

//Steps Type
type Steps []Step

func (s Steps) Len() int {
	return len(s)
}

func (s Steps) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Steps) Less(i, j int) bool {
	return s[i].ID < s[j].ID
}



