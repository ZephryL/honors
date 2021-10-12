package common

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

//------------------------------------------------------------------
// Register a user for the first time
//------------------------------------------------------------------
func RegisterCreateHandler(s *System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Decode a register from request body
		d := json.NewDecoder(r.Body)
		d.DisallowUnknownFields()
		var vRegister = new(Register)
		err := d.Decode(&vRegister)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Registration failed at deserialization: %v", err.Error())))
			return
		}

		// Check for existing register
		if err := CheckRegister(s, vRegister); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(fmt.Sprintf("Registration failed at register check: %v", err.Error())))
			return
		}

		// Check for existing user
		var vUser = new(User)
		vUser.Usr_Email = vRegister.Reg_Email
		err, vFound := CheckUser(s, vUser)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Registration failed at user check: %v", err.Error())))
			return
		}
		if vFound {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(fmt.Sprintf("Registration failed at user check: %v", "User already registered")))
			return
		}

		// Create a register
		err = CreateRegister(s, vRegister)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Registration failed at new user: %v", err.Error())))
			return
		}

		// Send sanitised Register as json
		vRegister.Reg_Password = ""
		vRegister.Reg_Token = ""
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
		d := json.NewDecoder(r.Body)
		d.DisallowUnknownFields()
		var vVerify = new(Verify)
		err := d.Decode(&vVerify)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Verification failed at deserialization: %v", err.Error())))
			return
		}

		// Get a register with key
		var vRegister = new(Register)
		vRegister.Reg_Key = vVerify.Vrf_Key
		err = ReadRegister(s, vRegister)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(fmt.Sprintf("Verification failed at register retrieve: %v", err.Error())))
			return
		}

		// Test whether incoming token matches persisted token
		if vRegister.Reg_Token != vVerify.Vrf_Token {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(fmt.Sprintf("Verification failed. Token mismatch with input token %v", vVerify.Vrf_Token)))
			return
		}

		// Test whether the email exists on user
		var vUser = new(User)
		vUser.Usr_Email = vRegister.Reg_Email
		err, vFound := CheckUser(s, vUser)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Verification failed at user check '%v'", err.Error())))
			return
		}
		if vFound {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(fmt.Sprintf("Verification failed. User account already exists for '%v'", vRegister.Reg_Email)))
			return
		}

		// Insert a user with register attributes, delete the register
		vUser.Usr_Name = vRegister.Reg_Name
		vUser.Usr_Surname = vRegister.Reg_Surname
		vUser.Usr_Email = vRegister.Reg_Email
		vUser.Usr_Password = vRegister.Reg_Password
		vUser.Usr_Identity = ""
		vUser.Usr_Token = ""
		err = MergeUserRegister(s, vUser)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(fmt.Sprintf("Verification failed: %v", err.Error())))
			return
		}

		// Reply with the sanitized user
		vUser.Usr_Password = ""
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
		var vRegister = new(Register)
		vRegKey, err := strconv.Atoi(params["reg-key"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v - Register url '%v' does not contain an identifier: %v", http.StatusInternalServerError, params, err.Error())))
			return
		}
		vRegister.Reg_Key = vRegKey

		// Get a register with key
		err = ReadRegister(s, vRegister)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(fmt.Sprintf("Verification failed: %v", err.Error())))
			return
		}

		// Send sanitised Register as json
		vRegister.Reg_Password = ""
		vRegister.Reg_Token = ""
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
		d := json.NewDecoder(r.Body)
		d.DisallowUnknownFields()
		var vLogin = new(Login)
		err := d.Decode(&vLogin)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v - Decode Login Error: %v", http.StatusInternalServerError, err.Error())))
			return
		}

		// Get a UsrKey and Token with Login email and password
		vUsrKey, vToken, err := getKeyToken(s, vLogin)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(fmt.Sprintf("HTTP %v - Not registered: %v", http.StatusUnauthorized, err.Error())))
			return
		}

		// Encode user key and token as a cookie value - Hash keys should be at least 32 bytes long, Block keys should be 16 bytes (AES-128) or 32 bytes (AES-256) long.
		value := map[string]string{"key": strconv.Itoa(vUsrKey), "token": vToken}
		encoded, err := s.Cookie.Encode("zephacad-cookie", value)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("HTTP %v - Encode Error: %v", http.StatusInternalServerError, err.Error())))
			return
		}

		// Create a cookie, set all attributes
		cookie := &http.Cookie{Name: "zephacad-cookie"}
		cookie.Value = encoded
		SetCookieDefaults(cookie)
		http.SetCookie(w, cookie)

		// Send the passthru argument from login straight back as an unnamed type on the fly
		var vArbitrary = new(Arbitrary)
		vArbitrary.Passthru = vLogin.Passthru
		if err := json.NewEncoder(w).Encode(vArbitrary); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Your call was politely taken, but the JSON Encoder failed to format a reply!"))
		}
	}
}

