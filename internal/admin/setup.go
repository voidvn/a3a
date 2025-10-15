package admin

import (
	"net/http"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/qor/admin"
)

// SetupAdmin initializes the admin interface
func SetupAdmin(databaseURL string) (http.Handler, error) {
	// Initialize GORM DB connection
	db, err := gorm.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	// Initialize the admin
	Admin := admin.New(&admin.AdminConfig{
		SiteName: "S4S Admin",
		DB:       db,
	})

	// Add models to admin interface
	// User model
	user := Admin.AddResource(&User{})
	user.IndexAttrs("ID", "Email", "FirstName", "LastName", "Role", "IsActive", "CreatedAt")
	user.EditAttrs("Email", "FirstName", "LastName", "Role", "IsActive")

	// Product model
	product := Admin.AddResource(&Product{})
	product.IndexAttrs("ID", "Name", "Price", "StockQuantity", "Category", "IsActive")
	product.EditAttrs("Name", "Description", "Price", "StockQuantity", "Category", "ImageURL", "IsActive")

	// Order model
	order := Admin.AddResource(&Order{})
	order.IndexAttrs("ID", "User", "TotalAmount", "Status", "CreatedAt")
	order.EditAttrs("User", "Status", "ShippingAddress", "PaymentMethod")

	// Mount admin to the router
	mux := http.NewServeMux()
	Admin.MountTo("/admin", mux)

	return mux, nil
}

// User represents the user model for admin
// Note: This should match your actual user model
// You may need to adjust the fields based on your actual model
type User struct {
	gorm.Model
	Email     string `gorm:"unique;not null"`
	Password  string `gorm:"-"` // This field is not stored in the database
	FirstName string
	LastName  string
	Role      string `gorm:"default:'user'"`
	IsActive  bool   `gorm:"default:true"`
}

// Product represents the product model for admin
type Product struct {
	gorm.Model
	Name          string `gorm:"not null"`
	Description   string
	Price         float64 `gorm:"not null"`
	StockQuantity int     `gorm:"not null;default:0"`
	Category      string
	ImageURL      string
	IsActive      bool `gorm:"default:true"`
}

// Order represents the order model for admin
type Order struct {
	gorm.Model
	UserID          uint    `gorm:"not null"`
	User            User    `gorm:"foreignkey:UserID"`
	TotalAmount     float64 `gorm:"not null"`
	Status          string  `gorm:"default:'pending'"`
	ShippingAddress string
	PaymentMethod   string
}
