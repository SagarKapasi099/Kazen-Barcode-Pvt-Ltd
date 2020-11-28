// static-files.go
package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"os"
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt2 "github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
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

var templatesPath = "templates/*.gohtml"
var tpl *template.Template

const AppKey = "klshifhjKLHLHGl;sdjhfl'kshjfgkSghsFHJKSHGSHGslhgh"
const statusClosed = "c"
const statusOpen = "o"

type Product struct {
	Id         string   `json:"_id" bson:"_id"`
	Name       string   `json:"name" bson:"name"`
	URL        string   `json:"url" bson:"url"`
	Properties []string `json:"properties" bson:"properties"`
	Price      int      `json:"price" bson:"price"`
	Active     bool     `json:"active" bson:"active"`
}

type Enquiry struct {
	Id          primitive.ObjectID `json:"id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Email       string             `json:"email" bson:"email"`
	Mobile      string             `json:"mobile" bson:"mobile"`
	Comments    string             `json:"comments" bson:"comments"`
	OTP         string             `json:"-" bson:"otp"`
	ProductId   []string           `json:"-" bson:"product_id"`
	Status      string             `json:"status" bson:"status"`
	CreatedDate time.Time          `json:"created_date" bson:"created_date"`
}

type DataTableResponse struct {
	Draw            int        `json:"draw"`
	RecordsTotal    int        `json:"recordsTotal"`
	RecordsFiltered int        `json:"recordsFiltered"`
	Data            [][]string `json:"data"`
}

type DatatableView struct {
	Add      template.URL
	Create   template.URL
	ReadJson template.URL
	Read     template.URL
	Update   template.URL
	MarkDone template.URL
}

type TableColumns struct {
	Columns []string
}

/**
Dynamic Header Dropdown
*/
type ProductType struct {
	Text []string
}

type Headers struct {
	ProductType ProductType
}

func main() {

	var err error
	tpl = template.Must(template.New("").Funcs(template.FuncMap{"args": args, "headers": getHeader}).ParseGlob(templatesPath))

	// checking to see if database connection is alive
	c := GetClient()
	pingError := c.Ping(context.Background(), readpref.Primary())
	if pingError != nil {
		log.Println("Couldn't connect to the database", err)
	} else {
		log.Println("Connected!")
	}

	// open
	mySQLClient, err := getMySQLClient("Main Calling")
	// ping
	if err != nil {
		log.Fatal("cannot connect to MySQL in main()", err)
	}
	// close
	if err := mySQLClient.Close(); err != nil {
		log.Println("MySQL Cannot db.Close On Defer main()", err)
	} else {
		log.Println("close called main()", err)
	}
	fmt.Printf("start MySQL Stats: %+v\n", mySQLClient.Stats())
	log.Println("MySQL Connection: ", mySQLClient.Ping())

	r := mux.NewRouter()
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	r.HandleFunc("/", HomeHandler).Methods("GET")
	r.HandleFunc("/index.html.var", HomeHandler).Methods("GET") // apache reverse proxy calls this url for index
	r.HandleFunc("/about", AboutHandler).Methods("GET")
	r.HandleFunc("/contact", ContactUsHandler).Methods("GET")
	r.HandleFunc("/enquiry", EnquiryHandler).Methods("POST")
	r.HandleFunc("/verifyOTP", VerifyOTPHandler).Methods("POST")
	r.HandleFunc("/products", ShowProductsHandler).Methods("POST")
	r.HandleFunc("/products/{type}", ShowProductTypesHandler).Methods("GET")

	// Admin Section
	r.HandleFunc("/manage", ManageHandler).Methods("GET")
	r.HandleFunc("/getToken", GenerateJWT)
	r.Handle("/administrator", AuthMiddleware(http.HandlerFunc(AdministratorHandler))).Methods("GET")
	// Enquiries
	r.Handle("/administrator/enquiries", AuthMiddleware(http.HandlerFunc(AdminEnquiriesHandler))).Methods("GET")
	r.Handle("/administrator/getEnquiriesJson", AuthMiddleware(http.HandlerFunc(AdminEnquiriesJsonHandler))).Methods("POST")
	r.Handle("/administrator/viewEnquiry/{id}", AuthMiddleware(http.HandlerFunc(AdminViewEnquiryHandler))).Methods("GET")
	r.Handle("/administrator/updateEnquiry", AuthMiddleware(http.HandlerFunc(AdminUpdateEnquiryHandler))).Methods("POST")
	// Products
	r.Handle("/administrator/products", AuthMiddleware(http.HandlerFunc(AdminProductsHandler))).Methods("GET")
	r.Handle("/administrator/getProductsJson", AuthMiddleware(http.HandlerFunc(AdminProductsJsonHandler))).Methods("POST")
	r.Handle("/administrator/viewProduct/{id}", AuthMiddleware(http.HandlerFunc(AdminViewProductHandler))).Methods("GET")
	r.Handle("/administrator/saveProduct", AuthMiddleware(http.HandlerFunc(AdminSaveProductHandler))).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// used for passing data to child templates
func args(vs ...interface{}) []interface{} { return vs }

// GetClient returns a MongoDB Client Singleton
func GetClient() *mongo.Client {
	clientOptions := options.Client().ApplyURI("mongodb+srv://SagarKapasi099:3wqzTsSvNQkovuxi@projectautodidact-5vr5f.gcp.mongodb.net/test?retryWrites=true&w=majority")
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Println("error creating mongo.NewClient()", err)
	}
	err = client.Connect(context.Background())
	if err != nil {
		log.Println(err)
	}
	return client
}

// GetMySQLClient returns a MySQL Client Singleton
func getMySQLClient(logMessage string) (*sql.DB, error) {
	fmt.Printf("G E N E R A T E D %v\n", logMessage)
	db, err := sql.Open("mysql", os.Getenv("KBPL_DATABASE_USER")+":"+os.Getenv("KBPL_DATABASE_PASSWORD")+"@tcp(127.0.0.1:3306)/kbpl")
	if err != nil {
		log.Println("error opening MySQL on port 3306")
		return nil, err
	}

	if db == nil {
		err = errors.New("db is nil in getMySQLClient")
		return nil, err
	}
	db.SetConnMaxLifetime(time.Minute * 4)
	db.SetMaxIdleConns(50)
	db.SetMaxOpenConns(70)
	/*
		if err := db.PingContext(context.TODO()); err != nil {
			log.Fatal(err)
		}
	*/
	return db, nil
}

func closeDatabase(db *sql.DB, logMessage string) {
	if err := db.Close(); err != nil {
		fmt.Println("close err ", logMessage, err)
	} else {
		fmt.Println("close success ", logMessage, err)
	}
	fmt.Printf("ping?%+v\n", db.Ping())
	fmt.Printf("stats: %+v\n", db.Stats())
}

func GenerateOTP(length int) string {
	const charset = "0123456789"
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
			if err != nil || accessTokenCookie.Value == "" {
				return "", errors.New("Error Obtaining Cookie access_token")
			}
			jwtCookie := accessTokenCookie.Value
			return jwtCookie, nil
		},
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, errMsg string) {
			if r.Method == "POST" {
				http.Error(w, "", http.StatusUnauthorized)
			} else {
				http.Redirect(w, r, "/manage", http.StatusSeeOther)
			}
		},
		ValidationKeyGetter: func(token *jwt2.Token) (interface{}, error) {
			return []byte(AppKey), nil
		},
		SigningMethod: jwt2.SigningMethodHS256,
	})
	return jwtMiddleware.Handler(next)
}

// get dynamic header properties for navigation bar
func getHeader() Headers {
	var headers Headers
	db, err := getMySQLClient("HeaderIncoming")
	if err != nil {
		log.Println("error getting db connection in header", err)
	}
	rows, err := db.Query("select pt_text from product_type")
	if err != nil {
		log.Println("error getting product type in getHeader()", err)
	}
	defer rows.Close()
	for rows.Next() {
		var currentType string
		err = rows.Scan(&currentType)
		if err != nil {
			log.Println("get header rows.Scan(&currentType) error: ", err)
		}
		headers.ProductType.Text = append(headers.ProductType.Text, currentType)
	}
	err = rows.Err()
	if err != nil {
		log.Println(err)
	}

	closeDatabase(db, "HeaderIncoming")
	return headers
}

func getDataProducts(productTypeId int, enquiryId int64) ([]Product, error) {
	var products []Product

	db, err := getMySQLClient("getDataProductsIncoming")
	if err != nil {
		log.Println("error getting db connection in getDataProducts", err)
	}

	rows, err := db.Query("SELECT p_id, p_name, p_url, p_price, p_active from product where p_active = ?", 1)
	if err != nil {
		log.Println("default query inside getDataProducts (from product where p_active = 1) returned error: ", err)
	}
	if productTypeId != 0 {
		err = rows.Close()
		if err != nil {
			log.Println("where p_active = 1 unable to close row", err)
		}
		rows, err = db.Query("SELECT p_id, p_name, p_url, p_price, p_active from product where p_active = ? and pt_id = ?", 1, productTypeId)
	}
	if enquiryId != 0 {
		err = rows.Close()
		if err != nil {
			log.Println("where p_active = 1 unable to close row", err)
		}
		rows, err = db.Query("SELECT product.p_id, product.p_name, product.p_url, product.p_price, product.p_active " +
			"from product " +
			"left join enquiry_product_mapping as epm on epm.p_id = product.p_id " +
			"left join enquiry on enquiry.e_id = epm.e_id " +
			"where product.p_active = ? and epm.e_id = ?", 1, enquiryId)
	}
	if err != nil {
		log.Println("Error on Finding all the documents in home handler", err)
	}
	for rows.Next() {
		var product Product
		err = rows.Scan(&product.Id, &product.Name, &product.URL, &product.Price, &product.Active)
		if err != nil {
			log.Println("Error on Decoding the product in home handler", err)
		}
		// adding product properties[]
		propertiesRow, err := db.Query("SELECT prop_text from product_properties where p_id = ?", product.Id)
		if err != nil {
			log.Println("Error on getting prop_text from product properties where p_id is ", product.Id, err)
		}
		for propertiesRow.Next() {
			var currentProperty string
			err = propertiesRow.Scan(&currentProperty)
			product.Properties = append(product.Properties, currentProperty)
		}
		products = append(products, product)
	}
	if err := rows.Close; err != nil {
		fmt.Print("E R R O R Rows")
		fmt.Printf("%v %p\n", &err, err)
	}

	closeDatabase(db, "getDataProductsIncoming")
	//fmt.Printf("getDataProducts Stats: %+v\n", db.Stats())
	////if err := db.Close(); err != nil {
	////	log.Println("MySQL Cannot db.Close On Defer inside getDataProducts", err)
	////} else {
	////	log.Println("close called getDataProducts", err)
	////}
	//fmt.Printf("close?%+v\n", db.Ping())
	//fmt.Printf("getDataProducts Stats: %+v\n", db.Stats())

	return products, err
}

func AboutHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	err := tpl.ExecuteTemplate(w, "about", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	products, err := getDataProducts(0, 0)

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

	success := true

	var productIds []string
	if err := json.Unmarshal([]byte(productIdsJson), &productIds); err != nil {
		success = false
		log.Println(err)
	}

	otp := GenerateOTP(6)

	db, err := getMySQLClient("enquiryHandler")

	if err != nil {
		success = false
		log.Println("error getting db connection in enquiryHandler", err)
	}

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := db.Prepare("INSERT INTO enquiry set e_name=?, e_email=?, e_mobile=?, e_comments=?, e_otp=?, e_status=?")
	if err != nil {
		success = false
		log.Println("error preparing insert enquiry statement in enquiryHandler", err)
	}

	res, err := stmt.Exec(name, email, number, comments, otp, "o")
	if err != nil {
		success = false
		log.Println("error stmt.Exec insert in enquiryHandler", err)
	}

	enquiryId, err := res.LastInsertId()
	if err != nil {
		success = false
		log.Println("error getting last insert id in enquiryHandler")
	}

	err = stmt.Close()
	if err != nil {
		log.Println("cannot close stmt enquiryHandler 1234", err)
	}

	// insert into enquiry product mapping
	sqlStr := "INSERT INTO enquiry_product_mapping(e_id, p_id) VALUES "
	vals := []interface{}{}

	for _, row := range productIds {
		sqlStr += "(?, ?),"
		vals = append(vals, enquiryId, row)
		log.Println(vals)
	}
	//trim the last ,
	sqlStr = sqlStr[0 : len(sqlStr)-1]
	log.Println(sqlStr)
	//prepare the statement
	stmtEPM, err := db.Prepare(sqlStr)
	if err != nil {
		success = false
		log.Println("Error preparing sql statement for stmtEPM (enquiry_product_mapping)", err)
	}

	//format all vals at once
	resEPM, err := stmtEPM.Exec(vals...)
	if err != nil {
		success = false
		log.Println("Error executing sql statement resEMP for stmtEPM (enquiry_product_mapping) resEPM's Value:", resEPM, err)
	}
	err = stmtEPM.Close()
	if err != nil {
		log.Println("cannot close stmtEPM enquiryHandler 12345", err)
	}

	if success {
		err = tx.Commit()
		if err != nil {
			log.Println("error committing ", err)
		}
	} else {
		err = tx.Rollback()
		log.Println("error enquiryHandler transaction rollback", err)
	}

	err = db.Close()
	if err != nil {
		log.Println("cannot close db enquiryHandler 123456", err)
	}

	jsonResponse := `{
		"success": ` + strconv.FormatBool(success) + `,
		"message":{
			"superText":"Please Enter OTP Sent On ` + number + `",
			"subText": "",
			"buttonText": "Verify"
		},
		"number": "` + number + `",
		"enquiryId": "` + strconv.FormatInt(enquiryId, 10) + `",
		"method": "post"
	}`

	n, err := fmt.Fprintf(w, jsonResponse)
	if err != nil {
		log.Println(n, err)
	}
}

func VerifyOTPHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: avoid brute force attack

	enquiryId := r.FormValue("enquiryId")

	db, err := getMySQLClient("VerifyOTPHandler")
	if err != nil {
		log.Println("error getting database instance for VerifyOTPHandler. enquiryId:", enquiryId, err)
	}

	rows := db.QueryRow("SELECT e_otp FROM enquiry where e_id=?", enquiryId)
	var otp string
	err = rows.Scan(&otp)
	if err != nil {
		log.Println("error getting otp for enqiuryId:", enquiryId, err)
	}

	showProducts := "false" // would be updated to true by the below query if more than one product is present
	err = db.QueryRow("SELECT IF(COUNT(*),'true','false') FROM enquiry_product_mapping WHERE e_id = ?", enquiryId).Scan(&showProducts)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Close()
	if err != nil {
		log.Println("cannot close stmt VerifyOTPHandler 1234", err)
	}

	otpReceived := r.FormValue("otp")

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
	enquiryIdInt64, err := strconv.ParseInt(enquiryId, 10, 64)
	if err != nil {
		log.Println("error strconv.ParseInt(enquiryId, 10, 64) in show products handler", err)
	}

	products, err := getDataProducts(0, enquiryIdInt64)

	type Result struct {
		Products []Product
		Title    string
	}

	result := Result{
		products,
		"Enquired Products With Prices",
	}

	w.WriteHeader(http.StatusOK)
	err = tpl.ExecuteTemplate(w, "products", result)
	if err != nil {
		log.Println(err)
	}
}

func ShowProductTypesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productType := vars["type"]

	db, err := getMySQLClient("showProductTypesHandlerIncoming")
	if err != nil {
		log.Println("error getting db connection in showProductTypesHandler", err)
	}

	rows, err := db.Query("select pt_id, pt_text from product_type where pt_text = ? limit 1", productType)
	if err != nil {
		// TODO: redirect to homepage OR 404
		log.Println("cannot find product_type for product_type pt_text: ", productType, err)
	}
	defer rows.Close()

	var productTypeId int
	var productTypeName string
	for rows.Next() {
		err = rows.Scan(&productTypeId, &productTypeName)
		if err != nil {
			// TODO: redirect to homepage OR 404
			log.Println("cannot rows.Scan product_type pt_text: ", productType, err)
		}
		// break // we only need a single match (limit 1 also used in query)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	// we now have the pt_id, now we fetch and display products of that type (pt_id)
	products, err := getDataProducts(productTypeId, 0)

	closeDatabase(db, "showProductTypesHandlerIncoming")
	//if err := db.Close(); err != nil {
	//	log.Println("MySQL Cannot db.Close On Defer inside showProductTypesHandler()", err)
	//} else {
	//	log.Println("i swear i closed the database connection", db.Ping())
	//}

	fmt.Printf("the end: %+v\n", db.Stats())

	type Result struct {
		Products []Product
		Title    string
	}

	result := Result{
		products,
		productTypeName,
	}

	w.WriteHeader(http.StatusOK)
	err = tpl.ExecuteTemplate(w, "products", result)
	if err != nil {
		log.Println(err)
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
	datatableViewData := DatatableView{
		template.URL(""),
		template.URL(""),
		template.URL("/administrator/getEnquiriesJson"),
		template.URL("/administrator/viewEnquiry"),
		template.URL("/administrator/updateEnquiry"),
		template.URL("/administrator/markDoneEnquiry"),
	}

	type Result struct {
		DatatableView
		TableColumns
	}

	data := Result{
		datatableViewData,
		TableColumns{
			[]string{
				"Name",
				"Mobile No",
				"Email",
				"Comments",
				"Created On",
				"Actions",
			},
		},
	}

	err := tpl.ExecuteTemplate(w, "adminDatatableView", data)
	if err != nil {
		log.Println("error parsing template adminEnquiries", err)
	}
}

func AdminEnquiriesJsonHandler(w http.ResponseWriter, r *http.Request) {
	// getting all enquiries
	var enquiries []Enquiry
	client := GetClient()
	collection := client.Database("kbpl").Collection("enquiries")

	start, err := strconv.Atoi(r.FormValue("start"))
	if err != nil {
		// TODO replace with json response
		log.Println("wrong start value for adminEnquiriesJsonHandler", err)
	}

	length, err := strconv.Atoi(r.FormValue("length"))
	if err != nil {
		// TODO replace with json response
		log.Println("wrong length value for adminEnquiriesJsonHandler", err)
	}

	draw, err := strconv.Atoi(r.FormValue("draw"))
	if err != nil {
		// TODO replace with json response
		log.Println("wrong draw value for adminEnquiriesJsonHandler", err)
	}

	filter := bson.M{}
	filterOptions := options.Find()
	filterOptions.SetSkip(int64(start))
	filterOptions.SetLimit(int64(length))
	filterOptions.SetSort(bson.M{"created_date": -1})
	cur, err := collection.Find(context.TODO(), filter, filterOptions)
	if err != nil {
		log.Println("error getting all enquiries in adminEnquiriesJsonHandler", err)
	}
	curCount, err := collection.CountDocuments(context.TODO(), filter, options.Count())
	if err != nil {
		log.Println("error getting count for all enquiries in adminEnquiriesJsonHandler", err)
	}

	if err = cur.All(context.TODO(), &enquiries); err != nil {
		log.Println("error putting queries into &queries", err)
	}

	var dataField [][]string
	for _, value := range enquiries {
		dataField = append(dataField, []string{value.Name, value.Mobile, value.Email, value.Comments, value.CreatedDate.Format("2006-01-02 15:04:05"), value.Status, value.Id.Hex()})
	}

	datatableResponse := DataTableResponse{
		draw,
		int(curCount),
		int(curCount),
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
	params := mux.Vars(r)
	id := params["id"]
	if id == "" || len(id) == 0 {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	client := GetClient()
	collection := client.Database("kbpl").Collection("enquiries")

	primitiveObjectFromHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("objectIdFromHexFailed", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}

	filter := bson.M{"_id": primitiveObjectFromHex}
	filterOptions := options.FindOne()

	cur := collection.FindOne(context.TODO(), filter, filterOptions)

	var enquiry Enquiry

	if err := cur.Decode(&enquiry); err != nil {
		log.Println("error decoding enquiry", err)
		http.Error(w, "Bad Request (Code: ERRONEDECODE)", http.StatusBadRequest)
		return
	}

	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		log.Println("Getting Timezone Asia/Kolkata is causing error in adminViewEnquiryHandler: ", err)
	}
	enquiry.CreatedDate = enquiry.CreatedDate.In(loc)

	type Result struct {
		Enquiry
	}

	result := Result{
		enquiry,
	}

	err = tpl.ExecuteTemplate(w, "adminSingleView", result)
	if err != nil {
		log.Println("error parsing template adminViewEnquiry", err)
	}
}

/**
set the status to closed
*/
func AdminUpdateEnquiryHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	status := r.FormValue("status")
	if id == "" || (status != statusOpen && status != statusClosed) {
		http.Error(w, "No Data Received", http.StatusBadRequest)
		return
	}

	var statusUpdate string
	if status == statusClosed {
		statusUpdate = statusClosed
	} else if status == statusOpen {
		statusUpdate = statusOpen
	} else {
		http.Error(w, "Valid Status Not Found", http.StatusBadRequest)
		return
	}

	client := GetClient()
	collection := client.Database("kbpl").Collection("enquiries")

	opts := options.FindOneAndUpdate().SetUpsert(true)
	objectIdFromHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("cannot convert id into objectId in updateEnquiryHandler", err)
	}
	filter := bson.D{{"_id", objectIdFromHex}}
	update := bson.D{{"$set", bson.D{{"status", statusUpdate}}}}
	var updatedDocument bson.M
	err = collection.FindOneAndUpdate(context.TODO(), filter, update, opts).Decode(&updatedDocument)
	if err != nil {
		log.Println("error updating enquiry", err)
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			log.Println("Error, there are no documents", err)
			return
		}
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
