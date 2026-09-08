//
// EPITECH PROJECT, 2026
// CartePro
// File description:
// helperfuncs
//

package users

import (
	"cartepro/database"
	"errors"
	"time"
)

func getNameFromUserID(userID uint) (string, error) {
	var user database.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return "", err
	}
	switch user.Role { //mmmmm tasty switch case
	case database.RoleClient:
		var client database.Client
		if err := database.DB.First(&client, "user_id = ?", userID).Error; err != nil {
			return "", err
		}
		return client.Name, nil
	case database.RolePartner:
		var partner database.Partner
		if err := database.DB.First(&partner, "user_id = ?", userID).Error; err != nil {
			return "", err
		}
		return partner.BusinessName, nil
	case database.RoleAdmin:
		var admin database.Admin
		if err := database.DB.First(&admin, "user_id = ?", userID).Error; err != nil {
			return "", err
		}
		return admin.Name, nil
	default:
		return "", errors.New("user role not recognized")
	}
}

type TransactionResponse struct {
	ID             uint                     `json:"ID"`
	Amount         int64                    `json:"Amount"`
	ReceiverName   string                   `json:"ReceiverName"`
	ReceiverUserID uint                     `json:"ReceiverUserID"`
	SenderName     string                   `json:"SenderName"`
	SenderUserID   uint                     `json:"SenderUserID"`
	CreatedAt      time.Time                `json:"CreatedAt"`
	Cancelled      bool                     `json:"Cancelled"`
	Type           database.TransactionType `json:"Type"`
	Comment        string                   `json:"Comment"`
}
