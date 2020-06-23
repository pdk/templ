package main

import (
	"fmt"
	"html/template"
	"log"
	"os"
)

// Templates:
// 1. page layout (very few)
// 2. page content (one per page)
// 3. components (collection of tools to be re-used)

// map which page use which layout

type TemplateLibrary map[string]*template.Template

var pageLayouts = map[string]string{
	"contact": "special",
}

var funcMap = template.FuncMap(map[string]interface{}{
	"htmlSafe": func(text string) template.HTML {
		return template.HTML(text)
	},
})
var components *template.Template

func (tl TemplateLibrary) Page(pageName string) *template.Template {

	t, ok := tl[pageName]
	if ok {
		return t
	}

	if components == nil {
		var err error
		components, err = template.ParseGlob("templates/components/*.html")
		if err != nil {
			log.Fatalf("error parsing templates/components/*.html: %v", err)
		}
		components.Funcs(funcMap)
	}

	layoutName, ok := pageLayouts[pageName]
	if !ok {
		layoutName = "default"
	}

	pageTemplate := template.Must(components.Clone())
	pageTemplate, err := pageTemplate.ParseFiles("templates/layouts/"+layoutName+".html", "templates/pages/"+pageName+".html")
	if err != nil {
		log.Fatalf("error parsing templates/layouts/%s.html: %v", layoutName, err)
	}

	tl[pageName] = pageTemplate

	return pageTemplate
}

func ConfigureServer() {

	Reg(Post("/newThing", "newThing", OneThingHandler).WithLayout("special"))
	Reg(Get("/newThing", "newThing", AnotherThingHandler).WithLayout("special"))

}

func Reg(x RequestSpec) {

}

type PageHandler func()

func Get(path, pageName string, pageHandler PageHandler) RequestSpec {
	return RequestSpec{
		Method:   "GET",
		Path:     path,
		PageName: pageName,
		Handler:  pageHandler,
		Layout:   "default",
	}
}

func Post(path, pageName string, pageHandler PageHandler) RequestSpec {
	return RequestSpec{
		Method:   "POST",
		Path:     path,
		PageName: pageName,
		Handler:  pageHandler,
		Layout:   "default",
	}
}

type RequestSpec struct {
	Method   string
	Path     string
	PageName string
	Handler  PageHandler
	Layout   string
}

func (rs RequestSpec) WithLayout(layoutName string) RequestSpec {
	rs.Layout = layoutName
	return rs
}

func OneThingHandler() {}

func AnotherThingHandler() {}

type CheckBox struct {
	ID       string
	Label    string
	Checked  bool
	Class    string
	Disabled bool
}

func (cb CheckBox) WithDisabled(v bool) CheckBox {
	cb.Disabled = v
	return cb
}

func Ckbx(id, label string, checked bool, class string) CheckBox {
	return CheckBox{
		ID:       id,
		Label:    label,
		Checked:  checked,
		Class:    class,
		Disabled: false,
	}
}

func main() {

	lib := make(TemplateLibrary)

	fmt.Println("=== about ===")
	t1 := lib.Page("about")
	t1.ExecuteTemplate(os.Stdout, "page", nil)

	fmt.Println("=== contact ===")
	t2 := lib.Page("contact")
	t2.ExecuteTemplate(os.Stdout, "page", nil)

	// boxes := []CheckBox{
	// 	Ckbx("red", "Red", false, ""),
	// 	Ckbx("ylw", "Yellow", true, "filled-in"),
	// 	Ckbx("fldin", "Filled In", true, "filled-in"),
	// 	Ckbx("ind", "Indeterminate Style", true, ""),
	// 	Ckbx("grn", "Green", false, "").WithDisabled(true),
	// 	Ckbx("brn", "Brown", true, "").WithDisabled(true),
	// }

	// templates := template.Must(template.ParseGlob("templates/*.html"))

	// templates.ExecuteTemplate(os.Stdout, "checkboxes.html", map[string]interface{}{
	// 	"intro":  "world",
	// 	"checks": boxes,
	// })
}
