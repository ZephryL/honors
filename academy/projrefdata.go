package academy

import (
	"database/sql"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/zephryl/honors/common"
)

//------------------------------------------------------------------------------
// CreateProjRef - insert a ProjRef/ProjRef intersect
//------------------------------------------------------------------------------
func CreateProjRef(s *common.System, aProjRef *ProjRef) (error) {
	// Open a db, defer close
	db, err := sql.Open("mysql", s.SqlString())
	if err != nil {
		return err;
	}
	defer db.Close()
	
	// Insert ProjRef, get key
	var vQuery strings.Builder
	vQuery.WriteString("insert into ProjRef(Prj_Key, Usr_Key, Ref_Key, Ref_Date)")
	vQuery.WriteString("values(?, ?, ?, ?)")
	res, err := db.Exec(vQuery.String(), s.Usr_Key, aProjRef.Prj_Key, aProjRef.Ref_Key, time.Now());
	if err != nil {
		return err;
	}
	var vKey int64;
	vKey, err = res.LastInsertId();
	if (vKey == 0 ) {
		return err;
	}
	aProjRef.Ref_Key = int(vKey);

	// all ok
	return nil;
}

//------------------------------------------------------------------------------
// ReadProjRefList - select all ProjRef by filter
//------------------------------------------------------------------------------
func ReadProjRefList(s *common.System, aProjRefList *ProjRefList) error {
	// Open db, defer close
	db, err := sql.Open("mysql", s.SqlString())
	if err != nil {
		return err
	}
	defer db.Close()
	// Fetch ProjRef rows
	var vQuery strings.Builder
	vQuery.WriteString("select Ref_Key, Ref_Author, Ref_Year, Ref_Era, Ref_Title, Ref_Publication, Ref_Publisher, ")
	vQuery.WriteString("       Ref_City, Ref_Url, Ref_Path ")
	vQuery.WriteString("from reference ")
	vQuery.WriteString("where Ref_Key in (select REF_Key from projref where Prj_Key = ? and Usr_Key = ?) ")
	vQuery.WriteString("order by Ref_Author")
	// Execute query into rows
	vRows, err := db.Query(vQuery.String(), aProjRefList.Prj_Key, s.Usr_Key)
	if err != nil {
		return err
	}
	// Iterate rows into ProjRefCollection
	defer vRows.Close()
	for vRows.Next() {
		var vReference Reference
		if err = vRows.Scan(
			&vReference.Ref_Key, &vReference.Ref_Author, &vReference.Ref_Year, &vReference.Ref_Era, &vReference.Ref_Title, &vReference.Ref_Publication,
			&vReference.Ref_Publisher, &vReference.Ref_City, &vReference.Ref_Url, &vReference.Ref_Path); err != nil {
			return err
		}
		aProjRefList.List = append(aProjRefList.List, vReference)
	}
	// get any error encountered during iteration
	err = vRows.Err()
	if err != nil {
		return err
	}
	// all ok
	return nil
}
