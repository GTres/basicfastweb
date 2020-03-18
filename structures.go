package main

type Server struct {
	Host				string		`json:"host,omitempty"`
	Server_Changed		bool		`json:"servers_changed,omitempty"`
	Ssl_Grade			string		`json:"ssl_grade,omitempty"`
	Previous_Ssl_Grade	string		`json:"previous_ssl_grade,omitempty"`
	Logo				string		`json:"logo,omitempty"`
	Title				string		`json:"title,omitempty"`
	Is_Down				bool		`json:"is_down,omitempty"`
	EndPoints			[]EndPoint	`json:"endpoints"`
}

type EndPoint struct {
	Address				string		`json:"ipAddress"`
	Ssl_Grade			string		`json:"grade"`
	Country				string		`json:"country"`
	Owner				string		`json:"owner"`
}