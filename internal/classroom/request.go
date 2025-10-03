package classroom

type CreateClassroomRequest struct {
	Name        string  `json:"name"`
	LanguageID  uint    `json:"language_id"`
	RegionID    *string `json:"region_id"`
	LocationID  *string `json:"location_id"`
	Description *string `json:"description"`
	Note        *string `json:"note"`
	Icon        *string `json:"icon"`
}

type UpdateClassroomRequest struct {
	Name        string  `json:"name"`
	LanguageID  uint    `json:"language_id"`
	RegionID    *string `json:"region_id"`
	LocationID  *string `json:"location_id"`
	Description *string `json:"description"`
	Note        *string `json:"note"`
	Icon        *string `json:"icon"`
}

type CreateAssignmentByTemplateRequest struct {
	ClassroomID string `json:"classroom_id"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
}