package dao

import (
	"book-service/src/models"
	"github.com/astaxie/beego/orm"
	"time"
)

func CreateBook(book *models.Book) (int64, error) {
	book.CreationTime = time.Now()
	book.UpdateTime = time.Now()
	bookID, err := GetOrmer().Insert(book)
	if err != nil {
		return 0, err
	}

	return bookID, nil
}

func UpdateBook(book *models.Book) error {
	book.UpdateTime = time.Now()
	_, err := GetOrmer().Update(book)
	return err
}

// GetBook specified by ID
func GetBook(id int64) (*models.Book, error) {
	book := &models.Book{ID: id}
	if err := GetOrmer().Read(book); err != nil {
		if err == orm.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return book, nil
}

// ListBooks list all books
func ListBooks(query *models.BookQuery) ([]*models.Book, error) {
	var err error
	var books []*models.Book

	qs := GetOrmer().QueryTable("book")

	if len(query.Name) > 0 {
		qs = qs.Filter("Name__icontains", query.Name)
	}

	if query.OwnerID != 0 {
		qs = qs.Filter("OwnerID", query.OwnerID)
	}

	qs = qs.Filter("Deleted", false)

	_, err = qs.All(&books)

	return books, err
}

// DeleteBook ...
func DeleteBook(id int64) error {
	book, err := GetBook(id)
	if err != nil {
		return err
	}
	book.UpdateTime = time.Now()
	book.Deleted = true
	_, err = GetOrmer().Update(book, "Name", "UpdateTime", "Deleted")
	return err
}
