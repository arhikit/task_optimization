package example

import (
	cryptoRand "crypto/rand"
	"math/rand"
	"sync"
)

const (
	DOCUMENTS_COUNT    = 10
	NEW_COMMENT_STRING = "newComment"
	COMMENT_STRING     = "comment"
	TYPE_STRING        = "type"
	DESCRIPTION_STRING = "description"
	FIELD_NAME_STRING  = "FieldName"
	FIELD_VALUE_STRING = "FieldValue"
	FILE_NAME_STRING   = "filename"
	randPoolSize       = 16 * 10000
)

var (
	poolMu  sync.Mutex
	poolPos = randPoolSize     // protected with poolMu
	pool    [randPoolSize]byte // protected with poolMu
)

type File struct {
	id             [16]byte
	documentTypeID [16]byte
	documentID     [16]byte
	uploadedBy     [16]byte
	fileName       string
	size           int64
}

func generateFiles_arr() [10]File {

	var arrFiles [10]File
	for i := 0; i < 10; i++ {
		arrFiles[i] = generateFile()
	}
	return arrFiles
}

func generateFiles_slice() []File {

	var arrFiles [10]File
	for i := 0; i < 10; i++ {
		arrFiles[i] = generateFile()
	}
	return arrFiles[:]
}

func generateFile() File {
	return File{
		id:             newRandomFromPool(),
		documentTypeID: newRandomFromPool(),
		documentID:     newRandomFromPool(),
		uploadedBy:     newRandomFromPool(),
		fileName:       FILE_NAME_STRING,
		size:           int64(rand.Int()),
	}
}

func newRandomFromPool() [16]byte {
	var uuid [16]byte
	if poolPos == randPoolSize {
		_, err := cryptoRand.Read(pool[:])
		if err != nil {
			return uuid
		}
		poolPos = 0
	}
	copy(uuid[:], pool[poolPos:(poolPos+16)])
	poolPos += 16

	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant is 10
	return uuid
}
