package contract

import (
	"math"
	"strings"
)

/* ---------------------------------------------------------------- */
/*                              filter                              */
/* ---------------------------------------------------------------- */

type Filters struct {
	Page     int    `validate:"omitempty,min=0,max=1000"`
	PageSize int    `validate:"omitempty,min=0,max=100"`
	Sort     string `validate:"omitempty,oneof=game_id title price year stock -game_id -title -price -year -stock ''"`
}

func (f Filters) SortColumn() string {
	return strings.TrimPrefix(f.Sort, "-")
}

func (f Filters) SortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}

	return "ASC"
}

func (f Filters) Limit() int {
	return f.PageSize
}

func (f Filters) Offset() int {
	return (f.Page - 1) * f.PageSize
}

/* ---------------------------------------------------------------- */
/*                             metadata                             */
/* ---------------------------------------------------------------- */

type Metadata struct {
	CurrentPage  int `json:"current_page,omitempty"`
	PageSize     int `json:"page_size,omitempty"`
	FirstPage    int `json:"first_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
}

func CalculateMetadata(totalRecords, page, pageSize int) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}

	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}
}
