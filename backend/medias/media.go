package medias

// MediasConfig encapsulates all configuration data
//
// swagger:model
type Configuration struct {
	LastID int     `json:"lastid"`
	Medias []Media `json:"medias"`
}

func NewConfiguration() *Configuration {
	var configuration = new(Configuration)

	configuration.LastID = 0
	configuration.Medias = []Media{}

	return configuration
}

// Media represents a media configuration
//
// swagger:model
type Media struct {
	// the id for this media
	//
	// required: true
	// unique: true
	// min: 1
	ID          int                    `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Rows        int                    `json:"rows"`
	Cols        int                    `json:"cols"`
	Stylesvar   map[string]interface{} `json:"stylesvar"`
	Plugins     []MediaPlugin          `json:"plugins"`
}

func NewMedia() *Media {
	var media = new(Media)

	media.Stylesvar = make(map[string]interface{})
	media.Plugins = []MediaPlugin{}

	return media
}

// MediaPlugin represents a plugin configuration for the media
//
// Properties and configuration for a plugin used in the media
//
// swagger:model
type MediaPlugin struct {
	InstanceId string              `json:"instanceId"`
	Name       string              `json:"name"`
	FrontEnd   MediaPluginFrontEnd `json:"frontend"`
	BackEnd    MediaPluginBackEnd  `json:"backend"`
}

type MediaPluginFrontEnd struct {
	Files   []string               `json:"files"`
	EltName string                 `json:"eltName"`
	X       int                    `json:"x"`
	Y       int                    `json:"y"`
	Rows    int                    `json:"rows"`
	Cols    int                    `json:"cols"`
	Props   map[string]interface{} `json:"props"`
}

type MediaPluginBackEnd struct {
	Ports []int                  `json:"ports"`
	Props map[string]interface{} `json:"props"`
}
