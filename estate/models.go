package estate

type Estate struct {
	Reg_Key int     `json:"key"`
	Reg_Name string `json:"name"`
	Reg_Code string `json:"code"`
}
type EstateCollection struct {
	List []Estate `json:"list"`
}