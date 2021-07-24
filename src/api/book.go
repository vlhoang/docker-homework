package api

import (
	"book-service/src/dao"
	"book-service/src/models"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

// BookAPI handles requests for book management
type BookAPI struct {
	book *models.Book
	BaseAPI
}

// Prepare ...
func (b *BookAPI) Prepare() {
	method := b.Ctx.Request.Method

	if method == http.MethodPut || method == http.MethodDelete {
		id, err := b.GetInt64FromPath(":id")
		if err != nil || id <= 0 {
			b.SendBadRequestError(errors.New("invalid book ID"))
			return
		}

		book, err := dao.GetBook(id)
		if err != nil {
			b.SendInternalServerError(fmt.Errorf("failed to get book %d: %v", id, err))
			return
		}

		if book == nil {
			b.SendNotFoundError(fmt.Errorf("book %d not found", id))
			return
		}

		b.book = book
	}
}

// Post creates a book
func (b *BookAPI) Post() {
	book := &models.Book{}
	isValid, err := b.DecodeJSONReqAndValidate(book)
	if !isValid {
		b.SendBadRequestError(err)
		return
	}

	// check duplicate book
	books, err := dao.ListBooks(&models.BookQuery{
		Name: book.Name,
	})
	if err != nil {
		b.SendInternalServerError(fmt.Errorf("failed to list books: %v", err))
		return
	}
	if len(books) > 0 {
		b.SendConflictError(errors.New("conflict book"))
		return
	}

	id, err := dao.CreateBook(book)
	if err != nil {
		b.SendInternalServerError(fmt.Errorf("failed to create book: %v", err))
		return
	}

	b.Redirect(http.StatusCreated, strconv.FormatInt(id, 10))
}

// Get the book specified by ID
func (b *BookAPI) Get() {
	id, err := b.GetInt64FromPath(":id")
	if err != nil || id <= 0 {
		b.SendBadRequestError(fmt.Errorf("invalid book ID: %s", b.GetStringFromPath(":id")))
		return
	}

	book, err := dao.GetBook(id)
	if err != nil {
		b.SendInternalServerError(fmt.Errorf("failed to get book %d: %v", id, err))
		return
	}

	if book == nil || book.Deleted {
		b.SendNotFoundError(fmt.Errorf("book %d not found", id))
		return
	}

	b.Data["json"] = book
	b.ServeJSON()
}

// List books according to the query strings
func (b *BookAPI) List() {
	ownerId, _ := b.GetInt64("owner_id")
	query := &models.BookQuery{
		Name:    b.GetString("name"),
		OwnerID: ownerId,
	}

	books, err := dao.ListBooks(query)
	if err != nil {
		b.SendInternalServerError(fmt.Errorf("failed to list books: %v", err))
		return
	}

	b.Data["json"] = books
	b.ServeJSON()
}

// Put updates the label
func (b *BookAPI) Put() {
	book := &models.Book{}
	if err := b.DecodeJSONReq(book); err != nil {
		b.SendBadRequestError(err)
		return
	}

	oldName := b.book.Name

	// only name and description can be changed
	b.book.Name = book.Name
	b.book.Description = book.Description

	if b.book.Name != oldName {
		// check duplicate book
		books, err := dao.ListBooks(&models.BookQuery{
			Name: book.Name,
		})
		if err != nil {
			b.SendInternalServerError(fmt.Errorf("failed to list books: %v", err))
			return
		}
		if len(books) > 0 {
			b.SendConflictError(errors.New("conflict book"))
			return
		}
	}

	if err := dao.UpdateBook(b.book); err != nil {
		b.SendInternalServerError(fmt.Errorf("failed to update book %d: %v", b.book.ID, err))
		return
	}

	b.Redirect(http.StatusOK, strconv.FormatInt(b.book.ID, 10))
}

// Delete the book
func (b *BookAPI) Delete() {
	id := b.book.ID
	if err := dao.DeleteBook(id); err != nil {
		b.SendInternalServerError(fmt.Errorf("failed to delete book %d: %v", id, err))
		return
	}

	b.Redirect(http.StatusOK, strconv.FormatInt(id, 10))
}
