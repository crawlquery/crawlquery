package factory

type Page struct {
	ID          string
	URL         string
	Title       string
	Description string
	HTML        string
}

var ThreePages = map[string]*Page{
	"home1": {
		ID:          "home1",
		URL:         "https://example.com",
		Title:       "Home",
		Description: "Welcome to our official website where we offer the latest updates and information.",
		HTML:        "<html><head><title>Home</title><meta name='description' content='Welcome to our official website where we offer the latest updates and information.'></head><body><h1>Welcome to our official website</h1><p>Here you will find the latest updates and information.</p></body></html>",
	},
	"about1": {
		ID:          "about1",
		URL:         "https://example.com/about",
		Title:       "About",
		Description: "Learn more about our company and our mission.",
		HTML:        "<html><head><title>About</title><meta name='description' content='Learn more about our company and our mission.'></head><body><h1>About Us</h1><p>Our mission is to provide the best services to our customers.</p></body></html>",
	},
	"contact1": {
		ID:          "contact1",
		URL:         "https://example.com/contact",
		Title:       "Contact",
		Description: "Get in touch with us for any inquiries or feedback.",
		HTML:        "<html><head><title>Contact</title><meta name='description' content='Get in touch with us for any inquiries or feedback.'></head><body><h1>Contact Us</h1><p>Feel free to contact us for any questions or feedback.</p></body></html>",
	},
}
