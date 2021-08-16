// Go HTTP router based on a simple custom match() function

package match

import (
	"fmt"
	"net/http"
	"strconv"
)

func Serve(w http.ResponseWriter, r *http.Request) {
	var h http.Handler
	var slug string
	var id int

	p := r.URL.Path
	switch {
	case match(p, "/"):
		h = get(home)
	case match(p, "/contact"):
		h = get(contact)
	case match(p, "/api/widgets") && r.Method == "GET":
		h = get(apiGetWidgets)
	case match(p, "/api/widgets"):
		h = post(apiCreateWidget)
	case match(p, "/api/widgets/+", &slug):
		h = post(apiWidget{slug}.update)
	case match(p, "/api/widgets/+/parts", &slug):
		h = post(apiWidget{slug}.createPart)
	case match(p, "/api/widgets/+/parts/+/update", &slug, &id):
		h = post(apiWidgetPart{slug, id}.update)
	case match(p, "/api/widgets/+/parts/+/delete", &slug, &id):
		h = post(apiWidgetPart{slug, id}.delete)
	case match(p, "/+", &slug):
		h = get(widget{slug}.widget)
	case match(p, "/+/admin", &slug):
		h = get(widget{slug}.admin)
	case match(p, "/+/image", &slug):
		h = post(widget{slug}.image)
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
	path = path[1:]
	pattern = pattern[1:]
	var pathEndsWithSlash bool
	var patternEndsWithSlash bool
	for path != "" && pattern != "" {
		var pathSegment string
		for i := 0; i < len(path); i++ {
			if path[i] == '/' {
				if i == len(path)-1 {
					pathEndsWithSlash = true
				}
				pathSegment = path[:i]
				path = path[i+1:]
				break
			}
		}
		if pathSegment == "" {
			pathSegment = path
			path = ""
		}

		var patternSegment string
		for i := 0; i < len(pattern); i++ {
			if pattern[i] == '/' {
				if i == len(pattern)-1 {
					patternEndsWithSlash = true
				}
				patternSegment = pattern[:i]
				pattern = pattern[i+1:]
				break
			}
		}
		if patternSegment == "" {
			patternSegment = pattern
			pattern = ""
		}

		if patternSegment == "+" {
			// '+' matches till next slash in path
			switch p := vars[0].(type) {
			case *string:
				*p = pathSegment
			case *int:
				n, err := strconv.Atoi(pathSegment)
				if err != nil || n < 0 {
					return false
				}
				*p = n
			default:
				panic("vars must be *string or *int")
			}
			vars = vars[1:]
		} else if patternSegment != pathSegment {
			return false
		}
	}
	return path == "" && pattern == "" &&
		pathEndsWithSlash == patternEndsWithSlash
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

func (h apiWidget) update(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "apiUpdateWidget %s\n", h.slug)
}

func (h apiWidget) createPart(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "apiCreateWidgetPart %s\n", h.slug)
}

type apiWidgetPart struct {
	slug string
	id   int
}

func (h apiWidgetPart) update(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "apiUpdateWidgetPart %s %d\n", h.slug, h.id)
}

func (h apiWidgetPart) delete(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "apiDeleteWidgetPart %s %d\n", h.slug, h.id)
}

type widget struct {
	slug string
}

func (h widget) widget(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "widget %s\n", h.slug)
}

func (h widget) admin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "widgetAdmin %s\n", h.slug)
}

func (h widget) image(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "widgetImage %s\n", h.slug)
}
