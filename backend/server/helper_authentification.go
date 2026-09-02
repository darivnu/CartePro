//
// EPITECH PROJECT, 2026
// CartePro
// File description:
// This file contains a bunch of functions to help the other handlers
//

package server

import (
	"cartepro/database"
	"net/http"
)

func GetClientFromSession(r *http.Request) (*database.User, *database.Client, error) {
	session, err := GetSessionFromRequest(r)
	if err != nil {
		return nil, nil, err
	}

	var user database.User
	if err := database.DB.First(&user, session.UserID).Error; err != nil {
		return nil, nil, err
	}

	//check if user is a client
	if user.Role != database.RoleClient {
		return nil, nil, http.ErrNoCookie //returning ErrNoCookie to indicate that the user is not a client
	}

	var client database.Client
	if err := database.DB.Where("user_id = ?", user.ID).First(&client).Error; err != nil {
		return nil, nil, err
	}

	return &user, &client, nil
}

func GetPartnerFromSession(r *http.Request) (*database.User, *database.Partner, error) {
	session, err := GetSessionFromRequest(r)
	if err != nil {
		return nil, nil, err
	}

	var user database.User
	if err := database.DB.First(&user, session.UserID).Error; err != nil {
		return nil, nil, err
	}

	//check if user is a partner
	if user.Role != database.RolePartner {
		return nil, nil, http.ErrNoCookie //returning ErrNoCookie to indicate that the user is not a partner
	}

	var partner database.Partner
	if err := database.DB.Where("user_id = ?", user.ID).First(&partner).Error; err != nil {
		return nil, nil, err
	}

	return &user, &partner, nil
}

func GetAdminFromSession(r *http.Request) (*database.User, *database.Admin, error) {
	session, err := GetSessionFromRequest(r)
	if err != nil {
		return nil, nil, err
	}

	var user database.User
	if err := database.DB.First(&user, session.UserID).Error; err != nil {
		return nil, nil, err
	}

	//check if user is an admin
	if user.Role != database.RoleAdmin {
		return nil, nil, http.ErrNoCookie //returning ErrNoCookie to indicate that the user is not an admin
	}

	var admin database.Admin
	if err := database.DB.Where("user_id = ?", user.ID).First(&admin).Error; err != nil {
		return nil, nil, err
	}

	return &user, &admin, nil
}
