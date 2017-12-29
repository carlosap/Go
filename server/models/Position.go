package models

const (
	 PositionsLinkType = "positions"
)

//Position entity
type Position struct {
	ID       string `json:"fileID"`
	Starting int    `json:"starting"`
	Length   int    `json:"length"`
	Text     string `json:"text"`
	Type     string `json:"type"`
}

//Positions Array Collection
type Positions []Position

func (p Positions) Len() int {
	return len(p)
}

func (p Positions) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p Positions) Less(i, j int) bool {
	return p[i].ID < p[j].ID
}