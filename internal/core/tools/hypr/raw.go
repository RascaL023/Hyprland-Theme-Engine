package hypr

type Raw struct {
	Gaps struct {
		In				int 	`json:"in"`
		Out 			int		`json:"out"`
	} `json:"gaps"`

	Border struct {
		Size 			int 	`json:"size"`
		Active struct {
			Degree		int		`json:"degree"`
			Primary		string	`json:"primary"`
			Secondary	string	`json:"secondary"`
		} `json:"active"`
		Inactive struct {
			Primary 	string	`json:"primary"`
		} `json:"inactive"`
	} `json:"border"`

	Blur struct {
		Option			bool	`json:"option"`
		Size			int		`json:"size"`
		Passes			int		`json:"passes"`
	} `json:"blur"`

	Rounding			int		`json:"rounding"`

	Shadow struct {
		Option			bool	`json:"option"`
		Range			int		`json:"range"`
		RenderPower		int		`json:"renderPower"`
		Color			string	`json:"color"`
	} `json:"shadow"`


	WindowRule struct {
		Foot struct {
			Opacity		float64	`json:"opacity"`
		} `json:"foot"`
	} `json:"windowRule"`
}
