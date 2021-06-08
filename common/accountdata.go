package common

import (
	"database/sql"
	"strings"
	"fmt"
	"errors"
	"bytes"
	"text/template"
	_ "github.com/go-sql-driver/mysql"
)

//------------------------------------------------------------------
// Check whether a register already exists for email
//------------------------------------------------------------------
func CheckRegister(s *System, aRegister *Register) (error) {
	// Open a db, defer close
	db, err := sql.Open("mysql", s.SqlString());
	if err != nil {
		return err;
	}
	defer db.Close();

	// Check register
	var vCount int;
	var vQuery strings.Builder
	vQuery.WriteString("select count(*) ")
	vQuery.WriteString("from register ")
	vQuery.WriteString("where Reg_Email = ? ")
	err = db.QueryRow(vQuery.String(), aRegister.Reg_Email).Scan(&vCount);
	if err != nil {
		return err; // Unknown db error
	}
	if vCount > 0 {
		return errors.New("Account already registered");
	}
	
	return nil; // All good
}

//------------------------------------------------------------------
// Create a register, send a message
//------------------------------------------------------------------
func CreateRegister(s *System, aRegister *Register) (error) {
	// Open a db, defer close
	db, err := sql.Open("mysql", s.SqlString());
	if err != nil {
		return err;
	}
	defer db.Close();

	// Insert register
	var vKey int64;
	var vQuery strings.Builder
	vQuery.WriteString("insert into register(Reg_Firstnames, Reg_Surname, Reg_Email, Reg_Password, Reg_Token, Reg_Date)")
	vQuery.WriteString("values(?, ?, ?, ?, uuid(), now())")
	res, err := db.Exec(vQuery.String(), aRegister.Reg_Firstnames, aRegister.Reg_Surname, aRegister.Reg_Email, aRegister.Reg_Password);
	if err != nil {
		return err; // Unknown db error
	}
	vKey, err = res.LastInsertId();
	if (vKey == 0 ) {
		return err; // Unknown db error
	}
	aRegister.Reg_Key = int(vKey);

	// Get the generated token/date, build a link string
	vQuery.Reset()
	vQuery.WriteString("select Reg_Token, Reg_Date ")
	vQuery.WriteString("from register ")
	vQuery.WriteString("where Reg_Key = ? ")
	err = db.QueryRow(vQuery.String(), aRegister.Reg_Key).Scan(&aRegister.Reg_Token, &aRegister.Reg_Date);
	if err != nil {
		return err; // Unknown db error
	}
	var vSDate = fmt.Sprintf("%d%02d%02d%02d%02d%02d",
		aRegister.Reg_Date.Year(), aRegister.Reg_Date.Month(), aRegister.Reg_Date.Day(), 
		aRegister.Reg_Date.Hour(), aRegister.Reg_Date.Minute(), aRegister.Reg_Date.Second());
	aRegister.Reg_Passthru = fmt.Sprintf("%v?dte=%v&tkn=%v&key=%v", aRegister.Reg_Passthru, vSDate, aRegister.Reg_Token, aRegister.Reg_Key);
	
	// Parse the verify template, pass to sendmail as body
	t, err := template.ParseFiles(fmt.Sprintf("templates/%v", "register.text"));
	if err != nil {
		return err
	};

	buffer := new(bytes.Buffer)
	if err = t.Execute(buffer, aRegister); err != nil {
		return err
	}
	// Send a register verify email
	err = SendMailText(aRegister.Reg_Email, buffer.String());
	if err != nil {
		return err; // Unknown mail error
	}
	
	return nil; // All good
}

//------------------------------------------------------------------
// Read a Register with key
//------------------------------------------------------------------
func ReadRegister(s *System, aRegister *Register) (error) {
	// Open a db, defer close
	db, err := sql.Open("mysql", s.SqlString());
	if err != nil {
		return err;
	}
	defer db.Close();

	// Test register
	var vQuery strings.Builder
	vQuery.WriteString("select Reg_Key, Reg_Firstnames, Reg_Surname, Reg_Email, Reg_Password, Reg_Token, Reg_Date ")
	vQuery.WriteString("from register ")
	vQuery.WriteString("where Reg_Key = ? ")
	if err := db.QueryRow(vQuery.String(), aRegister.Reg_Key).Scan(
		&aRegister.Reg_Key, &aRegister.Reg_Firstnames, &aRegister.Reg_Surname, &aRegister.Reg_Email, &aRegister.Reg_Password,
		&aRegister.Reg_Token, &aRegister.Reg_Date); err != nil {
		return err; // Unknown db error
	}
	
	return nil; // All good
}


//------------------------------------------------------------------
// Check whether a user already exists for email
//------------------------------------------------------------------
func CheckUser(s *System, aUser *User) (error) {
	// Open a db, defer close
	db, err := sql.Open("mysql", s.SqlString());
	if err != nil {
		return err;
	}
	defer db.Close();

	// Check user
	var vCount int;
	var vQuery strings.Builder
	vQuery.WriteString("select count(*) ")
	vQuery.WriteString("from user ")
	vQuery.WriteString("where Usr_Email = ? ")
	err = db.QueryRow(vQuery.String(), aUser.Usr_Email).Scan(&vCount);
	if err != nil {
		return err; // Unknown db error
	}
	if vCount > 0 {
		return errors.New("User already activated"); // User usered
	}
	
	return nil; // All good
}

//------------------------------------------------------------------
// MergeUserRegister - Create a user, delete a register, in a transaction
//------------------------------------------------------------------
func MergeUserRegister(s *System, aUser *User) (error) {
	// Open a db, defer close
	db, err := sql.Open("mysql", s.SqlString());
	if err != nil {
		return err;
	}
	defer db.Close();

	// Create a transaction
	tx, err := db.Begin();
	if err != nil {
		return err; // Unknown db error
	}

	// Insert user
	var vKey int64;
	var vQuery strings.Builder
	vQuery.WriteString("insert into user(Usr_Firstnames, Usr_Surname, Usr_Email, Usr_Password, Usr_Identity, Usr_Token)")
	vQuery.WriteString("values(?, ?, ?, ?, ?, ?)")
	res, err := tx.Exec(vQuery.String(), aUser.Usr_Firstnames, aUser.Usr_Surname, aUser.Usr_Email, aUser.Usr_Password, aUser.Usr_Identity, aUser.Usr_Token);
	if err != nil {
		tx.Rollback();
		return err; // Unknown db error
	}
	vKey, err = res.LastInsertId();
	if (vKey == 0 ) {
		tx.Rollback();
		return err; // Unknown db error
	}
	aUser.Usr_Key = int(vKey);

	// Delete the Register
	vQuery.Reset();
	vQuery.WriteString("delete from register where Reg_Email = ?");
	res, err = tx.Exec(vQuery.String(), aUser.Usr_Email);
	if err != nil {
		tx.Rollback();
		return err; // Unknown db error
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return err; // Unknown db error
	}

	// All good
	return nil;
}
