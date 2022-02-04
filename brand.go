package easypost

// Brand contains data for a user brand
type Brand struct {
	BackgroundColor string `json:"background_color,omitempty"`
	Color           string `json:"color,omitempty"`
	Logo            string `json:"logo,omitempty"`
	LogoHref        string `json:"logo_href,omitempty"`
	Ad              string `json:"ad,omitempty"`
	AdHref          string `json:"ad_href,omitempty"`
	Name            string `json:"name,omitempty"`
	UserID          string `json:"user_id,omitempty"`
	Theme           string `json:"theme,omitempty"`
}
