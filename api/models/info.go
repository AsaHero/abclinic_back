package models

type Chapter struct {
	GUID string `json:"guid"`
	Name string `json:"name"`
}

type Article struct {
	GUID string `json:"guid"`
	Text string `json:"text"`
	Img  string `json:"img"`
	Side string `json:"side"`
}

type CreateArticleRequest struct {
	ChapterID string `json:"chapter_id"`
	Text      string `json:"text"`
	Img       string `json:"img"`
	Side      string `json:"side"`
}

type UpdateArticleRequest struct {
	Text string `json:"text"`
	Img  string `json:"img"`
	Side string `json:"side"`
}

type CreateChapterRequest struct {
	Name string `json:"name"`
}
