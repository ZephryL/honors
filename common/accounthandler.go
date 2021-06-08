package common

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"github.com/gorilla/mux"
)

// func RegisterHandler(s *System) func(http.ResponseWriter, *http.Request) {
// 	return func(w http.ResponseWriter, r *http.Request) {	
// 	}
// }

//------------------------------------------------------------------
// Register a user for the first time
//------------------------------------------------------------------
func RegisterCreateHandler(s *System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Decode a register from request body
		d := json.NewDecoder(r.Body);
		d.DisallowUnknownFields();
		var vRegister = new(Register);
		err := d.Decode(&vRegister)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError);
			w.Write([]byte(fmt.Sprintf("Registration failed at deserialization: %v", err.Error())));
			return;
		}

		// Check for existing register
		if err := CheckRegister(s, vRegister); err != nil {
			w.WriteHeader(http.StatusUnauthorized);
			w.Write([]byte(fmt.Sprintf("Registration failed at register check: %v", err.Error())));
			return;
		}
		
		// Check for existing user
		var vUser = new(User);
		vUser.Usr_Email = vRegister.Reg_Email;
		if err := CheckUser(s, vUser); err != nil {
			w.WriteHeader(http.StatusUnauthorized);
			w.Write([]byte(fmt.Sprintf("Registration failed at user check: %v", err.Error())));
			return;
		}

		// Create a register
		err = CreateRegister(s, vRegister);
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError);
			w.Write([]byte(fmt.Sprintf("Registration failed at new user: %v", err.Error())));
			return;
		}

		// Send sanitised Register as json
		vRegister.Reg_Password = "";
		vRegister.Reg_Token = "";
		if err := json.NewEncoder(w).Encode(vRegister); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Your call was politely taken, but the JSON Encoder failed to format a reply!"))
		}
		
	}
}

//------------------------------------------------------------------
// Verify a registration string
//------------------------------------------------------------------
func RegisterVerifyHandler(s *System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Decode a Verify from request body
		d := json.NewDecoder(r.Body);
		d.DisallowUnknownFields();
		var vVerify = new(Verify);
		err := d.Decode(&vVerify)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError);
			w.Write([]byte(fmt.Sprintf("Verification failed at deserialization: %v", err.Error())));
			return;
		}

		// Get a register with key
		var vRegister = new(Register);
		vRegister.Reg_Key = vVerify.Reg_Key;
		err = ReadRegister(s, vRegister);
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized);
			w.Write([]byte(fmt.Sprintf("Verification failed: %v", err.Error())));
			return;
		}

		// Test whether incoming token matches persisted token
		if vRegister.Reg_Token != vVerify.Reg_Token {
			w.WriteHeader(http.StatusUnauthorized);
			w.Write([]byte(fmt.Sprintf("Verification failed: Token mismatch with input token %v", vVerify.Reg_Token)));
			return;
		}

		// Test whether the email exists on user
		var vUser = new(User);
		vUser.Usr_Email = vRegister.Reg_Email;
		if err := CheckUser(s, vUser); err != nil {
			w.WriteHeader(http.StatusUnauthorized);
			w.Write([]byte(fmt.Sprintf("Verification failed. User account already exists for '%v'", vRegister.Reg_Email)));
			return;
		}

		// Insert a user with register attributes, delete the register
		vUser.Usr_Firstnames = vRegister.Reg_Firstnames;
		vUser.Usr_Surname = vRegister.Reg_Surname;
		vUser.Usr_Email = vRegister.Reg_Email;
		vUser.Usr_Password = vRegister.Reg_Password;
		vUser.Usr_Identity = "";
		vUser.Usr_Token = "";
		err = MergeUserRegister(s, vUser);
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized);
			w.Write([]byte(fmt.Sprintf("Verification failed: %v", err.Error())));
			return;
		}

		// Reply with the sanitized user
		vUser.Usr_Password = "";
		if err := json.NewEncoder(w).Encode(vUser); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Your call was politely taken, but the JSON Encoder failed to format a reply!"))
		}
	}
}

//------------------------------------------------------------------
// Get a register by key
//------------------------------------------------------------------
func RegisterReadHandler(s *System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the Register key from http Request
		params := mux.Vars(r)
		var vRegister = new(Register);
		vRegKey, err := strconv.Atoi(params["reg-key"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v - Register url '%v' does not contain an identifier: %v", http.StatusInternalServerError, params, err.Error())))
			return
		}
		vRegister.Reg_Key = vRegKey;

		// Get a register with key
		err = ReadRegister(s, vRegister);
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized);
			w.Write([]byte(fmt.Sprintf("Verification failed: %v", err.Error())));
			return;
		}

		// Send sanitised Register as json
		vRegister.Reg_Password = "";
		vRegister.Reg_Token = "";
		if err := json.NewEncoder(w).Encode(vRegister); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Your call was politely taken, but the JSON Encoder failed to format a reply!"))
		}
		
	}
}

