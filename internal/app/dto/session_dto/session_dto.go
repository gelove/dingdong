package session_dto

type Data struct {
	Scheme  string  `json:"scheme"`
	Host    string  `json:"host"`
	Path    string  `json:"path"`
	Query   string  `json:"query"`
	Request Request `json:"request"`
}

type Request struct {
	Sizes           any    `json:"sizes"`
	MimeType        any    `json:"mimeType"`
	Charset         any    `json:"charset"`
	ContentEncoding any    `json:"contentEncoding"`
	Header          Header `json:"header"`
}

type Header struct {
	Headers []Headers `json:"headers"`
}

type Headers struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
