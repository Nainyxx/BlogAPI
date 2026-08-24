package Gin

import "BlogAPI/service"

type BlogHandler struct {
	BlogService *service.BlogService
}
