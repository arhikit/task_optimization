package new2

import (
	cryptoRand "crypto/rand"
	"github.com/google/uuid"
	"math/rand"
	"runtime"
	"sync"
	"time"
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
	RAND_POOL_SIZE     = 16 * 10000
)

type paramsPool struct {
	poolPos int
	pool    [RAND_POOL_SIZE]byte
}

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
		params := &paramsPool{
			poolPos: RAND_POOL_SIZE,
		}
		go func() {
			for j := 0; j < countLaunches; j++ {
				documents[numGoroutine*countLaunches+j] = generateDocument(params)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	return documents
}

func generateDocument(params *paramsPool) Document {

	countFiles := rand.Intn(10)
	files := make([]File, countFiles, countFiles)
	for i := 0; i < countFiles; i++ {
		files[i] = generateFile(params)
	}

	countFields := rand.Intn(10)
	fields := make([]Field, countFields, countFields)
	for i := 0; i < countFields; i++ {
		fields[i] = generateField(params)
	}

	doc := Document{
		id:          newRandomFromPool(params),
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

func generateField(params *paramsPool) Field {
	return Field{
		id:         newRandomFromPool(params),
		documentID: newRandomFromPool(params),
		name:       FIELD_NAME_STRING,
		value:      FIELD_VALUE_STRING,
	}
}

func generateFile(params *paramsPool) File {
	return File{
		id:             newRandomFromPool(params),
		documentTypeID: newRandomFromPool(params),
		documentID:     newRandomFromPool(params),
		uploadedBy:     newRandomFromPool(params),
		fileName:       FILE_NAME_STRING,
		size:           int64(rand.Int()),
	}
}

func newRandomFromPool(params *paramsPool) uuid.UUID {
	var uuid uuid.UUID
	if params.poolPos == RAND_POOL_SIZE {
		_, err := cryptoRand.Read(params.pool[:])
		if err != nil {
			return uuid
		}
		params.poolPos = 0
	}
	copy(uuid[:], params.pool[params.poolPos:(params.poolPos+16)])
	params.poolPos += 16

	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant is 10
	return uuid
}
