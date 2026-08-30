package misskey

type CreateNoteRequest struct {
	Text           string   `json:"text"`
	Visibility     string   `json:"visibility"`
	VisibleUserIds []string `json:"visibleUserIds"`
}

type Note struct {
	Id string `json:"id"`
}

type CreateNoteResponse struct {
	CreatedNote Note `json:"createdNote"`
}
