package models

type Request struct {
	Url string `json:"url" validate:"required"` //не забыть добавить валидацию в целом
}

type Response struct {
	UrlOriginal  string `json:"url_original"`
	Shortcode    string `json:"shortcode"`
	UrlShortened string `json:"url_shortened"`
}
