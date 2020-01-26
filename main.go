// static-files.go
package main

import (
	"context"
	"encoding/json"
	"errors"
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt2 "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"io"
	"math/rand"
	"strconv"
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

const AppKey = "klshifhjKLHLHGl;sdjhfl'kshjfgkSghsFHJKSHGSHGslhgh"

type Product struct {
	Id         string   `json:"_id" bson:"_id"`
	Name       string   `json:"name" bson:"name"`
	URL        string   `json:"url" bson:"url"`
	Properties []string `json:"properties" bson:"properties"`
	Price      int      `json:"price" bson:"price"`
	Active     bool     `json:"active" bson:"active"`
}

type Enquiry struct {
	Name      string   `json:"name" bson:"name"`
	Email     string   `json:"email" bson:"email"`
	Mobile    string   `json:"mobile" bson:"mobile"`
	Comments  string   `json:"comments" bson:"comments"`
	OTP       string   `json:"-" bson:"otp"`
	ProductId []string `json:"-" bson:"product_id"`
}

type DataTableResponse struct {
	Draw int `json:"draw"`
	RecordsTotal int `json:"recordsTotal"`
	RecordsFiltered int `json:"recordsFiltered"`
	Data [][]string `json:"data"`
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

	// Admin Section
	r.HandleFunc("/manage", ManageHandler).Methods("GET")
	r.HandleFunc("/getToken", GenerateJWT)
	r.Handle("/administrator", AuthMiddleware(http.HandlerFunc(AdministratorHandler))).Methods("GET")
	// Enquiries
	r.Handle("/administrator/enquiries", AuthMiddleware(http.HandlerFunc(AdminEnquiriesHandler))).Methods("GET")
	r.Handle("/administrator/getEnquiriesJson", AuthMiddleware(http.HandlerFunc(AdminEnquiriesJsonHandler))).Methods("GET")
	r.Handle("/administrator/viewEnquiry/{id}", AuthMiddleware(http.HandlerFunc(AdminViewEnquiryHandler))).Methods("GET")
	r.Handle("/administrator/saveEnquiry", AuthMiddleware(http.HandlerFunc(AdminSaveEnquiryHandler))).Methods("POST")
	// Products
	r.Handle("/administrator/products", AuthMiddleware(http.HandlerFunc(AdminProductsHandler))).Methods("GET")
	r.Handle("/administrator/getProductsJson", AuthMiddleware(http.HandlerFunc(AdminProductsJsonHandler))).Methods("GET")
	r.Handle("/administrator/viewProduct/{id}", AuthMiddleware(http.HandlerFunc(AdminViewProductHandler))).Methods("GET")
	r.Handle("/administrator/saveProduct", AuthMiddleware(http.HandlerFunc(AdminSaveProductHandler))).Methods("POST")
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

// AuthMiddleware is our middleware to check our token is valid. Returning
// a 401 status to the client if it is not valid.
func AuthMiddleware(next http.Handler) http.Handler {

	if len(AppKey) == 0 {
		log.Fatal("HTTP server unable to start, expected an APP_KEY for JWT auth")
	}
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		CredentialsOptional: true,
		Extractor: func(r *http.Request) (string, error) {
			accessTokenCookie, err := r.Cookie("access_token")
			if err != nil {
				return "", errors.New("Error Obtaining Cookie access_token")
			}
			jwtCookie := accessTokenCookie.Value
			return jwtCookie, nil
		},
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, errMsg string) {
			http.Redirect(w, r, "/manage", http.StatusSeeOther)
		},
		ValidationKeyGetter: func(token *jwt2.Token) (interface{}, error) {
			return []byte(AppKey), nil
		},
		SigningMethod: jwt2.SigningMethodHS256,
	})
	return jwtMiddleware.Handler(next)
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
		err = cur.Decode(&product)
		if err != nil {
			log.Fatal("Error on Decoding the document", err)
		}
		products = append(products, product)
	}

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

	otp := GenerateOTP(6)

	currentEnquiry := Enquiry{name, email, number, comments, otp, productIds}

	client := GetClient()
	collection := client.Database("kbpl").Collection("enquiries")
	res, err := collection.InsertOne(context.TODO(), currentEnquiry)
	if err != nil {
		log.Println(err)
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
	enquiryId := r.FormValue("enquiryId")

	primitiveValueOfEnquiryId, err := primitive.ObjectIDFromHex(enquiryId)
	if err != nil {
		log.Println(err)
	}
	// TODO extract a view from products grid and show the template here
	// { _id : { $in : [ObjectId("5e27363e5e9f18b645d3e957")] } }
	client := GetClient()
	var enquiry Enquiry
	EnquiryCollection := client.Database("kbpl").Collection("enquiries")

	documentReturned := EnquiryCollection.FindOne(context.TODO(), bson.M{"_id": primitiveValueOfEnquiryId})
	err = documentReturned.Decode(&enquiry)
	if err != nil {
		log.Println(err)
	}

	selectedProductsObjectIds := make([]primitive.ObjectID, len(enquiry.ProductId))
	for i := range enquiry.ProductId {
		selectedProductsObjectIds[i], err = primitive.ObjectIDFromHex(enquiry.ProductId[i])
		log.Println(err)
	}

	var products []Product
	filter := bson.M{"active": true, "_id": bson.M{"$in": selectedProductsObjectIds}}
	ProductsCollection := client.Database("kbpl").Collection("products")
	cur, err := ProductsCollection.Find(context.TODO(), filter)
	if err != nil {
		log.Fatal("Error on Finding all the documents", err)
	}

	for cur.Next(context.TODO()) {
		var product Product
		err = cur.Decode(&product)
		if err != nil {
			log.Fatal("Error on Decoding the document", err)
		}
		products = append(products, product)
	}

	type Result struct {
		Products []Product
	}

	result := Result{
		products,
	}

	w.WriteHeader(http.StatusOK)
	err = tpl.ExecuteTemplate(w, "products", result)
	if err != nil {
		log.Fatal(err)
	}
}

