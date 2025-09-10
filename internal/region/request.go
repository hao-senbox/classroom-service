package region

type CreateRegionRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdateRegionRequest struct {
	Name string `json:"name" binding:"required"`
}