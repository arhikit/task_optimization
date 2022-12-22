package new1

import (
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/google/uuid"
)

const (
	DOCUMENTS_COUNT    = 1000000
	NEW_COMMENT_STRING = "newComment"
	COMMENT_STRING     = "comment"
	TYPE_STRING        = "type"
	DESCRIPTION_STRING = "description"
	FIELD_NAME_STRING  = "FieldName"
	FIELD_VALUE_STRING = "FieldValue"
	FILE_NAME_STRING   = "filename"
)

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
	for i := 0; i < len(documents); i++ {
		transformOneDocument(documents[i])
	}
	return documents
}

func transformOneDocument(doc Document) {
	doc.isApproved = false
	doc.isRequired = true
	doc.comment = NEW_COMMENT_STRING
}

func getDocumentsFromIntegration() []Document {

	wg := &sync.WaitGroup{}
	countCPU := runtime.NumCPU()
	countLaunches := int(DOCUMENTS_COUNT / countCPU)
	documents := make([]Document, countCPU*countLaunches, countCPU*countLaunches)

	for i := 0; i < countCPU; i++ {
		wg.Add(1)
		numGoroutine := i
		go func() {
			for j := 0; j < countLaunches; j++ {
				documents[numGoroutine*countLaunches+j] = generateDocument()
			}
			wg.Done()
		}()
	}

	wg.Wait()
	return documents
}

func generateDocument() Document {

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
		docType:     TYPE_STRING,
		approvedAt:  time.Now(),
		fields:      fields,
		files:       files,
		description: DESCRIPTION_STRING,
		comment:     COMMENT_STRING,
		isApproved:  true,
		isRequired:  false,
	}

	return doc
}

func generateField() Field {
	return Field{
		id:         uuid.New(),
		documentID: uuid.New(),
		name:       FIELD_NAME_STRING,
		value:      FIELD_VALUE_STRING,
	}
}

func generateFile() File {
	return File{
		id:             uuid.New(),
		documentTypeID: uuid.New(),
		documentID:     uuid.New(),
		uploadedBy:     uuid.New(),
		fileName:       FILE_NAME_STRING,
		size:           int64(rand.Int()),
	}
}
