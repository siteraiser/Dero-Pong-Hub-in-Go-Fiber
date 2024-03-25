package productmodel

import "strings"

func sameAddress(iaddr string, user map[string]interface{}) bool {

	//fmt.Printf("s%", user)
	//fmt.Printf("s%", iaddr)
	waddr := B2S(user["wallet"].([]uint8))

	ia_frag := strings.TrimLeft(iaddr, "deroi1")
	wa_frag := strings.TrimLeft(waddr, "dero1")
	ia_frag = ia_frag[0:53]
	wa_frag = wa_frag[0:53]

	return ia_frag == wa_frag
}

/**/
func B2S(bs []uint8) string {
	ba := []byte{}
	for _, b := range bs {
		ba = append(ba, byte(b))
	}
	return string(ba)
}
