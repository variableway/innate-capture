package idea

// CreateInput is the write request model for an idea.
type CreateInput struct {
	Title       string
	Description string
	Context     string
	Source      string
}

// Entry is a single inbox markdown file.
type Entry struct {
	Path     string
	Filename string
	Title    string
	Date     string
}
