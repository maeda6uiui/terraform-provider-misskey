package misskey

type ErrorInfo struct {
	Message string `json:"message"`
	Code    string `json:"code"`
	Id      string `json:"id"`
}

type ErrorResponse struct {
	Error ErrorInfo `json:"error"`
}

type CreateNoteRequest struct {
	Text           string   `json:"text"`
	Visibility     string   `json:"visibility"`
	VisibleUserIds []string `json:"visibleUserIds"`
}

type Note struct {
	Id             string   `json:"id"`
	Text           string   `json:"text"`
	Visibility     string   `json:"visibility"`
	VisibleUserIds []string `json:"visibleUserIds"`
}

type CreateNoteResponse struct {
	CreatedNote Note `json:"createdNote"`
}

type DeleteNoteRequest struct {
	NoteId string `json:"noteId"`
}

type ShowNoteRequest struct {
	NoteId string `json:"noteId"`
}
