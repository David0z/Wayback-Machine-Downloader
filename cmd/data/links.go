package data

type Link struct {
	Urlkey     string `json:"urlkey"`
	Timestamp  string `json:"timestamp"`
	Original   string `json:"original"`
	Mimetype   string `json:"mimetype"`
	Statuscode string `json:"statuscode"`
	Downloaded bool   `json:"downloaded"`
	WebsiteURL string `json:"websiteurl"`
}
