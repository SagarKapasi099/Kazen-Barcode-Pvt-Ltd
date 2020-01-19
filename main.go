// static-files.go
package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
)

var templatesPath = "templates/*.html"
var tpl *template.Template

type Product struct {
	Id         string
	Name       string
	URL        string
	Properties []string
}

type Enquiry struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Mobile   string `json:"mobile"`
	Comments string `json:"comments"`
	OTP      string `json:"otp"`
}

func main() {

	var err error
	tpl, err = template.ParseGlob(templatesPath)
	if err != nil {
		log.Fatal(err)
	}
	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	r.HandleFunc("/", HomeHandler).Methods("GET")
	r.HandleFunc("/about", AboutHandler).Methods("GET")
	r.HandleFunc("/contact", ContactUsHandler).Methods("GET")
	r.HandleFunc("/enquiry", EnquiryHandler).Methods("POST")
	r.HandleFunc("/verifyOTP", VerifyOTPHandler).Methods("POST")
	r.HandleFunc("/products", ShowProductsHandler).Methods("POST")
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
			"id-1",
			"CP2140 / CP3140",
			"https://web.archive.org/web/20181106170520im_/http://kazenbarcode.com/product_image/1cat2.jpg",
			[]string{
				"Printing Resolution: 203 DPI",
				"Print Speed: 4 Inches per Second",
				"Ribbon Length: 300MTRS.",
				"Interface: Serial(RS232), Parallel , USB",
				"Memory : 8MB DRAM, 4MB Flash ROM",
				"Downloable Fonts….Movable Sensor",
				"Zebra Emulation",
			},
		}, {
			"id-2",
			"CP2140 / CP3140",
			"https://web.archive.org/web/20181106170520im_/http://kazenbarcode.com/product_image/1cat2.jpg",
			[]string{
				"Printing Resolution: 203 DPI",
				"Print Speed: 4 Inches per Second",
				"Ribbon Length: 300MTRS.",
				"Interface: Serial(RS232), Parallel , USB",
				"Memory : 8MB DRAM, 4MB Flash ROM",
				"Downloable Fonts….Movable Sensor",
				"Zebra Emulation",
			},
		}, {
			"id-3",
			"CP2140 / CP3140",
			"https://web.archive.org/web/20181106170520im_/http://kazenbarcode.com/product_image/1cat2.jpg",
			[]string{
				"Printing Resolution: 203 DPI",
				"Print Speed: 4 Inches per Second",
				"Ribbon Length: 300MTRS.",
				"Interface: Serial(RS232), Parallel , USB",
				"Memory : 8MB DRAM, 4MB Flash ROM",
				"Downloable Fonts….Movable Sensor",
				"Zebra Emulation",
			},
		}, {
			"id-4",
			"CP2140 / CP3140",
			"https://web.archive.org/web/20181106170520im_/http://kazenbarcode.com/product_image/1cat2.jpg",
			[]string{
				"Printing Resolution: 203 DPI",
				"Print Speed: 4 Inches per Second",
				"Ribbon Length: 300MTRS.",
				"Interface: Serial(RS232), Parallel , USB",
				"Memory : 8MB DRAM, 4MB Flash ROM",
				"Downloable Fonts….Movable Sensor",
				"Zebra Emulation",
			},
		}, {
			"id-5",
			"CP2140 / CP3140",
			"https://web.archive.org/web/20181106170520im_/http://kazenbarcode.com/product_image/1cat2.jpg",
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

func ContactUsHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	err := tpl.ExecuteTemplate(w, "contact", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func EnquiryHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	email := r.FormValue("email")
	number := r.FormValue("number")
	comments := r.FormValue("comments")

	// TODO: generate OTP
	otp := "2995"

	currentEnquiry := Enquiry{name, email, number, comments, otp}
	fmt.Println(currentEnquiry)

	// TODO: save to database
	customerId := "cust-1" // customer id from database

	jsonResponse := `{
		"success": true,
		"message":{
			"superText":"Please Enter OTP Sent On `+number+`",
			"subText": "",
			"buttonText": "Verify"
		},
		"number": "`+number+`",
		"customerId": "`+ customerId +`",
		"url": "/verifyOTPHandler",
		"method": "post"
	}`

	_, err := fmt.Fprintf(w, jsonResponse)
	if err != nil {
		log.Fatal(err)
	}
}

func VerifyOTPHandler(w http.ResponseWriter, r *http.Request) {
	customerId := r.FormValue("customerId")
	// TODO get number from id and get OTP from id
	otp := "399"
	otpReceived := r.FormValue("otp")
	fmt.Println(otpReceived)
	// TODO: avoid brute force attack
	// TODO: check database for OTP
	jsonResponse := ``
	if otpReceived == otp {
		// OTP matches
		jsonResponse = `{
			"success": true,
			"message": {
				"superText": "Thank you for your enquiry.",
				"subText": "It has been forwarded to the relevant department and will be dealt with as soon as possible.",
				"buttonText": "Show Products With Prices"
			},
			"customerId": "`+customerId+`",
			"url": "/showProducts",
			"method": "post"
		}`
	} else {
		// OTP does not match
		jsonResponse = `{
			"success": false,
			"message": {
				"superText": "Something Went Wrong.",
				"subText": ""
			},
			"customerId": "`+customerId+`"
		}`
	}
	_, err := fmt.Fprintf(w, jsonResponse)
	if err != nil {
		log.Fatal(err)
	}
}

func ShowProductsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, r.FormValue("customerId"))
	// TODO extract a view from products grid and show the template here
	fmt.Fprintf(w, "here are your products")
}