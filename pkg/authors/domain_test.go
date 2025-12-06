package authors_test

import (
	"testing"

	"github.com/sgaunet/template-api/internal/apperror"
	"github.com/sgaunet/template-api/pkg/authors"
	"github.com/stretchr/testify/assert"
)

func TestNewAuthorName_Valid(t *testing.T) {
	name, err := authors.NewAuthorName("ValidName")
	assert.NoError(t, err)
	assert.Equal(t, "ValidName", name.String())
}

func TestNewAuthorName_TooShort(t *testing.T) {
	_, err := authors.NewAuthorName("abc")
	assert.Error(t, err)
	assert.True(t, apperror.IsValidationError(err))
}

func TestNewAuthorName_TooLong(t *testing.T) {
	_, err := authors.NewAuthorName("ThisNameIsWayTooLongForTheValidation")
	assert.Error(t, err)
	assert.True(t, apperror.IsValidationError(err))
}

func TestNewAuthorName_WithWhitespace(t *testing.T) {
	name, err := authors.NewAuthorName("  ValidName  ")
	assert.NoError(t, err)
	assert.Equal(t, "ValidName", name.String())
}

func TestNewAuthorBio_Valid(t *testing.T) {
	bio, err := authors.NewAuthorBio("This is a valid bio")
	assert.NoError(t, err)
	assert.Equal(t, "This is a valid bio", bio.String())
}

func TestNewAuthorBio_TooLong(t *testing.T) {
	longBio := string(make([]byte, 501))
	_, err := authors.NewAuthorBio(longBio)
	assert.Error(t, err)
	assert.True(t, apperror.IsValidationError(err))
}

func TestCreateAuthorRequest_Validate_Valid(t *testing.T) {
	req := authors.CreateAuthorRequest{
		Name: "ValidName",
		Bio:  "A valid bio",
	}
	err := req.Validate()
	assert.NoError(t, err)
}

func TestCreateAuthorRequest_Validate_InvalidName(t *testing.T) {
	req := authors.CreateAuthorRequest{
		Name: "abc",
		Bio:  "A valid bio",
	}
	err := req.Validate()
	assert.Error(t, err)
	assert.True(t, apperror.IsValidationError(err))
}

func TestCreateAuthorRequest_ToAuthor(t *testing.T) {
	req := authors.CreateAuthorRequest{
		Name: "ValidName",
		Bio:  "A valid bio",
	}
	author, err := req.ToAuthor()
	assert.NoError(t, err)
	assert.Equal(t, "ValidName", author.Name)
	assert.Equal(t, "A valid bio", author.Bio)
	assert.Equal(t, int64(0), author.ID)
}

func TestAuthor_ToResponse(t *testing.T) {
	author := &authors.Author{
		ID:   123,
		Name: "Test Author",
		Bio:  "Test bio",
	}
	response := author.ToResponse()
	assert.Equal(t, int64(123), response.ID)
	assert.Equal(t, "Test Author", response.Name)
	assert.Equal(t, "Test bio", response.Bio)
}
