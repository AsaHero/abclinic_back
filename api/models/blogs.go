package models

type Categories struct {
	GUID        string `json:"guid"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Img         string `json:"img"`
}

type Authors struct {
	GUID string `json:"guid"`
	Name string `json:"name"`
	Img  string `json:"img"`
}

type Contents struct {
	URL string `json:"url"`
}

type Publications struct {
	GUID       string     `json:"guid"`
	CategoryID string     `json:"category_id"`
	Author     Authors    `json:"author"`
	Title      string     `json:"title"`
	Text       string     `json:"text"`
	Type       string     `json:"type"`
	Video      string     `json:"video"`
	Img        []Contents `json:"img"`
}

type CreatePublicationRequest struct {
	AuthorID string     `json:"author_id"`
	Title    string     `json:"title"`
	Text     string     `json:"text"`
	Type     string     `json:"type"`
	Video    string     `json:"video"`
	Img      []Contents `json:"img"`
}

type UpdatePublicationRequest struct {
	Title string     `json:"title"`
	Text  string     `json:"text"`
	Type  string     `json:"type"`
	Video string     `json:"video"`
	Img   []Contents `json:"img"`
}

type CreateCategoryRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Img         string `json:"img"`
}

type CreateAuthorRequest struct {
	Name string `json:"name"`
	Img  string `json:"img"`
}
