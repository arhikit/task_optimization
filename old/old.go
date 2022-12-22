package old

import (
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
)

const DOCUMENTS_COUNT = 1000000

type Document struct {
	id          uuid.UUID
	docType     string
	comment     string
	isApproved  bool
	isRequired  bool
	approvedAt  time.Time
	fields      []Field
	files       []File
	description string
}

type Field struct {
	id         uuid.UUID
	documentID uuid.UUID
	name       string
	value      string
}

type File struct {
	id             uuid.UUID
	documentTypeID uuid.UUID
	documentID     uuid.UUID
	uploadedBy     uuid.UUID
	fileName       string
	size           int64
}

func GetTransformDocuments() {
	// слайс документов, полученных из другого сервиса
	documents := getDocumentsFromIntegration()

	transformDocuments(documents)
}

func transformDocuments(documents []Document) []Document {
	result := make([]Document, 0)

	for i := range documents {
		result = append(result, transformOneDocument(documents[i]))
	}

	return result
}

func transformOneDocument(doc Document) Document {
	doc.isApproved = false
	doc.isRequired = true
	doc.comment = "newComment"

	return doc
}

func getDocumentsFromIntegration() []Document {
	c := make(chan Document, 10)
	wg := &sync.WaitGroup{}
	result := make([]Document, 0, DOCUMENTS_COUNT)

	go func() {
		for {
			doc, open := <-c
			if !open {
				break
			}

			result = append(result, doc)
		}
	}()

	for i := 0; i <= DOCUMENTS_COUNT; i++ {
		wg.Add(1)
		go generateDocument(wg, c)
	}

	wg.Wait()

	return result
}

func generateDocument(wg *sync.WaitGroup, c chan Document) {
	files := make([]File, 0)
	fields := make([]Field, 0)

	for i := 0; i < rand.Intn(10); i++ {
		files = append(files, generateFile())
	}

	for i := 0; i < rand.Intn(10); i++ {
		fields = append(fields, generateField())
	}

	doc := Document{
		id:          uuid.New(),
		docType:     "type",
		comment:     "comment",
		isApproved:  true,
		isRequired:  false,
		approvedAt:  time.Now(),
		fields:      fields,
		files:       files,
		description: "description",
	}

	c <- doc
	wg.Done()
}

func generateField() Field {
	return Field{
		id:         uuid.New(),
		documentID: uuid.New(),
		name:       "FieldName",
		value:      "FieldValue",
	}
}

func generateFile() File {
	return File{
		id:             uuid.New(),
		documentTypeID: uuid.New(),
		documentID:     uuid.New(),
		uploadedBy:     uuid.New(),
		fileName:       "filename",
		size:           int64(rand.Int()),
	}
}
