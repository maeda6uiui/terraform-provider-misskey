package misskey

type CreateNoteRequest struct {
	Text           string   `json:"text"`
	Visibility     string   `json:"visibility"`
	VisibleUserIds []string `json:"visibleUserIds"`
}

type CreateNoteResponse struct {
	Id string `json:"id"`
}
