package conf

// Configure of System
type Configure struct {
	Favicon string
	Addr    string
}

// Default Configure
var Default = &Configure{
	// TODO: change favicon
	Favicon: "http://fanyi.bdstatic.com/static/translation/img/favicon/favicon-32x32_ca689c3.png",
	Addr:    "0.0.0.0:4000",
}
