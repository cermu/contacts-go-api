package models

import (
	"fmt"
	"github.com/badoux/checkmail"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	_ "github.com/twinj/uuid"
	"golang.org/x/crypto/bcrypt"
	utl "my-contacts/utils"
	"os"
)

/*
JWT claims struct, has one claim, UserId
*/
type Token struct {
	UserId uint
	AuthUUID string
	jwt.StandardClaims
}

//  Account struct
type Account struct {
	gorm.Model
	Email string `gorm:"size:255;not null;unique" json:"email"`
	Password string `json:"password"`
	Token string `json:"token";sql:"-"`
}

// Validate incoming data from client
func(account *Account) Validate() (map[string] interface{}, bool) {
	//if !strings.Contains(account.Email, "@") {
	//	return utl.Message(false, "Email address is required"), false
	//}

	if err := checkmail.ValidateFormat(account.Email); err != nil {
		return utl.Message(false, "Provide a valid email address"), false
	}

	if len(account.Password) < 6 {
		return utl.Message(false, "Password should be more than six characters"), false
	}

	// Email must be unique
	temp := &Account{}

	// fmt.Println(account.Email)
	// fmt.Println(temp.Email)
	// fmt.Println("=================================")

	// check for errors and duplicate emails
	err := GetDB().Table("accounts").Where("email = ?", account.Email).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		fmt.Println(err)
		// fmt.Println("check here")
		return utl.Message(false, "Connection error, please retry"), false
	}

	if temp.Email != "" {
		return utl.Message(false, "Email address is already taken"), false
	}
	return utl.Message(true, "Requirements passed"), true

}

// Create method that create a new account and generates a JWT token
func(account *Account) Create() map[string] interface{} {
	if resp, ok := account.Validate(); !ok {
		return resp
	}

	// Hash the plain password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)

	GetDB().Create(account) // Save the account in DB

	if account.ID <= 0 {
		return utl.Message(false, "Failed to create account, try again.")
	}

	// Save auth details
	authData, authError := SaveAuthDetails(account.ID)
	if authError != nil {
		return utl.Message(false, "The following error occurred: "+ authError.Error())
	}

	// Create new JWT token for newly registered account
	tk := &Token{UserId: account.ID, AuthUUID: authData.AuthUUID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk) // Add claim 'tk' to the token
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString

	account.Password = "" // delete password

	response := utl.Message(true, "Account created.")
	response["account"] = account
	return response
}

// Login function to authenticate a user and generate a JWT token
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

	// Save auth details
	authData, authError := SaveAuthDetails(accountPointer.ID)
	if authError != nil {
		return utl.Message(false, "The following error occurred: " + authError.Error())
	}

	// Create JWT token
	tk := &Token{UserId: accountPointer.ID, AuthUUID: authData.AuthUUID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	accountPointer.Token = tokenString // Store the token in the response

	resp := utl.Message(true, "Logged In")
	resp["account"] = accountPointer
	return resp
}

// Logout function to delete auth details and log out a user
func Logout(id uint) map[string]interface{} {
	err := DeleteAuthDetails(id)
	if err != nil {
		return utl.Message(false, "The following error occurred: " + err.Error())
	}

	resp := utl.Message(true, "Logged out successfully")
	return resp
}

// GetUser function to fetch a user account, returns an Account pointer
func GetUser(u uint) *Account {
	accountPointer := &Account{}
	err := GetDB().Table("accounts").Where("id = ?", u).First(accountPointer).Error
	if err != nil {
		return nil
	}

	if accountPointer.Email == "" { // User not found
		return nil
	}

	accountPointer.Password = ""
	return accountPointer
}

// GetUserByMail function to fetch and return an account using the passed email address
func GetUserByMail(email string) *Account {
	accountPointer := &Account{}
	err := GetDB().Table("accounts").Where("email = ?", email).First(accountPointer).Error
	if err != nil {
		return nil
	}

	accountPointer.Password = ""
	return accountPointer
}
