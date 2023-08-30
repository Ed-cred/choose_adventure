package cyoa

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"text/template"
)

type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

type Option struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}

type handler struct {
	story map[string]Chapter
}

var defaultHandlerTmpl = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="utf-8">
		<title>Choose Your Own Adventure</title>
		<meta name="viewport" content="width=device-width, initial-scale=1">
	</head>
	<body>
		<h1>{{.Title}}</h1>
		{{range .Paragraphs}}
			<p>{{.}}</p>
		{{end}}
		<ul>
			{{range .Options}}
				<li><a href="/{{.Chapter}}">{{.Text}}</a></li>
			{{end}}
		</ul>
	</body>
</html>
`
var tmpl *template.Template

func init() {
	tmpl = template.Must(template.New("").Parse(defaultHandlerTmpl))
}

func NewHandler(s map[string]Chapter) http.Handler {
	return handler{story: s}
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := tmpl.Execute(w, h.story["intro"])
	if err != nil {
		log.Fatal("could not execute template", err)
	}
}

func JsonStoryFromFile(file *os.File) (map[string]Chapter, error) {
	story := make(map[string]Chapter)
	err := json.NewDecoder(file).Decode(&story)
	if err != nil {
		log.Fatal("failed to decode json", err)
		return nil, err
	}
	return story, nil
}
