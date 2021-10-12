package academy

import (
	"database/sql"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/zephryl/honors/common"
)

//------------------------------------------------------------------------------
// CreateInstitution - insert a Institution
//------------------------------------------------------------------------------
func CreateInstitution(s *common.System, aInstitution *Institution) (error) {
	// Open a db, defer close
	db, err := sql.Open("mysql", s.SqlString())
	if err != nil {
		return err;
	}
	defer db.Close()
	
	// Insert Institution, get key
	var vQuery strings.Builder
	vQuery.WriteString("insert into Institution(Ins_Name, Ins_Shortname, Ins_Slang)")
	vQuery.WriteString("values(?, ?, ?)")
	res, err := db.Exec(vQuery.String(), &aInstitution.Ins_Name, &aInstitution.Ins_Shortname, &aInstitution.Ins_Slang);
	if err != nil {
		return err;
	}
	var vKey int64;
	vKey, err = res.LastInsertId();
	if (vKey == 0 ) {
		return err;
	}
	aInstitution.Ins_Key = int(vKey);

	// all ok
	return nil;
}

//------------------------------------------------------------------------------
// ReadInstitutionList - get a DB, call row reader
//------------------------------------------------------------------------------
func ReadInstitutionList(s *common.System, aInstitutionList *InstitutionList) error {
	// Open db, defer close
	db, err := sql.Open("mysql", s.SqlString())
	if err != nil {
		return err
	}
	defer db.Close()
	// Fetch Institution rows
	if err := ReadInstitutionRows(db, aInstitutionList); err != nil {
		return err;
	}
	return nil;
}

//------------------------------------------------------------------------------
// ReadInstitutionRows - select all Institution by filter
//------------------------------------------------------------------------------
func ReadInstitutionRows(db *sql.DB, aInstitutionList *InstitutionList) error {
	// Fetch Institution rows
	var vQuery strings.Builder
	vQuery.WriteString("select Ins_Key, Ins_Name, Ins_Shortname, Ins_Slang ")
	vQuery.WriteString("from institution ")
	vQuery.WriteString("order by Ins_Name")
	// Execute query into rows
	vRows, err := db.Query(vQuery.String())
	if err != nil {
		return err
	}
	// Iterate rows into InstitutionCollection
	defer vRows.Close()
	for vRows.Next() {
		var vInstitution Institution
		if err = vRows.Scan(
			&vInstitution.Ins_Key, &vInstitution.Ins_Name, &vInstitution.Ins_Shortname, &vInstitution.Ins_Slang); err != nil {
			return err
		}
		aInstitutionList.List = append(aInstitutionList.List, vInstitution)
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
// ReadInstitution - select the Institution by key
//------------------------------------------------------------------------------
func ReadInstitution(s *common.System, aInstitution *Institution) (error) {
	// Open db, defer close
	db, err := sql.Open("mysql", s.SqlString())
	if err != nil {
		return err;
	}
	defer db.Close();
	// Fetch row
	var vQuery strings.Builder
	vQuery.WriteString("select Ins_Name, Ins_Shortname, Ins_Slang ")
	vQuery.WriteString("from institution ")
	vQuery.WriteString("where Ins_Key = ? ");
	// Execute query into rows
	if err := db.QueryRow(vQuery.String(), &aInstitution.Ins_Key).Scan(&aInstitution.Ins_Name, &aInstitution.Ins_Shortname, &aInstitution.Ins_Slang); err != nil {
			return err;
		}
	// all ok
	return nil;
}

//------------------------------------------------------------------------------
// Update Institution
//------------------------------------------------------------------------------
func UpdateInstitution(s *common.System, aInstitution *Institution) (error) {
	// Open db, defer close
	db, err := sql.Open("mysql", s.SqlString())
	if err != nil {
		return err;
	}
	defer db.Close();
	
	// Update
	var vQuery strings.Builder;
	vQuery.WriteString("update institution ");
	vQuery.WriteString("set Ins_Name=?, Ins_Shortname=?, Ins_Slang=? ");
	vQuery.WriteString("where Ins_Key = ? ");
	_, err = db.Exec(vQuery.String(), &aInstitution.Ins_Name, &aInstitution.Ins_Shortname, &aInstitution.Ins_Slang, &aInstitution.Ins_Key );
	if err != nil {
		return err;
	}
	return nil;
}

//------------------------------------------------------------------------------
// Delete Institution
//------------------------------------------------------------------------------
func DeleteInstitution(s *common.System, aInstitution *Institution) (error) {
	// Open db, defer close
	db, err := sql.Open("mysql", s.SqlString())
	if err != nil {
		return err;
	}
	defer db.Close();
	
	// Delete
	var vQuery strings.Builder;
	vQuery.WriteString("delete from institution ");
	vQuery.WriteString("where Ins_Key = ? ");
	_, err = db.Exec(vQuery.String(), &aInstitution.Ins_Key );
	if err != nil {
		return err;
	}
	return nil;
}
