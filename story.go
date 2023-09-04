package cyoa

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
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
	story  map[string]Chapter
	tmpl   *template.Template
	pathFn func(r *http.Request) string
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
		<section class="page">
			<h1>{{.Title}}</h1>
			{{range .Paragraphs}}
				<p>{{.}}</p>
			{{end}}
			<ul>
				{{range .Options}}
					<li><a href="/{{.Chapter}}">{{.Text}}</a></li>
				{{end}}
			</ul>
		</section>
		<style>
		body {
		  font-family: helvetica, arial;
		}
		h1 {
		  text-align:center;
		  position:relative;
		}
		.page {
		  width: 80%;
		  max-width: 500px;
		  margin: auto;
		  margin-top: 40px;
		  margin-bottom: 40px;
		  padding: 80px;
		  background: #FECEA8;
		  border: 1px solid #eee;
		  box-shadow: 0 10px 6px -6px #777;
		}
		ul {
		  border-top: 1px dotted #ccc;
		  padding: 10px 0 0 0;
		  -webkit-padding-start: 0;
		}
		li {
		  padding-top: 10px;
		}
		a,
		a:visited {
		  text-decoration: none;
		  color: #6295b5;
		}
		a:active,
		a:hover {
		  color: #7792a2;
		}
		p {
		  text-indent: 1em;
		}
	  </style>
	</body>
</html>
`
var tpl *template.Template

func init() {
	tpl = template.Must(template.New("").Parse(defaultHandlerTmpl))
}

func NewHandler(s map[string]Chapter, opts ...HandlerOpt) http.Handler {
	h := handler{story: s, tmpl: tpl, pathFn: defaultPathFn}
	for _, opt := range opts {
		opt(&h)
	}
	return h
}

type HandlerOpt func(h *handler)

func WithTemplate(t *template.Template) HandlerOpt {
	return func(h *handler) {
		h.tmpl = t
	}
}

func WithPathFn(fn func(r *http.Request) string) HandlerOpt {
	return func(h *handler) {
		h.pathFn = fn
	}
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := h.pathFn(r)
	if chapter, ok := h.story[path]; ok {
		err := h.tmpl.Execute(w, chapter)
		if err != nil {
			log.Fatal("could not execute template", err)
			http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "Chapter not found.", http.StatusNotFound)
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

func defaultPathFn(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "" || path == "/" {
		path = "/intro"
	}
	return path[1:] // gets rid of the slash
}

