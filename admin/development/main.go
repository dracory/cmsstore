package main

// import (
// 	"database/sql"
// 	"errors"
// 	"log"
// 	"net/http"
// 	"time"

// 	"github.com/dracory/env"
// 	"github.com/dracory/hb"

// 	_ "github.com/go-sql-driver/mysql"
// 	_ "github.com/lib/pq"
// 	_ "modernc.org/sqlite"
// )

// var db *sql.DB

// func main() {
// 	log.Println("1. Initializing environment variables...")
// 	env.Load(".env")

// 	log.Println("2. Initializing database...")
// 	var err error
// 	db, err = mainDb(env.GetString("DB_DRIVER"), env.GetString("DB_HOST"), env.GetString("DB_PORT"), env.GetString("DB_DATABASE"), env.GetString("DB_USERNAME"), env.GetString("DB_PASSWORD"))

// 	if err != nil {
// 		log.Panic("Database is NIL: " + err.Error())
// 		return
// 	}

// 	if db == nil {
// 		log.Panic("Database is NIL")
// 		return
// 	}

// 	log.Println("4. Starting server on http://" + env.GetString("SERVER_HOST") + ":" + env.GetString("SERVER_PORT") + " ...")
// 	log.Println("URL: http://" + env.GetString("APP_URL") + " ...")
// 	mux := http.NewServeMux()
// 	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		menu := hb.NewHTML("<a href='/cms'>Standalone CMS</a> <br /> <br /> <a href='/cmswithlayout'>CMS WIth Layout</a> <br /> <br /> <a href='/embeddedcms'>Embedded CMS in IFRAME</a>")
// 		w.Write([]byte(menu.ToHTML()))
// 	})

// 	mux.HandleFunc("/cms", func(w http.ResponseWriter, r *http.Request) {
// 		//admin.New().Handle(w, r)
// 	})

// 	srv := &http.Server{
// 		Handler: mux,
// 		Addr:    env.GetString("SERVER_HOST") + ":" + env.GetString("SERVER_PORT"),
// 		// Good practice: enforce timeouts for servers you create!
// 		WriteTimeout:      15 * time.Second,
// 		ReadTimeout:       15 * time.Second,
// 		IdleTimeout:       30 * time.Second,
// 		ReadHeaderTimeout: 2 * time.Second,
// 	}

// 	log.Fatal(srv.ListenAndServe())
// }

// func mainDb(driverName string, dbHost string, dbPort string, dbName string, dbUser string, dbPass string) (*sql.DB, error) {
// 	var db *sql.DB
// 	var err error
// 	if driverName == "sqlite" {
// 		dsn := dbName
// 		db, err = sql.Open("sqlite", dsn)
// 	}
// 	if driverName == "mysql" {
// 		dsn := dbUser + ":" + dbPass + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"
// 		db, err = sql.Open("mysql", dsn)
// 	}
// 	if driverName == "postgres" {
// 		dsn := "host=" + dbHost + " user=" + dbUser + " password=" + dbPass + " dbname=" + dbName + " port=" + dbPort + " sslmode=disable TimeZone=Europe/London"
// 		db, err = sql.Open("postgres", dsn)
// 	}
// 	if err != nil {
// 		return nil, err
// 	}
// 	if db == nil {
// 		return nil, errors.New("database for driver " + driverName + " could not be initialized")
// 	}

// 	return db, nil
// }