//------------------------------------------------------------------
// Logout
//------------------------------------------------------------------
func LogoutHandler(s *System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Test if a cookie was received
		cookie, err := r.Cookie("zephacad-cookie")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(fmt.Sprintf("HTTP %v - Couldn't find me no zephacad-cookie: %v", http.StatusUnauthorized, err.Error())))
			return
		}
		// User is authentic and authorized - Rewrite cookie defaults, set response cookie
		cookie.Expires = time.Now().Add(-7 * 24 * time.Hour) // a week in the past
		cookie.MaxAge = 0                                    // can't age
		http.SetCookie(w, cookie)
	}
}

//------------------------------------------------------------------
// Manage a forgot password request
//------------------------------------------------------------------
func ForgotCreateHandler(s *System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Decode a forgot from request body
		d := json.NewDecoder(r.Body)
		d.DisallowUnknownFields()
		var vForgot = new(Forgot)
		err := d.Decode(&vForgot)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Forgot Password failed at deserialization: %v", err.Error())))
			return
		}

		// Check for existing user
		var vUser = new(User)
		vUser.Usr_Email = vForgot.Fgt_Email
		err, vForgot.Fgt_Exists = CheckUser(s, vUser)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Forgot Password failed at check user: %v", err.Error())))
			return
		}

		// Create a forgot
		err = CreateForgot(s, vForgot)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Forgot Password failed at create forgot: %v", err.Error())))
			return
		}

		// Send sanitized Forgot as json
		vForgot.Fgt_Token = ""
		if err := json.NewEncoder(w).Encode(vForgot); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Your call was politely taken, but the JSON Encoder failed to format a reply!"))
		}

	}
}

//------------------------------------------------------------------
// Verify a registration string
//------------------------------------------------------------------
func ForgotVerifyHandler(s *System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Decode a Forgot from request body
		d := json.NewDecoder(r.Body)
		d.DisallowUnknownFields()
		var vNewPassword = new(NewPassword)
		err := d.Decode(&vNewPassword)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Verification failed at deserialization: %v", err.Error())))
			return
		}

		// Get a forgot with key
		var vForgot = new(Forgot)
		vForgot.Fgt_Key = vNewPassword.Fgt_Key
		err = ReadForgot(s, vForgot)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(fmt.Sprintf("Verification failed at forgot retrieve: %v", err.Error())))
			return
		}

		// Test whether incoming token matches persisted token
		if vForgot.Fgt_Token != vNewPassword.Npw_Token {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(fmt.Sprintf("Verification failed. Token mismatch with input token %v", vNewPassword.Npw_Token)))
			return
		}

		// Test whether the email exists on user
		var vUser = new(User)
		vUser.Usr_Email = vForgot.Fgt_Email
		err, vFound := CheckUser(s, vUser)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Verification failed at user check '%v'", err.Error())))
			return
		}
		if !vFound {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(fmt.Sprintf("Verification failed. User account does not exist for '%v'", vForgot.Fgt_Email)))
			return
		}

		// Update a user's password, delete the forgot
		vUser.Usr_Email = vForgot.Fgt_Email
		vUser.Usr_Password = vNewPassword.Npw_Password
		err = MergeUserForgot(s, vUser)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(fmt.Sprintf("Verification failed: %v", err.Error())))
			return
		}

		// Reply with the sanitized user
		vUser.Usr_Password = ""
		if err := json.NewEncoder(w).Encode(vUser); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Your call was politely taken, but the JSON Encoder failed to format a reply!"))
		}
	}
}
