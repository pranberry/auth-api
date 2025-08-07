package models

/*
	File contains shared models
*/

type ServiceUser struct {
	User_Name string `json:"username"`
	Password  string `json:"password"`
	Location  string
	IP_addr   string
}

type ResponseStruct struct {
	Message  string `json:"message,omitempty"`
	Username string `json:"username,omitempty"`
	AccessToken string `json:"access_token,omitempty"`
	TokenType   string `json:"token_type,omitempty"`
}