// static-files.go
package main

import (
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
)

var templatesPath = "templates/*.html"
var tpl *template.Template

type Product struct {
	Name string
	URL string
	Price string
	Properties []string
}

func main() {


	var err error
	tpl, err = template.ParseGlob(templatesPath)
	if err != nil {
		log.Fatal(err)
	}
	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/about", AboutHandler)

	log.Fatal(http.ListenAndServe(":8080", r))
}

func AboutHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	err := tpl.ExecuteTemplate(w, "about", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {

	products := []Product{
		{
			"CP2140 / CP3140",
			"https://web.archive.org/web/20181106170520im_/http://kazenbarcode.com/product_image/1cat2.jpg",
			"Rs 399",
			[]string{
				"Printing Resolution: 203 DPI",
				"Print Speed: 4 Inches per Second",
				"Ribbon Length: 300MTRS.",
				"Interface: Serial(RS232), Parallel , USB",
				"Memory : 8MB DRAM, 4MB Flash ROM",
				"Downloable Fonts….Movable Sensor",
				"Zebra Emulation",
			},
		},{
			"CP2140 / CP3140",
			"https://web.archive.org/web/20181106170520im_/http://kazenbarcode.com/product_image/1cat2.jpg",
			"Rs 399",
			[]string{
				"Printing Resolution: 203 DPI",
				"Print Speed: 4 Inches per Second",
				"Ribbon Length: 300MTRS.",
				"Interface: Serial(RS232), Parallel , USB",
				"Memory : 8MB DRAM, 4MB Flash ROM",
				"Downloable Fonts….Movable Sensor",
				"Zebra Emulation",
			},
		},{
			"CP2140 / CP3140",
			"https://web.archive.org/web/20181106170520im_/http://kazenbarcode.com/product_image/1cat2.jpg",
			"Rs 399",
			[]string{
				"Printing Resolution: 203 DPI",
				"Print Speed: 4 Inches per Second",
				"Ribbon Length: 300MTRS.",
				"Interface: Serial(RS232), Parallel , USB",
				"Memory : 8MB DRAM, 4MB Flash ROM",
				"Downloable Fonts….Movable Sensor",
				"Zebra Emulation",
			},
		},{
			"CP2140 / CP3140",
			"https://web.archive.org/web/20181106170520im_/http://kazenbarcode.com/product_image/1cat2.jpg",
			"Rs 399",
			[]string{
				"Printing Resolution: 203 DPI",
				"Print Speed: 4 Inches per Second",
				"Ribbon Length: 300MTRS.",
				"Interface: Serial(RS232), Parallel , USB",
				"Memory : 8MB DRAM, 4MB Flash ROM",
				"Downloable Fonts….Movable Sensor",
				"Zebra Emulation",
			},
		},{
			"CP2140 / CP3140",
			"https://web.archive.org/web/20181106170520im_/http://kazenbarcode.com/product_image/1cat2.jpg",
			"Rs 399",
			[]string{
				"Printing Resolution: 203 DPI",
				"Print Speed: 4 Inches per Second",
				"Ribbon Length: 300MTRS.",
				"Interface: Serial(RS232), Parallel , USB",
				"Memory : 8MB DRAM, 4MB Flash ROM",
				"Downloable Fonts….Movable Sensor",
				"Zebra Emulation",
			},
		},
	}

	type Result struct {
		Products []Product
	}

	result := Result{
		products,
	}

	w.WriteHeader(http.StatusOK)
	err := tpl.ExecuteTemplate(w, "home", result)
	if err != nil {
		log.Fatal(err)
	}

}
