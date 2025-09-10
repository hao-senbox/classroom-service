package classroom

type CreateClassroomRequest struct {
	Name        string  `json:"name"`
	RegionID    *string  `json:"region_id"`
	LocationID  *string `json:"location_id"`
	Description *string `json:"description"`
	Note        *string `json:"note"`
	Icon        *string `json:"icon"`
}
