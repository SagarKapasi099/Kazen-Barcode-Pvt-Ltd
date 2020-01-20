// static-files.go
package main

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"math/rand"
	"time"

	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"html/template"
	"log"
	"net/http"
)

var templatesPath = "templates/*.html"
var tpl *template.Template

type Product struct {
	Id         string   `json:"id"`
	Name       string   `json:"name"`
	URL        string   `json:"url"`
	Properties []string `json:"properties"`
	Price      int      `json:"price"`
	Active     bool     `json:"active"`
}

type Enquiry struct {
	Name      string   `json:"name"`
	Email     string   `json:"email"`
	Mobile    string   `json:"mobile"`
	Comments  string   `json:"comments"`
	OTP       string   `json:"otp"`
	ProductId []string `json:"product_id"`
}

func main() {

	var err error
	tpl, err = template.ParseGlob(templatesPath)
	if err != nil {
		log.Fatal(err)
	}

	// checking to see if database connection is alive
	c := GetClient()
	pingError := c.Ping(context.Background(), readpref.Primary())
	if pingError != nil {
		log.Fatal("Couldn't connect to the database", err)
	} else {
		log.Println("Connected!")
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

// GetClient returns a MongoDB Client Singleton
func GetClient() *mongo.Client {
	clientOptions := options.Client().ApplyURI("mongodb+srv://SagarKapasi099:3wqzTsSvNQkovuxi@projectautodidact-5vr5f.gcp.mongodb.net/test?retryWrites=true&w=majority")
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func GenerateOTP(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func AboutHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	err := tpl.ExecuteTemplate(w, "about", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	var products []Product
	filter := bson.M{"active": true}
	client := GetClient()
	collection := client.Database("kbpl").Collection("products")
	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Fatal("Error on Finding all the documents", err)
	}

	for cur.Next(context.TODO()) {
		var product Product
		err = cur.Decode(product)
		if err != nil {
			log.Fatal("Error on Decoding the document", err)
		}
		products = append(products, product)
	}
	/*
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
	*/
	type Result struct {
		Products []Product
	}

	result := Result{
		products,
	}

	w.WriteHeader(http.StatusOK)
	err = tpl.ExecuteTemplate(w, "home", result)
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
	productIdsJson := r.FormValue("id")

	var productIds []string
	if err := json.Unmarshal([]byte(productIdsJson), &productIds); err != nil {
		log.Println(err)
	}

	// TODO: generate OTP
	otp := GenerateOTP(6)

	currentEnquiry := Enquiry{name, email, number, comments, otp, productIds}
	fmt.Println(currentEnquiry)

	client := GetClient()
	collection := client.Database("kbpl").Collection("enquiries")
	res, err := collection.InsertOne(context.TODO(), currentEnquiry)
	if err != nil {
		fmt.Println(err)
	}
	enquiryId := res.InsertedID.(primitive.ObjectID) // customer id from database

	jsonResponse := `{
		"success": true,
		"message":{
			"superText":"Please Enter OTP Sent On ` + number + `",
			"subText": "",
			"buttonText": "Verify"
		},
		"number": "` + number + `",
		"enquiryId": "` + enquiryId.Hex() + `",
		"method": "post"
	}`

	n, err := fmt.Fprintf(w, jsonResponse)
	if err != nil {
		log.Fatal(n, err)
	}
}

func VerifyOTPHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: avoid brute force attack

	enquiryId := r.FormValue("enquiryId")
	primitiveValueOfEnquiryId, err := primitive.ObjectIDFromHex(enquiryId)
	if err != nil {
		log.Println(err)
	}

	// getting enquiry object from enquiryId
	var enquiry Enquiry
	client := GetClient().Database("kbpl").Collection("enquiries")
	documentReturned := client.FindOne(context.TODO(), bson.M{"_id": primitiveValueOfEnquiryId})
	err = documentReturned.Decode(&enquiry)
	if err != nil {
		log.Println(err)
	}

	otp := enquiry.OTP
	otpReceived := r.FormValue("otp")
	fmt.Println(otpReceived)
	showProducts := "false"
	if len(enquiry.ProductId) > 0 {
		showProducts = "true"
	}

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
			"enquiryId": "` + enquiryId + `",
			"showProducts": ` + showProducts + `
		}`
	} else {
		// OTP does not match
		jsonResponse = `{
			"success": false,
			"message": {
				"superText": "Something Went Wrong.",
				"subText": ""
			},
			"enquiryId": "` + enquiryId + `"
		}`
	}
	n, err := fmt.Fprintf(w, jsonResponse)
	if err != nil {
		log.Fatal(n, err)
	}
}

func ShowProductsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, r.FormValue("enquiryId"))
	// TODO extract a view from products grid and show the template here
	fmt.Fprintf(w, "here are your products")
}
