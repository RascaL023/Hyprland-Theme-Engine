package state

type State struct {
	Version int 	 `json:"version"`

	Theme struct {
		Name 						string `json:"name"`
		Type						string `json:"type"`
	} `json:"theme"`

	Waybar  string `json:"waybar"`
}

