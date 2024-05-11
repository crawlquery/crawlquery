package factory

import (
	"crawlquery/pkg/domain"
)

var HomePage = domain.Page{
	ID:              "1",
	URL:             "https://example.com",
	Title:           "Home",
	Content:         "<html><body>Welcome to our homepage.</body></html>",
	MetaDescription: "Welcome to our official website where we offer the latest updates and information.",
}

func TenPages() []domain.Page {
	return []domain.Page{
		HomePage,
		{ID: "2", URL: "https://example.com/about", Title: "About Us", Content: "<html><body>Learn more about our company's history and mission.</body></html>", MetaDescription: "Discover more about our company and what drives us."},
		{ID: "3", URL: "https://example.com/contact", Title: "Contact Us", Content: "<html><body>Contact us via email or phone.</body></html>", MetaDescription: "Get in touch with us for more information and support."},
		{ID: "4", URL: "https://example.com/blog", Title: "Blog", Content: "<html><body>Read our latest blog posts and updates.</body></html>", MetaDescription: "Stay updated with our latest blog posts and news articles."},
		{ID: "5", URL: "https://example.com/services", Title: "Our Services", Content: "<html><body>Explore our range of services.</body></html>", MetaDescription: "Browse our comprehensive range of services and find out how we can help you."},
		{ID: "6", URL: "https://example.com/products", Title: "Our Products", Content: "<html><body>Discover our products and find something you love.</body></html>", MetaDescription: "Explore our wide range of products."},
		{ID: "7", URL: "https://example.com/team", Title: "Meet Our Team", Content: "<html><body>Meet the people behind our company.</body></html>", MetaDescription: "Learn more about our team and their professional backgrounds."},
		{ID: "8", URL: "https://example.com/careers", Title: "Careers", Content: "<html><body>Join our team and help us grow.</body></html>", MetaDescription: "Explore career opportunities and become part of our team."},
		{ID: "9", URL: "https://example.com/privacy", Title: "Privacy Policy", Content: "<html><body>Read our privacy policy.</body></html>", MetaDescription: "Understand how we collect, use, and protect your data."},
		{ID: "10", URL: "https://example.com/terms", Title: "Terms and Conditions", Content: "<html><body>Review our terms and conditions.</body></html>", MetaDescription: "Read the terms and conditions of using our website and services."},
	}
}
