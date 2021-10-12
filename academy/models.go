package academy

import (
	"time")


type Institution struct {
	Ins_Key int `json:"inskey"`
	Ins_Name string `json:"name"`
	Ins_Shortname string `json:"shortname"`
	Ins_Slang string `json:"slang"`
}

type InstitutionList struct {
	List    []Institution `json:"list"`
}

type Project struct {
	Prj_Key int `json:"prjkey"`
	Ins_Key int `json:"inskey"`
	Ins_Name string `json:"insname"`
	Prj_Name string `json:"name"`
	Prj_Desc string `json:"desc"`
	Prj_Due time.Time `json:"due"`
	Prj_Words int `json:"words"`
	Prj_Path string `json:"path"`
	Ins_List InstitutionList `json:"inslist"`
}

type ProjectList struct {
	Find_Part string `json:"findpart"`
	List    []Project `json:"list"`
}

type Reference struct {
	Ref_Key int `json:"refkey"`
	Ref_Author string `json:"author"`
	Ref_Year int `json:"year"`
	Ref_Era string `json:"era"`
	Ref_Title string `json:"title"`
	Ref_Publication string `json:"publication"`
	Ref_Publisher string `json:"publisher"`
	Ref_City string `json:"city"`
	Ref_Url string `json:"url"`
	Ref_Path string `json:"path"`
	Ref_Retrieved time.Time `json:"retrieved"`
}

type ReferenceList struct {
	Find_Part string `json:"findpart"`
	List []Reference `json:"list"`
}

type ProjRef struct {
	Prj_Key int `json:"prjkey"`
	Ref_Key int `json:"refkey"`
	Date    time.Time  `json:"date"`
}

type ProjRefList struct {
	Prj_Key int `json:"prjkey"`
	Find_Part string `json:"findpart"`
	List []Reference `json:"list"`
}