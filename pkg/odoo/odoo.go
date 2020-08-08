package odoo

import (
	odoo "github.com/skilld-labs/go-odoo"
)

var (
	// TravelPackage ...
	TravelPackage = "travel.package"
	// Partner ...
	Partner = "res.partner"
	// Guide ...
	Guide = "guides"
)

// Connection ...
type Connection struct {
	Host string
	DB   string
	User string
	Pass string
}

// Connect ...
func (m Connection) Connect() (*odoo.Client, error) {
	c, err := odoo.NewClient(&odoo.ClientConfig{
		Admin:    m.User,
		Password: m.Pass,
		Database: m.DB,
		URL:      m.Host,
	})

	return c, err
}
