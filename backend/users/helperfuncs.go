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
	ID                    uint                     `json:"ID"`
	Amount                int64                    `json:"Amount"`
	ReceiverName          string                   `json:"ReceiverName"`
	ReceiverUserID        uint                     `json:"ReceiverUserID"`
	SenderName            string                   `json:"SenderName"`
	SenderUserID          uint                     `json:"SenderUserID"`
	CreatedAt             time.Time                `json:"CreatedAt"`
	Cancelled             bool                     `json:"Cancelled"`
	Type                  database.TransactionType `json:"Type"`
	Comment               string                   `json:"Comment"`
	OriginalTransactionID *uint                    `json:"OriginalTransactionID"`
}

// toTransactionResponse resolves the sender's and receiver's display names and
// converts a database.Transaction into the shape the frontend gets served.
func toTransactionResponse(transaction database.Transaction) (TransactionResponse, error) {
	senderName, err := getNameFromUserID(transaction.SenderUserID)
	if err != nil {
		return TransactionResponse{}, err
	}
	receiverName, err := getNameFromUserID(transaction.ReceiverUserID)
	if err != nil {
		return TransactionResponse{}, err
	}
	comment := ""
	if transaction.Comment != nil {
		comment = *transaction.Comment
	}
	return TransactionResponse{
		ID:                    transaction.ID,
		Amount:                transaction.Amount,
		ReceiverName:          receiverName,
		ReceiverUserID:        transaction.ReceiverUserID,
		SenderName:            senderName,
		SenderUserID:          transaction.SenderUserID,
		CreatedAt:             transaction.CreatedAt,
		Cancelled:             transaction.Cancelled,
		Type:                  transaction.Type,
		Comment:               comment,
		OriginalTransactionID: transaction.OriginalTransactionID,
	}, nil
}

// toTransactionResponses converts a slice of database.Transaction, in place of a loop at each call site.
func toTransactionResponses(transactions []database.Transaction) ([]TransactionResponse, error) {
	responses := make([]TransactionResponse, len(transactions))
	for i, transaction := range transactions {
		response, err := toTransactionResponse(transaction)
		if err != nil {
			return nil, err
		}
		responses[i] = response
	}
	return responses, nil
}
