package models

/*
	- the struct is used by multiple packages (db, user, etc). so better to move it out in a central place
	- models is place for shared models only, not a junk drawer!
*/

type ServiceUser struct {
	User_Name string `json:"username"`
	Password string `json:"password"`
	Location string
	IP_addr  string
}

type ResponseStruct struct {
	Message  string `json:"message"`
	Username string `json:"username,omitempty"`
}
