package world

type Region struct {
	Reg_Key int     `json:"key"`
	Reg_Name string `json:"name"`
	Reg_Code string `json:"code"`
}
type RegionCollection struct {
	List []Region `json:"list"`
}