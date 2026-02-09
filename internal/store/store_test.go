package store

import (
	"testing"
	"time"

	"github.com/couchcryptid/storm-data-graphql-api/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestBuildWhereClause_TimeOnly(t *testing.T) {
	filter := &model.StormReportFilter{
		BeginTimeAfter:  time.Date(2024, 4, 26, 0, 0, 0, 0, time.UTC),
		BeginTimeBefore: time.Date(2024, 4, 27, 0, 0, 0, 0, time.UTC),
	}

	where, args, nextIdx := buildWhereClause(filter)

	assert.Len(t, where, 2)
	assert.Contains(t, where[0], "$1")
	assert.Contains(t, where[1], "$2")
	assert.Len(t, args, 2)
	assert.Equal(t, 3, nextIdx)
}

func TestBuildWhereClause_WithTypes(t *testing.T) {
	filter := &model.StormReportFilter{
		BeginTimeAfter:  time.Date(2024, 4, 26, 0, 0, 0, 0, time.UTC),
		BeginTimeBefore: time.Date(2024, 4, 27, 0, 0, 0, 0, time.UTC),
		Types:           []string{"hail", "tornado"},
	}

	where, args, nextIdx := buildWhereClause(filter)

	assert.Len(t, where, 3)
	assert.Contains(t, where[2], "type = ANY($3)")
	assert.Len(t, args, 3)
	assert.Equal(t, 4, nextIdx)
}

func TestBuildWhereClause_AllArrayFilters(t *testing.T) {
	mag := 1.5
	filter := &model.StormReportFilter{
		BeginTimeAfter:  time.Date(2024, 4, 26, 0, 0, 0, 0, time.UTC),
		BeginTimeBefore: time.Date(2024, 4, 27, 0, 0, 0, 0, time.UTC),
		Types:           []string{"hail"},
		Severity:        []string{"severe"},
		States:          []string{"TX", "OK"},
		Counties:        []string{"Dallas"},
		MinMagnitude:    &mag,
	}

	where, args, nextIdx := buildWhereClause(filter)

	// 2 time + types + severity + states + counties + min_magnitude = 7
	assert.Len(t, where, 7)
	assert.Len(t, args, 7)
	assert.Equal(t, 8, nextIdx)
}

func TestBuildWhereClause_RadiusFilter(t *testing.T) {
	lat := 32.7767
	lon := -96.7970
	radius := 50.0
	filter := &model.StormReportFilter{
		BeginTimeAfter:  time.Date(2024, 4, 26, 0, 0, 0, 0, time.UTC),
		BeginTimeBefore: time.Date(2024, 4, 27, 0, 0, 0, 0, time.UTC),
		NearLat:         &lat,
		NearLon:         &lon,
		RadiusMiles:     &radius,
	}

	where, args, nextIdx := buildWhereClause(filter)

	// 2 time + bounding box (1 clause, 4 params) + haversine (1 clause, 4 params) = 4 clauses
	assert.Len(t, where, 4)
	// 2 time args + 4 bbox args + 4 haversine args = 10
	assert.Len(t, args, 10)
	assert.Equal(t, 11, nextIdx)
}

func TestSortColumn(t *testing.T) {
	tests := []struct {
		input model.SortField
		want  string
	}{
		{model.SortFieldBeginTime, "begin_time"},
		{model.SortFieldMagnitude, "magnitude"},
		{model.SortFieldState, "location_state"},
		{model.SortFieldType, "type"},
		{model.SortField("UNKNOWN"), "begin_time"},
	}

	for _, tt := range tests {
		t.Run(string(tt.input), func(t *testing.T) {
			assert.Equal(t, tt.want, sortColumn(tt.input))
		})
	}
}
