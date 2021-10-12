package academy

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/zephryl/honors/common"
)

//------------------------------------------------------------------------------
// CreateReference - insert a Reference
//------------------------------------------------------------------------------
func CreateReference(s *common.System, aReference *Reference) (error) {
	// Open a db, defer close
	db, err := sql.Open("mysql", s.SqlString())
	if err != nil {
		return err;
	}
	defer db.Close()
	
	// Insert Reference, get key
	var vQuery strings.Builder
	vQuery.WriteString("insert into Reference(Ref_Author, Ref_Year, Ref_Era, Ref_Title, Ref_Publication, ")
	vQuery.WriteString("    Ref_Publisher, Ref_City, Ref_Url, Ref_Path, Ref_Retrieved)")
	vQuery.WriteString("values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	res, err := db.Exec(vQuery.String(), aReference.Ref_Author, aReference.Ref_Year, aReference.Ref_Era, aReference.Ref_Title, aReference.Ref_Publication,
		aReference.Ref_Publisher, aReference.Ref_City, aReference.Ref_Url, aReference.Ref_Path, aReference.Ref_Retrieved)
	if err != nil {
		return err;
	}
	var vKey int64;
	vKey, err = res.LastInsertId();
	if (vKey == 0 ) {
		return err;
	}
	aReference.Ref_Key = int(vKey);

	// all ok
	return nil;
}

//------------------------------------------------------------------------------
// ReadReferenceList - select all Reference by filter
//------------------------------------------------------------------------------
func ReadReferenceList(s *common.System, aReferenceList *ReferenceList) error {
	// Open db, defer close
	db, err := sql.Open("mysql", s.SqlString())
	if err != nil {
		return err
	}
	defer db.Close()
	// Fetch Reference rows
	var vQuery strings.Builder
	vQuery.WriteString("select Ref_Key, Ref_Author, Ref_Year, Ref_Era, Ref_Title, Ref_Publication, Ref_Publisher, ")
	vQuery.WriteString("       Ref_City, Ref_Url, Ref_Path, Ref_Retrieved ")
	vQuery.WriteString("from reference ")
	if (aReferenceList.Find_Part != "") {		
		vQuery.WriteString(fmt.Sprintf("where (Ref_Title like '%%%[1]s%%' or Ref_Author like '%%%[1]s%%' ) ", aReferenceList.Find_Part))
	}
	vQuery.WriteString("order by Ref_Author")
	// Execute query into rows
	vRows, err := db.Query(vQuery.String())
	if err != nil {
		return err
	}
	// Iterate rows into ReferenceCollection
	defer vRows.Close()
	for vRows.Next() {
		var vReference Reference
		if err = vRows.Scan(
			&vReference.Ref_Key, &vReference.Ref_Author, &vReference.Ref_Year, &vReference.Ref_Era, &vReference.Ref_Title, &vReference.Ref_Publication,
			&vReference.Ref_Publisher, &vReference.Ref_City, &vReference.Ref_Url, &vReference.Ref_Path, &vReference.Ref_Retrieved); err != nil {
			return err
		}
		aReferenceList.List = append(aReferenceList.List, vReference)
	}
	// get any error encountered during iteration
	err = vRows.Err()
	if err != nil {
		return err
	}
	// all ok
	return nil
}

//------------------------------------------------------------------------------
// ReadReference - select the Reference by key
//------------------------------------------------------------------------------
func ReadReference(s *common.System, aReference *Reference) (error) {
	// Open db, defer close
	db, err := sql.Open("mysql", s.SqlString())
	if err != nil {
		return err;
	}
	defer db.Close();
	// Fetch row
	var vQuery strings.Builder
	vQuery.WriteString("select Ref_Author, Ref_Year, Ref_Era, Ref_Title, Ref_Publication, Ref_Publisher, ")
	vQuery.WriteString("       Ref_City, Ref_Url, Ref_Path, Ref_Retrieved ")
	vQuery.WriteString("from reference ")
	vQuery.WriteString("where Ref_Key = ? ");
	// Execute query into rows
	if err := db.QueryRow(vQuery.String(), &aReference.Ref_Key).Scan(
		&aReference.Ref_Author, &aReference.Ref_Year, &aReference.Ref_Era, &aReference.Ref_Title, &aReference.Ref_Publication,
		&aReference.Ref_Publisher, &aReference.Ref_City, &aReference.Ref_Url, &aReference.Ref_Path, &aReference.Ref_Retrieved); err != nil {
			return err;
		}
	// all ok
	return nil;
}

//------------------------------------------------------------------------------
// Update Reference
//------------------------------------------------------------------------------
func UpdateReference(s *common.System, aReference *Reference) (error) {
	// Open db, defer close
	db, err := sql.Open("mysql", s.SqlString())
	if err != nil {
		return err;
	}
	defer db.Close();
	
	// Update
	var vQuery strings.Builder;
	vQuery.WriteString("update reference ");
	vQuery.WriteString("set Ref_Author=?, Ref_Year=?, Ref_Era=?, Ref_Title=?, Ref_Publication=?, Ref_Publisher=?, ");
	vQuery.WriteString("    Ref_City=?, Ref_Url=?, Ref_Path=?, Ref_Retrieved=? ");
	vQuery.WriteString("where Ref_Key = ? ");
	_, err = db.Exec(vQuery.String(), &aReference.Ref_Author, &aReference.Ref_Year, &aReference.Ref_Era, &aReference.Ref_Title, &aReference.Ref_Publication,
		&aReference.Ref_Publisher, &aReference.Ref_City, &aReference.Ref_Url, &aReference.Ref_Path, &aReference.Ref_Retrieved, &aReference.Ref_Key );
	if err != nil {
		return err;
	}
	return nil;
}

//------------------------------------------------------------------------------
// Delete Reference
//------------------------------------------------------------------------------
func DeleteReference(s *common.System, aReference *Reference) (error) {
	// Open db, defer close
	db, err := sql.Open("mysql", s.SqlString())
	if err != nil {
		return err;
	}
	defer db.Close();
	
	// Delete
	var vQuery strings.Builder;
	vQuery.WriteString("delete from reference ");
	vQuery.WriteString("where Ref_Key = ? ");
	_, err = db.Exec(vQuery.String(), &aReference.Ref_Key );
	if err != nil {
		return err;
	}
	return nil;
}