func ManageHandler(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "adminManager", nil)
	if err != nil {
		log.Println(err)
	}
}

func GenerateJWT(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	// Check the credentials provided - if you store these in a database then
	// this is where your query would go to check.
	fmt.Println(r.FormValue("username"))
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	fmt.Println(username, password)
	if username != "myusername" || password != "mypassword" {
		w.WriteHeader(http.StatusUnauthorized)
		n, err := io.WriteString(w, `{"error":"invalid_credentials"}`)
		if err != nil {
			log.Println(n, err)
		}
		return
	}

	validTime := time.Now().Add(time.Hour * time.Duration(1)).Unix()

	// We are happy with the credentials, so build a token. We've given it
	// an expiry of 1 hour.
	token := jwt2.NewWithClaims(jwt2.SigningMethodHS256, jwt2.MapClaims{
		"user": username,
		"exp":  validTime,
		"iat":  time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(AppKey))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		n, err := io.WriteString(w, `{"error":"token_generation_failed"}`)
		log.Println(n, err)
		return
	}

	cookie := http.Cookie{
		Name:    "access_token",
		Value:   tokenString,
		Expires: time.Now().AddDate(0, 0, 1),
		Path:    "/",
	}
	http.SetCookie(w, &cookie)

	n, err := io.WriteString(w, `{"token":"`+tokenString+`"}`)
	if err != nil {
		log.Println(n, err)
	}
	return
}

func AdministratorHandler(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "adminHome", nil)
	if err != nil {
		log.Println("error parsing template adminHome", err)
	}
}

// Enquiries
func AdminEnquiriesHandler(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "adminEnquiries", nil)
	if err != nil {
		log.Println("error parsing template adminEnquiries", err)
	}
}

func AdminEnquiriesJsonHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fmt.Println(params)

	filter := bson.M{}

	// getting all enquiries
	var enquiries []Enquiry
	client := GetClient()
	collection := client.Database("kbpl").Collection("enquiries")

	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Println("error getting all enquiries in adminEnquiriesJsonHandler", err)
	}

	if err = cur.All(context.TODO(), &enquiries); err != nil {
		log.Fatal("error putting queries into &queries", err)
	}

	var dataField [][]string
	for _, value := range enquiries {
		dataField = append(dataField, []string{value.Name, value.Mobile, value.Email, value.Comments},
		)
	}

	datatableResponse := DataTableResponse{
		1,
		len(dataField),
		len(dataField),
		dataField,
	}

	encodedEnquiries, err := json.Marshal(datatableResponse)
	if err != nil {
		log.Println("error marshalling golang enquiries", err)
	}

	n, err := fmt.Fprintf(w, string(encodedEnquiries))
	if err != nil {
		log.Println("error in fmt.Fprintf encodedEnquiries", n, encodedEnquiries)
	}

}

func AdminViewEnquiryHandler(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "adminEnquiries", nil)
	if err != nil {
		log.Println("error parsing template adminViewEnquiry", err)
	}
}

func AdminSaveEnquiryHandler(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "adminEnquiries", nil)
	if err != nil {
		log.Println("error parsing template adminSaveEnquiry", err)
	}
}

// Products
func AdminProductsHandler(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "adminProducts", nil)
	if err != nil {
		log.Println("error parsing template adminProducts", err)
	}
}

func AdminProductsJsonHandler(w http.ResponseWriter, r *http.Request) {
	filter := bson.M{}

	// getting all enquiries
	var products []Product
	client := GetClient()
	collection := client.Database("kbpl").Collection("products")

	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Println("error getting all products in adminProductsJsonHandler", err)
	}

	if err = cur.All(context.TODO(), &products); err != nil {
		log.Fatal("error putting queries into &queries 2", err)
	}

	var dataField [][]string
	for _, value := range products {
		dataField = append(dataField, []string{value.Name, strconv.FormatInt(int64(value.Price), 10)})
	}

	datatableResponse := DataTableResponse{
		1,
		len(dataField),
		len(dataField),
		dataField,
	}

	encodedEnquiries, err := json.Marshal(datatableResponse)
	if err != nil {
		log.Println("error marshalling golang enquiries", err)
	}

	n, err := fmt.Fprintf(w, string(encodedEnquiries))
	if err != nil {
		log.Println("error in fmt.Fprintf encodedEnquiries", n, encodedEnquiries)
	}
}

func AdminViewProductHandler(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "adminProducts", nil)
	if err != nil {
		log.Println("error parsing template adminViewProduct", err)
	}
}

func AdminSaveProductHandler(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "adminProducts", nil)
	if err != nil {
		log.Println("error parsing template adminSaveProduct", err)
	}
}
