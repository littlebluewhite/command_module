package schedule

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
	command2 "new_command/model/command"
	"new_command/pkg/response"
	"strconv"
	"strings"
)

type Handler struct {
	DB       *gorm.DB
	Response response.Response
	C        *cache.Cache
}

// GetSchedules swagger
// @Summary Show all schedules
// @Description Get all schedules
// @Tags schedule
// @Produce json
// @Success 200 {array} schedule.SwaggerResponse
// @Router /schedule/api [get]
func (h *Handler) GetSchedules(c *gin.Context) {
	var schedules []Schedule
	cacheMap := map[int]Schedule{}
	if x, found := h.C.Get("schedules"); found {
		cacheMap = x.(map[int]Schedule)
	} else {
		h.Response.HttpResponse(c, http.StatusNotFound, "cache error").Printf(
			response.DefaultMessage, http.StatusNotFound, "cache error")
		return
	}
	for _, value := range cacheMap {
		schedules = append(schedules, value)
	}

	h.Response.HttpResponse(c, http.StatusOK, &schedules).Printf(response.DefaultMessage, http.StatusOK, &schedules)
}

// GetScheduleById swagger
// @Summary Show schedules
// @Description Get schedules by id
// @Tags schedule
// @Produce json
// @Param id path int true "schedule id"
// @Success 200 {object} schedule.SwaggerResponse
// @Router /schedule/api/{id} [get]
func (h *Handler) GetScheduleById(c *gin.Context) {
	id := c.Param("id")
	cacheMap := map[int]Schedule{}
	if x, found := h.C.Get("schedules"); found {
		cacheMap = x.(map[int]Schedule)
	} else {
		h.Response.HttpResponse(c, http.StatusBadRequest, gin.H{"message": "cache error"}).Printf(
			response.DefaultMessage, http.StatusBadRequest, gin.H{"message": "cache error"})
		return
	}
	idInt, _ := strconv.Atoi(id)
	schedule, ok := cacheMap[idInt]
	if !ok {
		h.Response.HttpResponse(c, http.StatusBadRequest, gin.H{"message": "id is not correct"}).Printf(
			response.DefaultMessage, http.StatusBadRequest, gin.H{"message": "id is not correct"})
		return
	}

	h.Response.HttpResponse(c, http.StatusOK, &schedule).Printf(
		response.DefaultMessage, http.StatusOK, &schedule)
}

// AddSchedule swagger
// @Summary Create schedules
// @Tags schedule
// @Accept json
// @Produce json
// @Param schedule body schedule.SwaggerCreate true "schedule body"
// @Success 200 {string} string "created success"
// @Router /schedule/api [post]
func (h *Handler) AddSchedule(c *gin.Context) {
	entry := &Schedule{}
	if err := c.ShouldBindJSON(entry); err != nil {
		h.Response.HttpResponse(c, http.StatusNotAcceptable, err).Printf(response.DefaultMessage, http.StatusNotAcceptable, err)
		return
	}
	if err := entry.CheckSchedule(); err != nil {
		h.Response.HttpResponse(c, http.StatusNotAcceptable, err).Printf(response.DefaultMessage, http.StatusNotAcceptable, err)
		return
	}
	e := h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(entry).Error; err != nil {
			return err
		}
		tx.Preload("Command", func(db *gorm.DB) *gorm.DB {
			return db.Preload("HttpsCommand")
		}).Preload("TimeData").Last(entry)
		if x, found := h.C.Get("schedules"); found {
			cacheMap := x.(map[int]Schedule)
			cacheMap[entry.ID] = *entry
			h.C.Set("schedules", cacheMap, cache.NoExpiration)
		} else {
			return errors.New("cache error")
		}
		return nil
	})
	if e != nil {
		h.Response.HttpResponse(c, 403, e).Printf(response.DefaultMessage, 403, e)
		return
	}
	h.Response.HttpResponse(c, http.StatusCreated, gin.H{
		"message": "created success",
	}).Printf(response.DefaultMessage, http.StatusCreated, gin.H{
		"message": "created success"})
}

