package main

import (
	productmodel "goserver/models/product"
	usermodel "goserver/models/user"
)

func register(
	obj RegisterObject, // submit Registration Object into database
) map[string]interface{} { // return a uuid as their registration password

	//	fmt.Println("len(userObj):", len(userObj))
	//	fmt.Println("obj.Method:", obj.Method)

	var reg = usermodel.InsertUser(
		obj.Params.Username,
		obj.Params.Wallet,
	)

	if reg != "failed" {
		return map[string]interface {
		}{
			"success": true,
			"reg":     reg,
		}
	}

	return map[string]interface {
	}{
		"success": false,
	}
}

func submitProduct(
	userObj map[string]interface{}, // curl -u secret:ca2154ca-5fb9-4884-b1ec-eee084de40fb
	obj SubmitProductObject, // and we are submitting a product
) map[string]interface{} {

	response := map[string]interface {
	}{
		"success": true,
	}
	if productmodel.SubmitProduct(userObj, obj.Params) {
		return response
	}
	response = map[string]interface {
	}{
		"success": false,
	}
	return response
}

func submitIAddress(
	userObj map[string]interface { // similar pattern, we have a user
	},
	obj SubmitIAddressObject, // and we are submitting a "check-out"
) map[string]interface {
} {
	response := map[string]interface {
	}{
		"success": true,
	}

	if productmodel.SubmitIAddress(
		userObj,
		obj.Params,
	) {
		return response
	}

	response = map[string]interface {
	}{
		"success": false,
	}

	return response
}

func submitTx(
	userObj map[string]interface { // similar pattern, we have a user
	},
	obj SubmitTxObject, // and we are submitting a "paid"
) map[string]interface {
} {
	response := map[string]interface {
	}{
		"success": true,
	}
	if productmodel.SubmitTx(
		userObj,
		obj.Params,
	) {
		return response
	}

	response = map[string]interface {
	}{
		"success": false,
	}

	return response
}
