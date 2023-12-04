package repository

import (
	"notes/model"

	"gorm.io/gorm"
)

type notesRepositoryImpl struct {
	q *gorm.DB
}

func NewNoteRepository(db *gorm.DB) NoteRepository {
	return &notesRepositoryImpl{
		q: db,
	}
}

func (repos *notesRepositoryImpl) CreateNote(note *model.Note) error {
	exec := repos.q.Create(note)

	if exec.Error != nil {
		return exec.Error
	}

	return nil
}

func (repos *notesRepositoryImpl) FindNote(id string) *model.Note {
	var data model.Note

	exec := repos.q.Where("id = ?", id).First(&data)

	if exec.Error != nil {
		return nil
	}

	return &data

}
