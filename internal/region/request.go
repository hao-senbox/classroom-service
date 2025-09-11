package region

type CreateRegionRequest struct {
	Name string `json:"name" binding:"required"`
	OrganizationID string `json:"organization_id" binding:"required"`
}

type UpdateRegionRequest struct {
	Name           string `json:"name" binding:"required"`
}