//------------------------------------------------------------------
// Login - return cookie only (i.e. 201 with no result), or error
//------------------------------------------------------------------
func LoginHandler(s *System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Decode a Login from request body
		d := json.NewDecoder(r.Body);
		d.DisallowUnknownFields();
		var vLogin = new(Login);
		err := d.Decode(&vLogin)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError);
			w.Write([]byte(fmt.Sprintf("HTTP %v - Decode Login Error: %v", http.StatusInternalServerError, err.Error())));
			return;
		}

		// Get a UsrKey and Token with Login email and password
		vUsrKey, vToken, err := getKeyToken(s, vLogin);
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized);
			w.Write([]byte(fmt.Sprintf("HTTP %v - Not registered: %v", http.StatusUnauthorized, err.Error())));
			return;
		}

		// Encode user key and token as a cookie value - Hash keys should be at least 32 bytes long, Block keys should be 16 bytes (AES-128) or 32 bytes (AES-256) long.	
		value := map[string]string{"key": strconv.Itoa(vUsrKey), "token": vToken}
		encoded, err := s.Cookie.Encode("zeph-cookie", value); 
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError);
			w.Write([]byte(fmt.Sprintf("HTTP %v - Encode Error: %v", http.StatusInternalServerError, err.Error())));
			return;
		}

		// Create a cookie, set all attributes
		cookie := &http.Cookie{Name: "zeph-cookie"};
		cookie.Value = encoded;
		SetCookieDefaults(cookie);
		http.SetCookie(w, cookie)

		// Send the passthru argument from login straight back as an unnamed type on the fly
		var vArbitrary = new(Arbitrary);
		vArbitrary.Passthru = vLogin.Passthru;
		if err := json.NewEncoder(w).Encode(vArbitrary); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Your call was politely taken, but the JSON Encoder failed to format a reply!"))
		}
	}
}

//------------------------------------------------------------------
// Get a Usr_Key and Usr_Token with Login details
//------------------------------------------------------------------
func getKeyToken(s *System, aLogin *Login) (int, string, error) {
	// Open a db, defer close
	db, err := sql.Open("mysql", s.SqlString());
	if err != nil {
		return 0, "", err;
	}
	defer db.Close();
	// Create a key and an Account
	var vUsrKey int;
	var vToken string;
	var vQuery strings.Builder
	// Update the GUID of a User with a matching auth set of UserID and Password
	vQuery.WriteString("update user set Usr_Token = uuid() ")
	vQuery.WriteString("where Usr_Email = ? ")
	vQuery.WriteString("and   Usr_Password = ? ")
	res, err := db.Exec(vQuery.String(), aLogin.Email, aLogin.Token);
	if err != nil {
		return 0, "", err;
	}
	affected, err := res.RowsAffected();
	if err != nil {
		return 0, "", err;
	}
	if affected == 0 {
		return 0, "", errors.New("");
	}
	// Then return the GUID
	vQuery.Reset();
	vQuery.WriteString("select Usr_Key, Usr_Token ")
	vQuery.WriteString("from user ")
	vQuery.WriteString("where Usr_Email = ? ")
	vQuery.WriteString("and   Usr_Password = ? ")
	err = db.QueryRow(vQuery.String(), aLogin.Email, aLogin.Token).Scan(&vUsrKey, &vToken);
	if err != nil {
		return 0, "", err;
	}
	return vUsrKey, vToken, nil;
}

//------------------------------------------------------------------
// Get a User of Key, return the User's Token
//------------------------------------------------------------------
func Authenticate(s *System, aToken *Token) (error) {
	// Open a db, defer close
	db, err := sql.Open("mysql", s.SqlString());
	if err != nil {
		return err;
	}
	defer db.Close();
	// Create a token placeholder and fetch a user token
	var vToken string = "";
	var vQuery strings.Builder
	vQuery.WriteString("select Usr_Token ")
	vQuery.WriteString("from user ")
	vQuery.WriteString("where Usr_Key = ? ")
	err = db.QueryRow(vQuery.String(), aToken.Key).Scan(&vToken);
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("Invalid User"); // User not found
		} else {
			return err; // Unknown db error
		}
	}
	if vToken != aToken.Token {
		return errors.New("Invalid User/Token"); // User/Token mismatch
	}
	return nil; // All good
}

//------------------------------------------------------------------
// Get user access to route
//------------------------------------------------------------------
func Authorize(s *System, aRoute string, aMethod string) (error) {
	// Open a db, defer close
	db, err := sql.Open("mysql", s.SqlString());
	if err != nil {
		return err;
	}
	defer db.Close();
	// Create a token placeholder and fetch a user token
	var vCount int = 0;
	var vQuery strings.Builder
	vQuery.WriteString("select count(*) ")
	vQuery.WriteString("from route ")
	vQuery.WriteString("where ? REGEXP Rte_Regex ")
	vQuery.WriteString("and   Rte_Method = ? ")
	err = db.QueryRow(vQuery.String(), aRoute, aMethod).Scan(&vCount);
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("User not authorized"); // User not found
		} else {
			return err; // Unknown db error
		}
	}
	if vCount == 0 {
		return errors.New("User not authorized");
	}
	return nil; // All good
}

//------------------------------------------------------------------
// Logout
//------------------------------------------------------------------
func LogoutHandler(s *System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

// func UserHandler(s *System) func(http.ResponseWriter, *http.Request) {
// 	return func(w http.ResponseWriter, r *http.Request) {	
// 	}
// }
