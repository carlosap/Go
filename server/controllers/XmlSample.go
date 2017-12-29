package controllers
import (
	"encoding/xml"
	"os"
	"io/ioutil"
	"fmt"
	//"strings"
	"net/http"
	"github.com/martini-contrib/encoder"
	"github.com/go-martini/martini"
	
) 

type Query struct {
  XMLName xml.Name `xml:"feed"`
	ID string `xml:"id"`
	Title string `xml:"title"`
	LastUpdate string `xml:"updated"`
	Author string `xml:"author>name"`
	TotalResults string `xml:"totalResults"`
	StartIndex string `xml:"startIndex"`
	ItemsPerPage string `xml:"itemsPerPage"`
	EntryIDs  []string `xml:"entry>id"`
	EntryLinks []Link `xml:"entry>link"`
	EntryTitles []string `xml:"entry>title"`
	EntryUpdates []string `xml:"entry>updated"`
	EntryScores []string `xml:"entry>score"`
}

// Xml Serialization XmlAttributeAttribute
type Link struct {
	Href   string `xml:"href,attr"`
}

func RegisterXmlEndpoints(h http.Handler) {
	m := h.(*martini.ClassicMartini)
	m.Group("/xml", func(r martini.Router) {
		m.Get("", GetXmlHandler)
	})
}

func GetXmlHandler(r *http.Request, w http.ResponseWriter,enc encoder.Encoder) (int, []byte) {

var q Query

xmlFile, err := os.Open("Harmony1.xml")
if err != nil {	
	fmt.Println("Error opening file: ", err)
	return http.StatusOK, encoder.Must(enc.Encode(q))
}

defer xmlFile.Close()


b, _ := ioutil.ReadAll(xmlFile)

xml.Unmarshal(b, &q)

fmt.Println(q.Title)

for _,item := range q.EntryLinks {
		fmt.Println(item)
}

	return http.StatusOK, encoder.Must(enc.Encode(q))
}