// // func entityList() []cms.CustomEntityStructure {
// // 	list := []cms.CustomEntityStructure{}
// // 	list = append(list, cms.CustomEntityStructure{
// // 		Group:     "Users",
// // 		Type:      "user",
// // 		TypeLabel: "User",
// // 		AttributeList: []cms.CustomAttributeStructure{
// // 			{
// // 				Name:             "first_name",
// // 				Type:             "string",
// // 				FormControlLabel: "First Name",
// // 				FormControlType:  "input",
// // 				FormControlHelp:  "The first name of the user",
// // 			},
// // 			{
// // 				Name:             "last_name",
// // 				Type:             "string",
// // 				FormControlLabel: "Last Name",
// // 				FormControlType:  "input",
// // 				FormControlHelp:  "The last name of the user",
// // 			},
// // 			{
// // 				Name:             "email",
// // 				Type:             "string",
// // 				FormControlLabel: "E-mail",
// // 				FormControlType:  "input",
// // 				FormControlHelp:  "The e-mail address of the user",
// // 			},
// // 			{
// // 				Name:             "image_url",
// // 				Type:             "string",
// // 				FormControlLabel: "Image URL",
// // 				FormControlType:  "input",
// // 				FormControlHelp:  "The image of the product",
// // 			},
// // 		},
// // 	})
// // 	list = append(list, cms.CustomEntityStructure{
// // 		Group:     "Shop",
// // 		Type:      "shop_product",
// // 		TypeLabel: "Product",
// // 		AttributeList: []cms.CustomAttributeStructure{
// // 			{
// // 				Name:             "title",
// // 				Type:             "string",
// // 				FormControlLabel: "Title",
// // 				FormControlType:  "input",
// // 				FormControlHelp:  "The title which will be displayed to the customer",
// // 			},
// // 			{
// // 				Name:             "description",
// // 				Type:             "string",
// // 				FormControlLabel: "Description",
// // 				FormControlType:  "textarea",
// // 				FormControlHelp:  "The description which will be displayed to the customer",
// // 			},
// // 			{
// // 				Name:             "price",
// // 				Type:             "string",
// // 				FormControlLabel: "Price",
// // 				FormControlType:  "input",
// // 				FormControlHelp:  "The price of the product",
// // 			},
// // 			{
// // 				Name:             "image_url",
// // 				Type:             "string",
// // 				FormControlLabel: "Image URL",
// // 				FormControlType:  "input",
// // 				FormControlHelp:  "The image of the product",
// // 			},
// // 		},
// // 	})
// // 	list = append(list, cms.CustomEntityStructure{
// // 		Type:      "makeawish",
// // 		TypeLabel: "Make-a-Wish",
// // 		AttributeList: []cms.CustomAttributeStructure{
// // 			{
// // 				Name:             "wish",
// // 				Type:             "string",
// // 				FormControlLabel: "Wish",
// // 				FormControlType:  "textarea",
// // 				FormControlHelp:  "The wish that was made",
// // 			},
// // 			{
// // 				Name:             "referral",
// // 				Type:             "string",
// // 				FormControlLabel: "Referral",
// // 				FormControlType:  "input",
// // 				FormControlHelp:  "Where the wish was made from",
// // 			},
// // 		},
// // 	})
// // 	list = append(list, cms.CustomEntityStructure{
// // 		Group:     "Shop",
// // 		Type:      "shop_order",
// // 		TypeLabel: "Order",
// // 		AttributeList: []cms.CustomAttributeStructure{
// // 			{
// // 				Name:             "user_id",
// // 				Type:             "string",
// // 				FormControlLabel: "User ID",
// // 				FormControlType:  "input",
// // 				FormControlHelp:  "The ID of the user who made the purchase",
// // 				BelongsToType:    "user",
// // 			},
// // 			{
// // 				Name:             "total",
// // 				Type:             "string",
// // 				FormControlLabel: "Total",
// // 				FormControlType:  "input",
// // 				FormControlHelp:  "Total amount of the order",
// // 			},
// // 		},
// // 	})
// // 	list = append(list, cms.CustomEntityStructure{
// // 		Group:     "Shop",
// // 		Type:      "shop_order_line_item",
// // 		TypeLabel: "Order Line Item",
// // 		AttributeList: []cms.CustomAttributeStructure{
// // 			{
// // 				Name:             "order_id",
// // 				Type:             "string",
// // 				FormControlLabel: "Order ID",
// // 				FormControlType:  "input",
// // 				FormControlHelp:  "The order the item belongs to",
// // 				BelongsToType:    "shop_order",
// // 			},
// // 			{
// // 				Name:             "product_id",
// // 				Type:             "string",
// // 				FormControlLabel: "Product ID",
// // 				FormControlType:  "input",
// // 				FormControlHelp:  "The ID of the product that is ordered",
// // 			},
// // 			{
// // 				Name:             "quantity",
// // 				Type:             "string",
// // 				FormControlLabel: "Quantity",
// // 				FormControlType:  "input",
// // 				FormControlHelp:  "How many products are ordered (quantity) in this order item",
// // 			},
// // 			{
// // 				Name:             "subtotal",
// // 				Type:             "string",
// // 				FormControlLabel: "Subtotal",
// // 				FormControlType:  "input",
// // 				FormControlHelp:  "Subtotal amount of the order item",
// // 			},
// // 		},
// // 	})
// // 	return list
// // 	// {
// // 	// 	"type":"product",
// // 	// 	"price":"12.00",
// // 	// }
// // 	// list := []map[string]interface{}{
// // 	// 	{
// // 	// 		"type": "product",
// // 	// 		"attributes": []map[string]interface{}{
// // 	// 			{
// // 	// 				"name": "title",
// // 	// 				"type": "string",
// // 	// 				"rule": "required",
// // 	// 				// type - one of text, textarea, select, hidden, html
// // 	// 				// name - name of the input field as seen in the request
// // 	// 				// label - publicly visible name
// // 	// 				// width - width of the field - min 1, max 12
// // 	// 				// rule - rules for the field, used when validating
// // 	// 				// value - value of the field
// // 	// 				// options - array of options (used by the select type)
// // 	// 				// html - raw HTML to be displayed as-is (used by the html type)

// // 	// 			},
// // 	// 			{
// // 	// 				"name": "price",
// // 	// 				"type": "float",
// // 	// 				"rule": "required",
// // 	// 			},
// // 	// 			{
// // 	// 				"name": "image_url",
// // 	// 				"type": "string",
// // 	// 			},
// // 	// 		},
// // 	// 	},
// // 	// }
// // 	// return list
// // }
