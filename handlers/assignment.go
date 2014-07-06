package handlers

import (
	"errors"

	m "github.com/Lanciv/GoGradeAPI/model"
	"github.com/Lanciv/GoGradeAPI/store"

	"github.com/gin-gonic/gin"
	"github.com/mholt/binding"
)

// CreateAssignment ...
func CreateAssignment(c *gin.Context) {
	a := new(m.Assignment)

	errs := binding.Bind(c.Req, a)
	if errs != nil {
		c.Error(errors.New("validation"), errs)
		c.JSON(400, errs)
		return
	}

	id, err := store.Assignments.Store(a)
	if err != nil {
		writeError(c.Writer, serverError, 500, err)
		return
	}
	a.ID = id

	c.JSON(201, &APIRes{"assignment": []m.Assignment{*a}})
	return
}

// GetAssignment ...
func GetAssignment(c *gin.Context) {

	id := c.Params.ByName("id")

	a := m.Assignment{}
	err := store.Assignments.FindByID(&a, id)
	if err == store.ErrNotFound {
		writeError(c.Writer, notFoundError, 404, nil)
		return
	}
	if err != nil {
		writeError(c.Writer, serverError, 500, nil)
		return
	}

	c.JSON(200, &APIRes{"assignment": []m.Assignment{a}})
	return
}

// UpdateAssignment ...
func UpdateAssignment(c *gin.Context) {
	id := c.Params.ByName("id")

	a := new(m.Assignment)

	errs := binding.Bind(c.Req, a)
	if errs != nil {
		writeError(c.Writer, errs, 400, nil)
		return
	}

	a.ID = id
	err := store.Assignments.Update(a, id)

	if err != nil {
		writeError(c.Writer, "error updating Assignment", 500, err)
		return
	}

	c.JSON(200, &APIRes{"assignment": []m.Assignment{*a}})
	return
}

// GetAllAssignments ...
func GetAllAssignments(c *gin.Context) {
	filter := map[string]string{}
	if c.Req.URL.Query().Get("classId") != "" {
		filter["classId"] = c.Req.URL.Query().Get("classId")
	}
	if c.Req.URL.Query().Get("termId") != "" {
		filter["termId"] = c.Req.URL.Query().Get("termId")
	}
	if c.Req.URL.Query().Get("typeId") != "" {
		filter["typeId"] = c.Req.URL.Query().Get("typeId")
	}

	assignment := []m.Assignment{}
	err := store.Assignments.Filter(&assignment, filter)
	if err != nil {
		writeError(c.Writer, serverError, 500, err)
		return
	}

	c.JSON(200, &APIRes{"assignment": assignment})
	return
}