package main

import (
	"cmd1/pkg/models"
	

	
)


type templateData struct {
    Snippet  *models.Snippet
    Snippets []*models.Snippet
	Form interface{}
	Flash string
	IsAuthenticated bool
	CSRFToken string
	Version string
}

