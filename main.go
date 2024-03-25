package main

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	productmodel "goserver/models/product"
	usermodel "goserver/models/user"
	"html/template"
	"regexp"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

// "net/url"
func main() {

	// Create a new engine
	engine := html.New("./views", ".html")
	engine.AddFunc(
		// add unescape function
		"unescape", func(s string) template.HTML {
			return template.HTML(s)
		},
	)

	// Pass the engine to the Views
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		var p = home()

		//	fmt.Printf("s%", p)

		// Render index within layouts/main
		return c.Render("index", fiber.Map{
			"Title":    "Hello, World!",
			"Products": p}, "layouts/main")
		//return c.SendString()
	})
	app.Get("/order/:id", func(c *fiber.Ctx) error {
		fmt.Fprintf(c, "%s\n", c.Params("id"))
		var p = order(c.Params("id"))

		//	fmt.Printf("s%", p)

		// Render order within layouts/main
		return c.Render("order", fiber.Map{
			"Title":   "Hello, World!",
			"Product": p}, "layouts/main")
	})
	app.Use("/papi", func(c *fiber.Ctx) error {
		auth := c.Get(fiber.HeaderAuthorization)
		var ba = ""
		if raw, err := base64.StdEncoding.DecodeString(auth[6:]); err == nil {
			ba = string(raw)
		}
		s := strings.Split(ba, ":")
		user, pass := s[0], s[1]
		//	fmt.Println("user:", user)
		//	fmt.Println("pass:", pass)

		var userObj = usermodel.CheckUser(user, pass)

		response := map[string]interface{}{"success": false}
		var r map[string]interface{}

		var tobj TestJsonObject

		err := c.BodyParser(&tobj)
		if err != nil {
			fmt.Printf("Err: s%", err)
		}

		if tobj.Method == "register" && len(userObj) == 0 && len(pass) == 0 {
			var obj RegisterObject
			err = c.BodyParser(&obj)
			r = register(obj)
		} else if tobj.Method == "checkIn" {
			fmt.Println("Checking in")
			usermodel.CheckIn(pass)
			return c.Send([]byte{})
		} else if tobj.Method == "submitProduct" && len(userObj) != 0 {
			var obj SubmitProductObject
			err = c.BodyParser(&obj)
			r = submitProduct(userObj, obj)
		} else if tobj.Method == "submitIAddress" && len(userObj) != 0 {
			var obj SubmitIAddressObject
			err = c.BodyParser(&obj)
			r = submitIAddress(userObj, obj)
		}

		if err != nil {
			fmt.Println(err)
		}

		if r["success"] == true && tobj.Method == "register" {
			response["success"] = r["success"]
			response["reg"] = r["reg"]
		} else if r["success"] == true {
			response["success"] = r["success"]
		}

		return c.JSON(response)
	})
	app.Listen(":3000")
}

type TestJsonObject struct {
	Method string `json:"method"`
}

type RegisterObject struct {
	Method string `json:"method"`
	Params struct {
		Username string
		Wallet   string
	}
}

type SubmitProductObject struct {
	Method string `json:"method"`
	Params struct {
		Id        int
		P_type    string `json:"p_type,omitempty"`
		Label     string `json:"label,omitempty"`
		Details   string `json:"details,omitempty"`
		Scid      string `json:"scid,omitempty"`
		Inventory int    `json:"inventory,omitempty"`
		Image     string `json:"image,omitempty"`
		Action    string `json:"action,omitempty"`
	}
}
type SubmitIAddressObject struct {
	Method string `json:"method"`
	Params struct {
		Id           int
		Product_id   int    `json:"product_id,omitempty"`
		Iaddr        string `json:"iaddr,omitempty"`
		Ask_amount   int    `json:"ask_amount,omitempty"`
		Comment      string `json:"comment,omitempty"`
		Ia_scid      string `json:"ia_scid,omitempty"`
		Status       int    `json:"status,omitempty"`
		Ia_inventory int    `json:"ia_inventory,omitempty"`
		Action       string `json:"action,omitempty"`
	}
}

func register(obj RegisterObject) map[string]interface{} {

	//	fmt.Println("len(userObj):", len(userObj))
	//	fmt.Println("obj.Method:", obj.Method)

	var reg = usermodel.InsertUser(obj.Params.Username, obj.Params.Wallet)
	if reg != "failed" {
		return map[string]interface{}{"success": true, "reg": reg}
	}
	return map[string]interface{}{"success": false}
}