// UpdateSchedule swagger
// @Summary Update schedules
// @Tags schedule
// @Accept json
// @Produce json
// @Param id path int true "schedule id"
// @Param schedule body schedule.SwaggerUpdate true "modify schedule body"
// @Success 200 {object} schedule.SwaggerResponse
// @Router /schedule/api/{id} [patch]
func (h *Handler) UpdateSchedule(c *gin.Context) {
	entry := &SchedulePatch{}
	if err := c.ShouldBindJSON(entry); err != nil {
		h.Response.HttpResponse(c, http.StatusNotAcceptable, err).Printf(response.DefaultMessage, http.StatusNotAcceptable, err)
		return
	}
	var schedule Schedule
	id := c.Param("id")
	h.DB.Preload(clause.Associations).First(&schedule, id)
	if schedule.ID == 0 {
		h.Response.HttpResponse(c, http.StatusBadRequest, gin.H{"message": "id is not correct"}).Printf(
			response.DefaultMessage, http.StatusBadRequest, gin.H{"message": "id is not correct"})
		return
	}
	schedule.UpdateData(*entry)
	if err := schedule.CheckSchedule(); err != nil {
		h.Response.HttpResponse(c, http.StatusNotAcceptable, err).Printf(response.DefaultMessage, http.StatusNotAcceptable, err)
		return
	}
	err := h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&schedule).Error; err != nil {
			return err
		}
		var command command2.Command
		if entry.CommandID != 0 {
			tx.Preload(clause.Associations).First(&command, entry.CommandID)
			schedule.Command = &command
		}
		if err := tx.Save(schedule.TimeData).Error; err != nil {
			return err
		}
		if x, found := h.C.Get("schedules"); found {
			cacheMap := x.(map[int]Schedule)
			cacheMap[schedule.ID] = schedule
			h.C.Set("schedules", cacheMap, cache.NoExpiration)
		} else {
			return errors.New("cache error")
		}
		return nil
	})
	if err != nil {
		h.Response.HttpResponse(c, http.StatusBadRequest, err).Printf(
			response.DefaultMessage, http.StatusBadRequest, err)
		return
	}
	h.Response.HttpResponse(c, http.StatusOK, &schedule).Printf(
		response.DefaultMessage, http.StatusOK, &schedule)
}

// DeleteSchedule swagger
// @Summary Delete schedules
// @Tags schedule
// @Produce json
// @Param id path int true "schedule id"
// @Success 200 {string} string "delete successfully"
// @Router /schedule/api/{id} [delete]
func (h *Handler) DeleteSchedule(c *gin.Context) {
	idString := c.Param("id")
	var schedule Schedule
	h.DB.Preload(clause.Associations).First(&schedule, idString)
	err := h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&schedule).Error; err != nil {
			return err
		}
		if err := tx.Delete(schedule.TimeData).Error; err != nil {
			return err
		}
		if x, found := h.C.Get("schedules"); found {
			cacheMap := x.(map[int]Schedule)
			delete(cacheMap, schedule.ID)
			h.C.Set("schedules", cacheMap, cache.NoExpiration)
		} else {
			return errors.New("cache error")
		}
		return nil
	})
	if err != nil {
		h.Response.HttpResponse(c, http.StatusBadRequest, err).Printf(
			response.DefaultMessage, http.StatusBadRequest, err)
		return
	}
	var sb strings.Builder
	sb.WriteString("id: ")
	sb.WriteString(idString)
	sb.WriteString(" has been deleted successfully")
	h.Response.HttpResponse(c, http.StatusOK, gin.H{"message": sb.String()}).Printf(
		response.DefaultMessage, http.StatusOK, gin.H{"message": sb.String()})
}
