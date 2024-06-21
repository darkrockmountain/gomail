package common

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"strings"

	"github.com/darkrockmountain/gomail/sanitizer"
)

// Attachment represents an email attachment with its filename and content.
// Use this struct to specify files to be attached to the email.
type Attachment struct {
	filename string // The name of the file.
	content  []byte // The content of the file.
}

// NewAttachment creates a new Attachment instance with the specified filename and content.
// It initializes the private fields of the Attachment struct with the provided values.
//
// Example:
//
//	content := []byte("file content")
//	attachment := NewAttachment("document.pdf", content)
//	fmt.Println(attachment.GetFilename()) // Output: document.pdf
//	fmt.Println(string(attachment.GetContent())) // Output: file content
func NewAttachment(filename string, content []byte) *Attachment {
	return &Attachment{
		filename: filename,
		content:  content,
	}
}

// NewAttachmentFromFile creates a new Attachment instance from the specified file path.
// It reads the content of the file and initializes the private fields of the Attachment struct.
//
// Example:
//
//	attachment, err := NewAttachmentFromFile("path/to/document.pdf")
//	if err != nil {
//	    fmt.Println("Error creating attachment:", err)
//	    return
//	}
//	fmt.Println(attachment.GetFilename()) // Output: document.pdf
//	fmt.Println(string(attachment.GetContent())) // Output: (file content)
func NewAttachmentFromFile(filePath string) (*Attachment, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	filename := extractFilename(filePath)
	return NewAttachment(
		filename,
		content,
	), nil
}

// extractFilename extracts the filename from the file path.
// This is a helper function to get the filename from a given file path.
func extractFilename(filePath string) string {
	// Implement this function based on your needs, for simplicity using base method
	parts := strings.Split(filePath, "/")
	return parts[len(parts)-1]
}

// SetFilename sets the filename of the attachment.
// It assigns the provided filename to the private filename field.
//
// Example:
//
//	var attachment Attachment
//	attachment.SetFilename("document.pdf")
//	fmt.Println(attachment.GetFilename()) // Output: document.pdf
func (a *Attachment) SetFilename(filename string) {
	a.filename = filename
}

// GetFilename safely returns the filename of the attachment.
// It uses the default text sanitizer to escape special characters and trim whitespace.
// If the attachment is nil, it returns an empty string.
//
// Returns:
//   - string: The sanitized filename.
func (a *Attachment) GetFilename() string {
	if a == nil {
		return "nil_attachment"
	}
	return sanitizer.DefaultTextSanitizer().Sanitize(a.filename)
}

// GetBase64StringContent returns the content of the attachment as a base64-encoded string.
// If the attachment is nil, it returns an empty string.
//
// Returns:
//   - string: The base64-encoded content of the attachment as a string.
//     Returns an empty string if the attachment is nil.
func (a *Attachment) GetBase64StringContent() string {
	if a == nil {
		return ""
	}
	return string(a.GetBase64Content())
}

// SetContent sets the content of the attachment.
// It assigns the provided content to the private content field.
//
// Example:
//
//	var attachment Attachment
//	content := []byte("file content")
//	attachment.SetContent(content)
//	fmt.Println(string(attachment.GetContent())) // Output: file content
func (a *Attachment) SetContent(content []byte) {
	a.content = content
}

// GetBase64Content returns the content of the attachment as a base64-encoded byte slice.
// If the attachment is nil, it returns an empty byte slice.
//
// Returns:
//   - []byte: The base64-encoded content of the attachment as a byte slice.
//     Returns an empty byte slice if the attachment is nil.
func (a *Attachment) GetBase64Content() []byte {
	if a == nil || len(a.content) == 0 {
		return []byte{}
	}
	buf := make([]byte, base64.StdEncoding.EncodedLen(len(a.content)))
	base64.StdEncoding.Encode(buf, a.content)
	return buf
}

// GetRawContent returns the content of the attachment as its raw byte slice.
// If the attachment is nil, it returns an empty byte slice.
//
// Returns:
//   - []byte: The content of the attachment as a byte slice.
//     Returns an empty byte slice if the attachment is nil.
func (a *Attachment) GetRawContent() []byte {
	if a == nil || len(a.content) == 0 {
		return []byte{}
	}
	return a.content
}

// jsonAttachment represents the JSON structure for an email attachment.
type jsonAttachment struct {
	Filename string `json:"filename"`
	Content  string `json:"content"` // Content will be base64 encoded
}

// MarshalJSON custom marshaler for Attachment
// This method converts the Attachment struct into a JSON representation.
// It creates an anonymous struct with exported fields and JSON tags,
// copies the values from the private fields, and then marshals it to JSON.
//
// Example:
//
//	attachment := Attachment{
//	    filename: "file.txt",
//	    content:  []byte("file content"),
//	}
//	jsonData, err := json.Marshal(attachment)
//	if err != nil {
//	    fmt.Println("Error marshaling to JSON:", err)
//	    return
//	}
//	fmt.Println("JSON output:", string(jsonData))
func (a Attachment) MarshalJSON() ([]byte, error) {
	return json.Marshal(&jsonAttachment{
		Filename: a.filename,
		Content:  base64.StdEncoding.EncodeToString(a.content), // Encode content to base64
	})
}

// UnmarshalJSON custom unmarshaler for Attachment
// This method converts a JSON representation into an Attachment struct.
// It creates an anonymous struct with exported fields and JSON tags,
// unmarshals the JSON data into this struct, and then copies the values
// to the private fields of the Attachment struct.
//
// Example:
//
//	jsonData := `{
//	    "filename": "file.txt",
//	    "content": "ZmlsZSBjb250ZW50" // base64 encoded "file content"
//	}`
//	var attachment Attachment
//	err := json.Unmarshal([]byte(jsonData), &attachment)
//	if err != nil {
//	    fmt.Println("Error unmarshaling from JSON:", err)
//	    return
//	}
//	fmt.Printf("Unmarshaled Attachment: %+v\n", attachment)
func (a *Attachment) UnmarshalJSON(data []byte) error {
	aux := &jsonAttachment{}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	a.filename = aux.Filename
	content, err := base64.StdEncoding.DecodeString(aux.Content) // Decode content from base64
	if err != nil {
		return err
	}
	a.content = content

	return nil
}
