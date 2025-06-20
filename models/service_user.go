package models

/*
	- moved the service model here because it avoids circular dependencies which came up when it was in user.go
	- the struct is used by multiple packages (db, user, etc). so better to move it out in a central place
	- models is place for shared models only, not a junk drawer!
*/

type ServiceUser struct {
	// in GO, field names should start with capital letters to be unmarshaled (decoded from JSON)
	User_Name string `json:"username"`
	// the bits in the back-tics are "struct-tags", this tells json.decode() what to look for
	Password string `json:"password"`
	Location string
	IP_addr  string
}

type ResponseStruct struct {
	Message  string `json:"message"`
	Username string `json:"username,omitempty"`
}