func submitProduct(userObj map[string]interface{}, obj SubmitProductObject) map[string]interface{} {
	response := map[string]interface{}{"success": true}
	if productmodel.SubmitProduct(userObj, obj.Params) {
		return response
	}
	response = map[string]interface{}{"success": false}
	return response
}
func submitIAddress(userObj map[string]interface{}, obj SubmitIAddressObject) map[string]interface{} {
	response := map[string]interface{}{"success": true}
	if productmodel.SubmitIAddress(userObj, obj.Params) {
		return response
	}
	response = map[string]interface{}{"success": false}
	return response
}

/**/
func home() []interface{} {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/pong_store")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	now := time.Now().UTC()
	count := 30
	then := now.Add(time.Duration(-count) * time.Minute)

	results, err := db.Query("SELECT pid,label,details,inventory,user,image,username FROM products INNER JOIN users ON products.user = users.userid WHERE users.checkin > ? ORDER BY products.id DESC", then.Format("2006-01-02 15:04:05"))
	if err != nil {
		panic(err.Error())
	}

	var products []interface{}

	for results.Next() {
		var (
			pid       int
			label     string
			details   string
			inventory int
			user      string
			image     string
			username  string
		)
		var product Product
		if err := results.Scan(&pid, &label, &details, &inventory, &user, &image, &username); err != nil {
			panic(err)
		}

		product.Label = label
		if len(label) > 50 {
			product.Label = label[0:50]
		}

		product.Details = details
		if len(details) > 100 {
			product.Details = details[0:50] + "..."
		}

		product.Inventory = inventory
		product.User = user

		product.Img = ""
		if len(image) > 22 {
			teststr := image
			if teststr[0:22] == "data:image/png;base64," {
				product.Img = "<img class='product_image' src='" + image + "'>"
			}
		}

		if image == "" {
			product.Image = false
		} else {
			product.Image = true
		}

		product.Username = username
		if username == "" {
			product.Username = "Username Not Provided"
		}

		use_ia_inventory := false
		if product.Inventory == 0 {
			use_ia_inventory = true
		}

		product.Use_p_inv = true
		if use_ia_inventory {
			product.Use_p_inv = false
		}

		results2, err := db.Query("SELECT id,comment,ia_inventory,ask_amount,status FROM i_addresses WHERE i_addresses.product_id = ? AND i_addresses.user = ?", pid, user)
		if err != nil {
			panic(err.Error())
		}

		for results2.Next() {
			var (
				id           int
				comment      string
				ia_inventory int
				ask_amount   int
				status       int
			)
			var iaddress IAddress
			if err := results2.Scan(&id, &comment, &ia_inventory, &ask_amount, &status); err != nil {
				panic(err)
			}

			iaddress.Id = id
			iaddress.Comment = comment
			iaddress.Ia_inventory = ia_inventory

			iaddress.Ask_amount = float64(ask_amount) * .00001
			iaddress.Status = status

			skip_ia := false
			if (use_ia_inventory && iaddress.Ia_inventory < 1) || iaddress.Status == 0 {
				skip_ia = true
			}

			if !use_ia_inventory {
				iaddress.Ia_inventory = 0
			}

			if !skip_ia {
				product.Iaddresses = append(product.Iaddresses, iaddress)
			}
		}

		if len(product.Iaddresses) > 0 {
			products = append(products, product)
		}
	}

	return products

}

type Product struct {
	Label      string
	Details    string
	Inventory  int
	User       string
	Image      bool
	Img        string
	Username   string
	Use_p_inv  bool
	Iaddresses []interface{}
}

type IAddress struct {
	Id           int
	Comment      string
	Ia_inventory int
	Ask_amount   float64
	Status       int
}

