package ui 

import (
	"embed"
)

// The bottom is a special comment directive
// When compiled, the comment directive tells Go to store 
// the files from ui/html, and ui/static in embed.FS.
//go:embed "html" "static"
var Files embed.FS
