package classroom

import (
	"classroom-service/internal/language"
	"classroom-service/pkg/constants"
)

func BuildDepartmentMessagesUpdate(classroomID string, req CreateClassroomRequest) language.UploadMessageLanguagesRequest {
	return language.UploadMessageLanguagesRequest{
		MessageLanguages: []language.UploadMessageRequest{
			{
				TypeID:     classroomID,
				Type:       "classroom",
				Key:        string(constants.ClassroomNoteKey),
				Value:      *req.Note,
				LanguageID: req.LanguageID,
			},
			{
				TypeID:     classroomID,
				Type:       "classroom",
				Key:        string(constants.ClassroomNameKey),
				Value:      req.Name,
				LanguageID: req.LanguageID,
			},
			{
				TypeID:     classroomID,
				Type:       "classroom",
				Key:        string(constants.ClassroomDescKey),
				Value:      *req.Description,
				LanguageID: req.LanguageID,
			},
		},
	}
}
