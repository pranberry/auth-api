package user

type ServiceUser struct {
    // in GO, field names should start with capital letters to be unmarshaled (decoded from JSON)
    User_Name string    `json:"username"`
    // the bits in the back-tics are "struct-tags", this tells json.decode() what to look for
    Password string     `json:"password"`
    Location string
    IP_addr string
}
var MasterUserDB = make(map[string]ServiceUser)

type ResponseStruct struct {
    Message string `json:"message"`
    Username string `json:"username,omitempty"`
}