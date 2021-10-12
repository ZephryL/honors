package academy

import (
	"database/sql"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/zephryl/honors/common"
)

//------------------------------------------------------------------------------
// CreateProject - insert a Project
//------------------------------------------------------------------------------
func CreateProject(s *common.System, aProject *Project) (error) {
	// Open a db, defer close
	db, err := sql.Open("mysql", s.SqlString())
	if err != nil {
		return err;
	}
	defer db.Close()
	
	// Insert Project, get key
	var vQuery strings.Builder
	vQuery.WriteString("insert into Project(Usr_Key, Ins_Key, Prj_Name, Prj_Desc, Prj_Due, Prj_Words, Prj_Path )")
	vQuery.WriteString("values(?, ?, ?, ?, ?, ?, ?) ")
	res, err := db.Exec(vQuery.String(), s.Usr_Key, aProject.Ins_Key, aProject.Prj_Name, aProject.Prj_Desc, 
	    aProject.Prj_Due, aProject.Prj_Words, aProject.Prj_Path);
	if err != nil {
		return err;
	}
	var vKey int64;
	vKey, err = res.LastInsertId();
	if (vKey == 0 ) {
		return err;
	}
	aProject.Prj_Key = int(vKey);

	// all ok
	return nil;
}

//------------------------------------------------------------------------------
// ReadProjectList - select all Project by filter
//------------------------------------------------------------------------------
func ReadProjectList(s *common.System, aProjectList *ProjectList) error {
	// Open db, defer close
	db, err := sql.Open("mysql", s.SqlString())
	if err != nil {
		return err
	}
	defer db.Close()
	// Fetch Project rows
	var vQuery strings.Builder
	vQuery.WriteString("select prj.Prj_Key, prj.Ins_Key, ins.Ins_Name, prj.Prj_Name, prj.Prj_Desc, prj.Prj_Due, prj.Prj_Words, prj.Prj_Path ")
	vQuery.WriteString("from project prj, institution ins ")
	vQuery.WriteString("where prj.Ins_Key = ins.Ins_Key ")
	vQuery.WriteString("and   prj.Usr_Key = ? ")
	vQuery.WriteString("order by prj.Prj_Name")
	// Execute query into rows
	vRows, err := db.Query(vQuery.String(), s.Usr_Key)
	if err != nil {
		return err
	}
	// Iterate rows into ProjectCollection
	defer vRows.Close()
	for vRows.Next() {
		var vProject Project
		if err = vRows.Scan(
			&vProject.Prj_Key, &vProject.Ins_Key, &vProject.Ins_Name, &vProject.Prj_Name, &vProject.Prj_Desc, 
			&vProject.Prj_Due, &vProject.Prj_Words, &vProject.Prj_Path); err != nil {
			return err
		}
		aProjectList.List = append(aProjectList.List, vProject)
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
// ReadProject - select the Project by key
//------------------------------------------------------------------------------
func ReadProject(s *common.System, aProject *Project) (error) {
	// Open db, defer close
	db, err := sql.Open("mysql", s.SqlString())
	if err != nil {
		return err;
	}
	defer db.Close();
	// Fetch row
	var vQuery strings.Builder
	vQuery.WriteString("select prj.Prj_Key, prj.Ins_Key, ins.Ins_Name, prj.Prj_Name, prj.Prj_Desc, prj.Prj_Due, prj.Prj_Words, prj.Prj_Path ");
	vQuery.WriteString("from project prj, institution ins ")
	vQuery.WriteString("where prj.Ins_Key = ins.Ins_Key ")
	vQuery.WriteString("and   prj.Prj_Key = ? ");
	vQuery.WriteString("and   prj.Usr_Key = ? ");
	// Execute query into rows
	if err := db.QueryRow(vQuery.String(), &aProject.Prj_Key, &s.Usr_Key).Scan(
		&aProject.Prj_Key, &aProject.Ins_Key, &aProject.Ins_Name, &aProject.Prj_Name, &aProject.Prj_Desc, &aProject.Prj_Due,
		&aProject.Prj_Words, &aProject.Prj_Path); err != nil {
		return err;
		}
	// fetch list of institutions for selection
	if err := ReadInstitutionRows(db, &aProject.Ins_List); err != nil {
		return err;
	}
	// all ok
	return nil;
}

//------------------------------------------------------------------------------
// Update Project
//------------------------------------------------------------------------------
func UpdateProject(s *common.System, aProject *Project) (error) {
	// Open db, defer close
	db, err := sql.Open("mysql", s.SqlString())
	if err != nil {
		return err;
	}
	defer db.Close();
	
	// Update
	var vQuery strings.Builder;
	vQuery.WriteString("update project ");
	vQuery.WriteString("set Ins_Key = ?, Prj_Name = ?, Prj_Desc = ?, Prj_Due = ?, Prj_Words = ?, Prj_Path = ? ");
	vQuery.WriteString("where Prj_Key = ? ");
	_, err = db.Exec(vQuery.String(), &aProject.Ins_Key, &aProject.Prj_Name, &aProject.Prj_Desc, 
	                                  &aProject.Prj_Due, &aProject.Prj_Words, &aProject.Prj_Path, &aProject.Prj_Key );
	if err != nil {
		return err;
	}
	// affected, err := res.RowsAffected();
	// if err != nil {
	// 	return err;
	// }
	// if affected == 0 {
	// 	return errors.New("No Project was updated. Are the keys supplied correctly?");
	// }
	// all ok
	return nil;
}

//------------------------------------------------------------------------------
// Update Project
//------------------------------------------------------------------------------
func DeleteProject(s *common.System, aProject *Project) (error) {
	// Open db, defer close
	db, err := sql.Open("mysql", s.SqlString())
	if err != nil {
		return err;
	}
	defer db.Close();
	
	// Delete
	var vQuery strings.Builder;
	vQuery.WriteString("delete from project ");
	vQuery.WriteString("where Prj_Key = ? ");
	_, err = db.Exec(vQuery.String(), &aProject.Prj_Key );
	if err != nil {
		return err;
	}
	return nil;
}
