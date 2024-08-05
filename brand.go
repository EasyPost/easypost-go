package easypost

// Brand contains data for a user brand
type Brand struct {
	ID              string `json:"id,omitempty" url:"id,omitempty"`
	BackgroundColor string `json:"background_color,omitempty" url:"background_color,omitempty"`
	Color           string `json:"color,omitempty" url:"color,omitempty"`
	Logo            string `json:"logo,omitempty" url:"logo,omitempty"`
	LogoHref        string `json:"logo_href,omitempty" url:"logo_href,omitempty"`
	Ad              string `json:"ad,omitempty" url:"ad,omitempty"`
	AdHref          string `json:"ad_href,omitempty" url:"ad_href,omitempty"`
	Name            string `json:"name,omitempty" url:"name,omitempty"`
	UserID          string `json:"user_id,omitempty" url:"user_id,omitempty"`
	Theme           string `json:"theme,omitempty" url:"theme,omitempty"`
}
