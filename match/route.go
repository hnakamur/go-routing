// Go HTTP router based on a simple custom match() function

package match

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type app struct {
	someConfig string
}

var theApp = app{
	someConfig: "value1",
}

func Serve(w http.ResponseWriter, r *http.Request) {
	var h http.Handler
	var apiWidget apiWidget
	var apiWidgetPart apiWidgetPart
	var widget widget

	p := r.URL.Path
	switch {
	case match(p, "/"):
		h = get(home)
	case match(p, "/contact"):
		h = get(contact)
	case match(p, "/api/widgets"):
		if r.Method == "GET" {
			h = get(apiGetWidgets)
		} else {
			h = post(apiCreateWidget)
		}
	case match(p, "/api/widgets/+", &apiWidget.slug):
		h = post(appWithApiWidget{app: &theApp, urlParams: &apiWidget}.update)
	case match(p, "/api/widgets/+/parts", &apiWidget.slug):
		h = post(appWithApiWidget{app: &theApp, urlParams: &apiWidget}.createPart)
	case match(p, "/api/widgets/+/parts/+/update", &apiWidgetPart.slug, &apiWidgetPart.id):
		h = post(appWithApiWidgetPart{app: &theApp, urlParams: &apiWidgetPart}.update)
	case match(p, "/api/widgets/+/parts/+/delete", &apiWidgetPart.slug, &apiWidgetPart.id):
		h = post(appWithApiWidgetPart{app: &theApp, urlParams: &apiWidgetPart}.delete)
	case match(p, "/+", &widget.slug):
		h = get(appWithWidget{app: &theApp, urlParams: &widget}.widget)
	case match(p, "/+/admin", &widget.slug):
		h = get(appWithWidget{app: &theApp, urlParams: &widget}.admin)
	case match(p, "/+/image", &widget.slug):
		h = post(appWithWidget{app: &theApp, urlParams: &widget}.image)
	default:
		http.NotFound(w, r)
		return
	}
	h.ServeHTTP(w, r)
}

// match reports whether path matches the given pattern, which is a
// path with '+' wildcards wherever you want to use a parameter. Path
// parameters are assigned to the pointers in vars (len(vars) must be
// the number of wildcards), which must be of type *string or *int.
func match(path, pattern string, vars ...interface{}) bool {
	for ; pattern != "" && path != ""; pattern = pattern[1:] {
		switch pattern[0] {
		case '+':
			// '+' matches till next slash in path
			slash := strings.IndexByte(path, '/')
			if slash < 0 {
				slash = len(path)
			}
			segment := path[:slash]
			path = path[slash:]
			switch p := vars[0].(type) {
			case *string:
				*p = segment
			case *int:
				n, err := strconv.Atoi(segment)
				if err != nil || n < 0 {
					return false
				}
				*p = n
			default:
				panic("vars must be *string or *int")
			}
			vars = vars[1:]
		case path[0]:
			// non-'+' pattern byte must match path byte
			path = path[1:]
		default:
			return false
		}
	}
	return path == "" && pattern == ""
}

func allowMethod(h http.HandlerFunc, method string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if method != r.Method {
			w.Header().Set("Allow", method)
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h(w, r)
	}
}

func get(h http.HandlerFunc) http.HandlerFunc {
	return allowMethod(h, "GET")
}

func post(h http.HandlerFunc) http.HandlerFunc {
	return allowMethod(h, "POST")
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "home\n")
}

func contact(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "contact\n")
}

func apiGetWidgets(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "apiGetWidgets\n")
}

func apiCreateWidget(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "apiCreateWidget\n")
}

type apiWidget struct {
	slug string
}

type appWithApiWidget struct {
	*app
	urlParams *apiWidget
}

func (h appWithApiWidget) update(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "apiUpdateWidget %s\n", h.urlParams.slug)
}

func (h appWithApiWidget) createPart(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "apiCreateWidgetPart %s\n", h.urlParams.slug)
}

type apiWidgetPart struct {
	slug string
	id   int
}

type appWithApiWidgetPart struct {
	*app
	urlParams *apiWidgetPart
}

func (h appWithApiWidgetPart) update(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "apiUpdateWidgetPart %s %d\n", h.urlParams.slug, h.urlParams.id)
}

func (h appWithApiWidgetPart) delete(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "apiDeleteWidgetPart %s %d\n", h.urlParams.slug, h.urlParams.id)
}

type widget struct {
	slug string
}

type appWithWidget struct {
	*app
	urlParams *widget
}

func (h appWithWidget) widget(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "widget %s\n", h.urlParams.slug)
}

func (h appWithWidget) admin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "widgetAdmin %s\n", h.urlParams.slug)
}

func (h appWithWidget) image(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "widgetImage %s\n", h.urlParams.slug)
}
