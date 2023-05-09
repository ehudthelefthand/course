package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ehudthelefthand/course/gorm"
	"github.com/ehudthelefthand/course/model"
	"github.com/ehudthelefthand/course/server"
	"github.com/stretchr/testify/assert"
)

type mockDB struct{}

func (m *mockDB) GetAllCourse() ([]model.Course, error) {
	return []model.Course{
		{ID: 1, Name: "This is a mock", Description: "Mock!"},
		{ID: 2, Name: "This is a mock 2", Description: "Mock 2!"},
	}, nil
}

// GET /courses
func TestListCourses_Success_ShoudlReturn200(t *testing.T) {
	gorm, err := gorm.NewDB()
	if err != nil {
		t.Error(err)
	}
	mdb := &mockDB{}

	r := server.Init(mdb, gorm)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/courses", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	res := []model.Course{}
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Error(err)
	}
	assert.Equal(t, 2, len(res))
	assert.Equal(t, uint(1), res[0].ID)
	assert.Equal(t, "This is a mock", res[0].Name)
	assert.Equal(t, "Mock!", res[0].Description)
}
