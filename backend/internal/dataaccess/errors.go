package dataaccess
import (
	//"database/sql"
	"errors"
)

var ErrorPostNotFound = errors.New("Post Not Found")
var ErrorUnauthorized = errors.New("Unauthorized to Action")
var ErrorTopicNotFound = errors.New("Topic Not Found")
var ErrorCommentNotFound = errors.New("Comment Not Found")
var ErrorUserNotFound = errors.New("User Not Found")