package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	utl "my-contacts/utils"
)

// Contact struct to store contact information
type Contact struct {
	gorm.Model
	Name string `json:"name"`
	Phone string `gorm:"size:255;not null;unique" json:"phone"`
	UserId uint `json:"user_id"` // The user that this contact belongs to
}

// Validate method to validate the request body data.
// returns message and true if the requirement is met
func(contact *Contact) Validate() (map[string] interface{}, bool) {
	if contact.Name == "" {
		return utl.Message(false, "Contact name is required"), false
	}

	if contact.Phone == "" {
		return utl.Message(false, "Phone number is required"), false
	}

	if contact.UserId <= 0 {
		return utl.Message(false, "User is not recognized"), false
	}

	// Phone should be unique
	temp := &Contact{}
	err := GetDB().Table("contacts").Where("phone = ?", contact.Phone).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		fmt.Println(err)
		return utl.Message(false, "Connection error, please try again"), false
	}

	if temp.Phone != "" {
		return utl.Message(false, "Phone number already exists"), false
	}
	return utl.Message(true, "success"), true
}

// Create method to add a new contact
func(contact *Contact) Create() map[string] interface{} {
	if resp, ok := contact.Validate(); !ok {
		return resp
	}

	GetDB().Create(contact)

	resp := utl.Message(true, "success")
	resp["contact"] = contact
	return resp
}

// GetContact function to retrieve a contact
func GetContact(id uint) *Contact {
	contactPointer := &Contact{}
	err := GetDB().Table("contacts").Where("id = ?", id).First(contactPointer).Error
	if err != nil {
		return nil
	}
	return contactPointer
}

// GetContacts function to retrieve a list of contacts
func GetContacts(user uint) []*Contact {
	contacts := make([]*Contact, 0)
	err := GetDB().Table("contacts").Where("user_id = ?", user).Find(&contacts).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return contacts
}