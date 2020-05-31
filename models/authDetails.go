package models

import (
	"github.com/jinzhu/gorm"
	"github.com/twinj/uuid"
)

type AuthDetails struct {
	gorm.Model
	UserID uint `gorm:"not null" json:"user_id"`
	AuthUUID string `gorm:"size:255;not null" json:"auth_uuid"`
}

// FetchAuthDetails function to fetch authentication details of a user
func FetchAuthDetails(tokenClaims *Token) (*AuthDetails, error) {
	authDetailsPointer := &AuthDetails{}
	err := GetDB().Table("auth_details").Where("user_id = ? AND auth_uuid = ?",
		tokenClaims.UserId, tokenClaims.AuthUUID).First(authDetailsPointer).Error
	if err != nil {
		return nil, err
	}
	return authDetailsPointer, nil
}

// DeleteAuthDetails function to delete authentication details once a user logs out
func DeleteAuthDetails(tokenClaims *Token) error {
	authDetailsPointer := &AuthDetails{}
	err := GetDB().Table("auth_details").Where("user_id = ? AND auth_uuid = ?",
		tokenClaims.UserId, tokenClaims.AuthUUID).Delete(authDetailsPointer).Error
	if err != nil {
		return err
	}
	return nil
}

// SaveAuthDetails function to save authentication details once a user login/sign up
func SaveAuthDetails(id uint) error {
	authDetailsPointer := &AuthDetails{}
	authDetailsPointer.UserID = id
	authDetailsPointer.AuthUUID = uuid.NewV4().String()

	err := GetDB().Create(*authDetailsPointer).Error
	if err != nil {
		return err
	}
	return nil
}

