package mm

type Member struct {
	Name       string   `json:"name"`
	PictureUrl string   `json:"profile_url"`
	Nicknames  []string `json:'nicknames'`
	BlogUrl    string   `json:"blog_url"`
}
