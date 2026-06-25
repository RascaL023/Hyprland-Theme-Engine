package hypr

type Hypr struct {
	GapsIn					int
	GapsOut					int

	BorderSize				int
	BorderActiveDegree		int
	BorderActivePrimary		string
	BorderActiveSecondary	string
	BorderInActivePrimary	string

	BlurOption				bool
	BlurSize				int
	BlurPasses				int

	Rounding				int

	ShadowOption			bool
	ShadowRange				int
	ShadowRenderPower		int
	ShadowColor				string

	WindowRuleFootOpacity	float64
}
