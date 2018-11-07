package model

// Author contains the author info
type Author struct {
	Name    string `json:"name,omitempty"`
	URL     string `json:"url,omitempty"`
	IconURL string `json:"icon_url,omitempty"`
}

// Field is a single field of values
type Field struct {
	Name   string `json:"name,omitempty"`
	Value  string `json:"value,omitempty"`
	Inline bool   `json:"inline,omitempty"`
}

// Thumbnail is the thumbnail of content image
type Thumbnail struct {
	URL string `json:"url,omitempty"`
}

// Image is the content image
type Image struct {
	URL string `json:"url,omitempty"`
}

// Footer is the footer of the hook object
type Footer struct {
	Text    string `json:"text,omitempty"`
	IconURL string `json:"icon_url,omitempty"`
}

// Embed is a single Embed block
type Embed struct {
	Author      Author    `json:"author"`
	Title       string    `json:"title"`
	URL         string    `json:"url,omitempty"`
	Description string    `json:"description"`
	Color       int       `json:"color,omitempty"`
	Fields      []Field   `json:"fields,omitempty"`
	Thumbnail   Thumbnail `json:"thumbnail,omitempty"`
	Image       Image     `json:"image,omitempty"`
	Footer      Footer    `json:"footer,omitempty"`
}

// Webhook is a the webhook object
type Webhook struct {
	Username  string  `json:"username,omitempty"`
	AvatarURL string  `json:"avatar_url,omitempty"`
	Content   string  `json:"content,omitempty"`
	Embeds    []Embed `json:"embeds"`
}
