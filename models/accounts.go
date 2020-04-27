package models

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	utl "my-contacts/utils"
	"os"
	"strings"
)

/*
JWT claims struct
*/
type Token struct {
	UserId uint
	jwt.StandardClaims
}

//  Account struct
type Account struct {
	gorm.Model
	Email string `json:"email"`
	Password string `json:"password"`
	Token string `json:"token";sql:"-"`
}

// Validate incoming user details
func(account *Account) Validate() (map[string] interface{}, bool) {
	if !strings.Contains(account.Email, "@") {
		return utl.Message(false, "Email address is required"), false
	}

	if len(account.Password) < 6 {
		return utl.Message(false, "Password should be more than six characters"), false
	}

	// Email must be unique
	temp := &Account{}

	// check for errors and duplicate emails
	err := GetDB().Table("accounts").Where("email = ?", account.Email).First(temp)
	if err != nil {
		fmt.Println(err)
		return utl.Message(false, "Connection error, please retry"), false
	}

	if temp.Email != "" {
		return utl.Message(false, "Email address is already taken"), false
	}
	return utl.Message(true, "Requirement passed"), true

}

// Create method that adds an account in DB
func(account *Account) Create() map[string] interface{} {
	if resp, ok := account.Validate(); !ok {
		return resp
	}

	// Hash the plain password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)

	GetDB().Create(account)

	if account.ID <= 0 {
		return utl.Message(false, "Failed to create account, try again.")
	}

	// Create new JWT token for newly registered account
	tk := &Token{UserId: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString

	account.Password = "" // delete password

	response := utl.Message(true, "Account created.")
	response["account"] = account
	return response
}

// Login function to sign in users
func Login(email, password string) map[string]interface{} {
	accountPointer := &Account{}
	err := GetDB().Table("accounts").Where("email =?", email).First(accountPointer).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utl.Message(false, "Email address not found")
		}
		return utl.Message(false, "Connection error, please retry")
	}

	err = bcrypt.CompareHashAndPassword([]byte(accountPointer.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return utl.Message(false, "Invalid login credentials")
	}

	// If all went well
	accountPointer.Password = ""

	// Create JWT token
	tk := &Token{UserId: accountPointer.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	accountPointer.Token = tokenString // Store the token in the response

	resp := utl.Message(true, "Logged In")
	resp["account"] = accountPointer
	return resp
}

// GetUser function to fetch a user account, returns an Account pointer
func GetUser(u uint) *Account {
	accountPointer := &Account{}
	GetDB().Table("accounts").Where("id = ?", u).First(accountPointer)
	if accountPointer.Email == "" { // User not found
		return nil
	}

	accountPointer.Password = ""
	return accountPointer
}
