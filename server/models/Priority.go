package models

//Note structure
type Priority struct {
	ID    string `json:"id"`
	Order int    `json:"order"`
}

//Priorities are a list of Priority to support sorting
type Priorities []Priority

func (p Priorities) Len() int {
	return len(p)
}

func (p Priorities) Less(i, j int) bool {
	return p[i].ID < p[j].ID
}

func (p Priorities) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}