package main

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

type SubmitTxObject struct {
	Method string `json:"method"`
	Params struct {
		Uuid  string `json:"uuid,omitempty"`
		Ia_id int    `json:"ia_id,omitempty"`
	}
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
