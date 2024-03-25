package main

import (
	"encoding/base64"
	"fmt"
	usermodel "goserver/models/user"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func handlePapi(c *fiber.Ctx) error {
	// Extract username and password from Basic Authentication header
	auth := c.Get(fiber.HeaderAuthorization)
	var ba = ""
	if raw, err := base64.StdEncoding.DecodeString(auth[6:]); err == nil {
		ba = string(raw)
	}
	s := strings.Split(ba, ":")
	user, pass := s[0], s[1]

	// fmt.Println("user:", user)
	// fmt.Println("pass:", pass)

	// Check user authentication
	userObj := usermodel.CheckUser(user, pass)

	response := map[string]interface{}{"success": false}
	var r map[string]interface{}

	var testobject TestJsonObject

	err := c.BodyParser(&testobject)
	if err != nil {
		fmt.Printf("Err: %s", err)
	}
	// fmt.Print(userObj)

	if len(userObj) != 0 {
		switch testobject.Method {
		case "checkIn":
			fmt.Println("Checking in")
			usermodel.CheckIn(pass)
			return c.Send([]byte{})
		case "submitProduct":
			var obj SubmitProductObject
			err = c.BodyParser(&obj)
			r = submitProduct(userObj, obj)
		case "submitIAddress":
			var obj SubmitIAddressObject
			err = c.BodyParser(&obj)
			r = submitIAddress(userObj, obj)
		case "newTX":
			fmt.Println("submitTX")
			var obj SubmitTxObject
			err = c.BodyParser(&obj)
			r = submitTx(userObj, obj)
		default:
			return c.JSON(response)
		}
	} else if testobject.Method == "register" && len(pass) == 0 {
		var obj RegisterObject
		err = c.BodyParser(&obj)
		r = register(obj)
	}

	if err != nil {
		fmt.Println(err)
	}

	// Update response based on method and success
	if r != nil && r["success"] == true {
		response["success"] = r["success"]
		if testobject.Method == "register" {
			response["reg"] = r["reg"]
		}
	}

	// Return JSON response
	return c.JSON(response)
}

// Define the route handler function
func handleOrder(c *fiber.Ctx) error {
	// Get the order ID from the URL parameters
	orderID := c.Params("id")

	// Print the order ID (optional)
	fmt.Fprintf(c, "%s\n", orderID)

	// Fetch order details based on the order ID
	product := order(orderID)

	// Render the order template within layouts/main
	return c.Render(
		"order",
		fiber.Map{
			"Title":   "Hello, World!",
			"Product": product,
		},
		"layouts/main",
	)
}

// Define the route handler function
func handleRestricted(c *fiber.Ctx) error {
	// Call home function to fetch products data
	ia_id := restricted(c.Params("uuid"))

	fmt.Fprintf(c, "%v\n", ia_id)
	// Render index template within layouts/main
	return c.Render(
		"restricted",
		fiber.Map{
			"Title": "Hello, Restricted World!",
			"Ia_id": ia_id, // products = []
		},
		"layouts/main",
	)
}

// Define the route handler function
func handleIndex(c *fiber.Ctx) error {
	// Call home function to fetch products data
	products := index()

	// Print products (optional)
	fmt.Printf("total products: %s", products)

	// Render index template within layouts/main
	return c.Render(
		"index",
		fiber.Map{
			"Title":    "Hello, World!",
			"Products": products, // products = []
		},
		"layouts/main",
	)
}
