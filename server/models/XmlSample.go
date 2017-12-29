package models

type Query struct {
	Series Show
	EpisodeList []Episode `xml:"Episode>"`
}


type Show struct {
	Title string `xml:"SeriesName>"`
	SeriesID int
	Keywords map[string] bool
}

type Episode struct {
	SeasonNumber int
	EpisodeNumber int
	EpisodeName string
	FirstAired string
}