package main

//JSONResponse HTTP Response
type JSONResponse struct {
	METHOD string      `json:"method,omitempty"`
	DATE   string      `json:"date,omitempty"`
	PATH   string      `json:"path,omitempty"`
	STATUS int         `json:"status,omitempty"`
	DATA   interface{} `json:"data,omitempty"`
}

//AuthKey DB
type AuthKey struct {
	ID        int    `json:"id,omitempty"`
	USERID    int    `json:"user_id,omitempty"`
	AUTHKEY   string `json:"authkey,omitempty"`
	TYPE      int    `json:"type,omitempty"`
	TIMESTAMP string `json:"timestamp,omitempty"`
}

//User DB
type User struct {
	ID        int    `json:"id,omitempty"`
	NAME      string `json:"name,omitempty"`
	SURNAME   string `json:"surname,omitempty"`
	EMAIL     string `json:"email,omitempty"`
	PASSWORD  string `json:"password,omitempty"`
	TYPE      int    `json:"type,omitempty"`
	TIMESTAMP string `json:"timestamp,omitempty"`
}
