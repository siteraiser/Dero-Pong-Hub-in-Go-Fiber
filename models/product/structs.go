package productmodel

// "reflect"
type ProductJSON struct {
	Id        int
	P_type    string `json:"p_type,omitempty"`
	Label     string `json:"label,omitempty"`
	Details   string `json:"details,omitempty"`
	Scid      string `json:"scid,omitempty"`
	Inventory int    `json:"inventory,omitempty"`
	Image     string `json:"image,omitempty"`
	Action    string `json:"action,omitempty"`
}

type IAddressJSON struct {
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