func order(iaid string) FullProduct {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/pong_store")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	rows, err := db.Query(`SELECT *,i_addresses.id ,i_addresses.status AS ia_status FROM i_addresses 
	RIGHT JOIN products ON i_addresses.product_id = products.pid
	RIGHT JOIN users ON i_addresses.user = users.userid
	WHERE i_addresses.id = ? AND i_addresses.user = products.user`, iaid)
	if err != nil {
		fmt.Printf("s%", err)
	}

	cols, _ := rows.Columns()
	//fmt.Println("columns", cols)
	var product FullProduct

	data := make(map[string]string)

	if rows.Next() {
		columns := make([]string, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i, _ := range columns {
			columnPointers[i] = &columns[i]
		}

		rows.Scan(columnPointers...)

		for i, colName := range cols {
			data[colName] = columns[i]
		}
	}

	product.Label = data["label"]

	product.Inventory = strToInt(data["inventory"])
	product.Ia_inventory = strToInt(data["ia_inventory"])
	product.Ia_status = strToInt(data["ia_status"])
	product.Ia_comment = data["ia_comment"]
	product.IAddress = data["iaddr"]

	if data["username"] == "" {
		product.Username = data["wallet"]
	} else {
		product.Username = data["username"]
	}
	product.Img = ""
	if len(data["image"]) > 22 {
		teststr := data["image"]
		if teststr[0:22] == "data:image/png;base64," {
			product.Img = "<img style='max-width: 100%;' src='" + data["image"] + "'>"
		}
	}

	if data["image"] == "" {
		product.Image = false
	} else {
		product.Image = true
	}

	ask_amount, _ := strconv.ParseFloat(data["ask_amount"], 5)
	product.Ask_amount = ask_amount * .00001

	//See if is a smart contract.
	scid := data["scid"]
	if data["ia_scid"] != "" {
		scid = data["ia_scid"]
	}
	product.Scid = scid

	r := regexp.MustCompile(`<.*?>`)
	product.Details = r.ReplaceAllString(data["details"], "")
	product.Details = strings.Replace(product.Details, "\n", "<br>", -1)
	product.P_type = data["p_type"]

	product.Is_physical = false
	if data["p_type"] == "physical" {
		product.Is_physical = true
	}

	stock := ""
	product_inventory := false
	if product.Ia_status != 0 {
		if product.Inventory > 0 {
			product_inventory = true
			stock = "Available:" + strconv.Itoa(product.Inventory)

		} else if product.Ia_inventory > 0 {
			stock = "Available:" + strconv.Itoa(product.Ia_inventory)

		} else {
			stock = "Out of Stock"
		}
	} else {
		stock = "Item Currently Unavailable."
	}
	product.Stock = stock
	product.Product_inv = product_inventory

	hidden := ""
	if product.P_type == "physical" || stock == "Out of Stock" || stock == "Item Currently Unavailable." {
		hidden = "hidden"
	}

	product.IAClass = hidden

	results2, err := db.Query("SELECT id,comment,ia_inventory,ask_amount,i_addresses.status AS ia_status FROM i_addresses WHERE i_addresses.product_id = ? AND i_addresses.user = ?", data["product_id"], data["user"])
	if err != nil {
		panic(err.Error())
	}

	//var iaddresses []interface{}

	for results2.Next() {
		var (
			id           int
			comment      string
			ia_inventory int
			ask_amount   int
			ia_status    int
		)
		var iaddress FullIAddress
		if err := results2.Scan(&id, &comment, &ia_inventory, &ask_amount, &ia_status); err != nil {
			panic(err)
		}

		iaddress.Id = id
		iaddress.Comment = comment
		iaddress.Ia_inventory = ia_inventory
		iaddress.Ask_amount = ask_amount
		iaddress.Ia_status = ia_status

		class := "greyed_out"
		if product_inventory {
			class = ""
		}
		if iaddress.Ia_inventory > 0 {
			class = ""
		}
		if iaddress.Ia_status == 0 {
			class = "greyed_out"
		}

		selected := false
		if iaddress.Id == strToInt(iaid) {
			selected = true
		}
		iaddress.Selected = selected
		iaddress.Class = class

		if iaddress.Ia_status != 0 || iaddress.Id == strToInt(iaid) {
			product.Iaddresses = append(product.Iaddresses, iaddress)
		}
	}

	return product
}

type FullProduct struct {
	Label        string
	Inventory    int
	Ia_inventory int
	Ia_status    int
	Ia_comment   string
	IAddress     string
	Username     string
	Image        bool
	Img          string
	Ask_amount   float64
	Scid         string
	Details      string
	P_type       string
	Is_physical  bool
	Stock        string
	Product_inv  bool
	Selected     bool
	IAClass      string
	Iaddresses   []interface{}
}
type FullIAddress struct {
	Id           int
	Comment      string
	Ia_inventory int
	Ask_amount   int
	Ia_status    int
	Selected     bool
	Class        string
}

func strToInt(str string) int {
	int, _ := strconv.Atoi(str)
	return int
}
