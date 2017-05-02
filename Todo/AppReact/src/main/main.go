package main 
import (
	"net/http"
	"os"
	"text/template"
	"controllers"
)

func main(){

	templates := getTemplates()
	controllers.Register(templates)
	http.ListenAndServer(":8000",nil)

}

func getTemplates() *template.Template{
	result := template.New("templates")
	basePath : "templates"

}