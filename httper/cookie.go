package httper

type CookieCfg struct {
	Name     string `yaml:"name"`
	Age      int    `yaml:"age"`
	Path     string `yaml:"string"`
	Host     string `yaml:"host"`
	Secure   bool   `yaml:"secure"`
	HttpOnly bool   `yaml:"http_only"`
}
