package misskey

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
