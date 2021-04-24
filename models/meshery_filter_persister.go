package models

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshkit/database"
)

// MesheryFilterPersister is the persister for persisting
// wasm filters on the database
type MesheryFilterPersister struct {
	DB *database.Handler
}

// MesheryFilterPage represents a page of performance profiles
type MesheryFilterPage struct {
	Page       uint64           `json:"page"`
	PageSize   uint64           `json:"page_size"`
	TotalCount int              `json:"total_count"`
	Filters    []*MesheryFilter `json:"filters"`
}

// GetMesheryFilters returns all of the performance profiles
func (mfp *MesheryFilterPersister) GetMesheryFilters(search, order string, page, pageSize uint64) ([]byte, error) {
	if order == "" {
		order = "updated_at desc"
	}

	count := int64(0)
	filters := []*MesheryFilter{}

	query := mfp.DB.Order(order)

	if search != "" {
		like := "%" + strings.ToLower(search) + "%"
		query = query.Where("(lower(meshery_filters.name) like ?)", like)
	}

	query.Table("meshery_filters").Count(&count)

	Paginate(uint(page), uint(pageSize))(query).Find(&filters)

	mesheryFilterPage := &MesheryFilterPage{
		Page:       page,
		PageSize:   pageSize,
		TotalCount: int(count),
		Filters:    filters,
	}

	return marshalMesheryFilterPage(mesheryFilterPage), nil
}

func (mfp *MesheryFilterPersister) SaveMesheryFilter(filter *MesheryFilter) ([]byte, error) {
	if filter.ID == nil {
		id, err := uuid.NewV4()
		if err != nil {
			return nil, fmt.Errorf("failed to create ID for the filter: %s", err)
		}

		filter.ID = &id
	}
	return marshalMesheryFilter(filter), mfp.DB.Save(filter).Error
}

func marshalMesheryFilter(mf *MesheryFilter) []byte {
	res, _ := json.Marshal(mf)

	return res
}

func marshalMesheryFilterPage(mfp *MesheryFilterPage) []byte {
	res, _ := json.Marshal(mfp)

	return res
}
