package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/MrAndreID/goechoms/applications"
	"github.com/MrAndreID/goechoms/applications/databases/models"
	"github.com/MrAndreID/goechoms/applications/types"
	"github.com/MrAndreID/goechoms/configs"
	"github.com/google/uuid"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type UserHandler struct {
	Config      *configs.Config
	Application *applications.Application
}

func NewUserHandler(cfg *configs.Config, app *applications.Application) *UserHandler {
	return &UserHandler{
		Config:      cfg,
		Application: app,
	}
}

func (uh *UserHandler) Index(c echo.Context) error {
	var (
		request   types.GetUserRequest
		paginator types.PaginatorResponse
		user      []models.User
		orderBy   map[string]string = map[string]string{
			"id":        "id",
			"name":      "name",
			"createdAt": "created_at",
			"updatedAt": "updated_at",
		}
		sortBy map[string]string = map[string]string{
			"asc":  "asc",
			"desc": "desc",
		}
		search []string = []string{"name"}
		total  int64
	)

	if err := uh.Application.BindRequest(c, &request); err != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   "Applications.Handlers.User.Index.01",
			"error": err.(*echo.HTTPError).Message,
		}).Error("invalid request data")

		return c.JSON(http.StatusBadRequest, types.MainResponse{
			Code:        fmt.Sprintf("%04d", http.StatusBadRequest),
			Description: strings.ToUpper(strings.ReplaceAll(http.StatusText(http.StatusBadRequest), " ", "_")),
			Data:        err.(*echo.HTTPError).Message,
		})
	}

	countTotal := uh.Application.Database.Model(&models.User{})
	queryBuilder := uh.Application.Database.Model(&models.User{}).Preload("Emails")

	if request.ID != "" {
		countTotal.Where("id = ?", request.ID)

		queryBuilder.Where("id = ?", request.ID)
	}

	uh.Application.DataTable(
		c.Request().Context(),
		queryBuilder,
		search,
		orderBy[request.OrderBy],
		sortBy[request.SortBy],
		orderBy["id"],
		"asc",
		request.Page,
		&request.Limit,
		request.Search,
		false,
	)

	queryBuilder.Find(&user)

	if request.DisableCalculateTotal != "true" {
		countTotal.Count(&total)
	}

	if len(user) >= request.Limit {
		paginator.NextPage = true
	}

	paginator.Data = user
	paginator.Total = total

	return c.JSON(http.StatusOK, types.MainResponse{
		Code:        fmt.Sprintf("%04d", http.StatusOK),
		Description: "SUCCESS",
		Data:        paginator,
	})
}

func (uh *UserHandler) Create(c echo.Context) error {
	var (
		tag     string = "Applications.Handlers.User.Create."
		request types.CreateUserRequest
		user    models.User
	)

	if err := uh.Application.BindRequest(c, &request); err != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "01",
			"error": err.(*echo.HTTPError).Message,
		}).Error("invalid request data")

		return c.JSON(http.StatusBadRequest, types.MainResponse{
			Code:        fmt.Sprintf("%04d", http.StatusBadRequest),
			Description: strings.ToUpper(strings.ReplaceAll(http.StatusText(http.StatusBadRequest), " ", "_")),
			Data:        err.(*echo.HTTPError).Message,
		})
	}

	UUID, err := uuid.NewRandom()

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "02",
			"error": err.Error(),
		}).Error("failed to generate uuid")

		return c.JSON(http.StatusInternalServerError, types.MainResponse{
			Code:        fmt.Sprintf("%04d", http.StatusInternalServerError),
			Description: strings.ToUpper(strings.ReplaceAll(http.StatusText(http.StatusInternalServerError), " ", "_")),
		})
	}

	user.ID = UUID.String()
	user.Name = request.Name
	user.CreatedAt = time.Now().In(uh.Application.TimeLocation)
	user.UpdatedAt = time.Now().In(uh.Application.TimeLocation)

	createUser := uh.Application.Database.Save(&user)

	if createUser.Error != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "03",
			"error": createUser.Error.Error(),
		}).Error("failed to create post")

		return c.JSON(http.StatusInternalServerError, types.MainResponse{
			Code:        fmt.Sprintf("%04d", http.StatusInternalServerError),
			Description: strings.ToUpper(strings.ReplaceAll(http.StatusText(http.StatusInternalServerError), " ", "_")),
		})
	}

	if createUser.RowsAffected == 0 {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "04",
			"error": "Failed to Create Post",
		}).Error("failed to create post")

		return c.JSON(http.StatusInternalServerError, types.MainResponse{
			Code:        fmt.Sprintf("%04d", http.StatusInternalServerError),
			Description: strings.ToUpper(strings.ReplaceAll(http.StatusText(http.StatusInternalServerError), " ", "_")),
		})
	}

	return c.JSON(http.StatusCreated, types.MainResponse{
		Code:        fmt.Sprintf("%04d", http.StatusCreated),
		Description: strings.ToUpper(strings.ReplaceAll(http.StatusText(http.StatusCreated), " ", "_")),
		Data:        user,
	})
}

func (uh *UserHandler) Edit(c echo.Context) error {
	var (
		tag     string = "Applications.Handlers.User.Edit."
		request types.EditUserRequest
		user    models.User
	)

	if err := uh.Application.BindRequest(c, &request); err != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "01",
			"error": err.(*echo.HTTPError).Message,
		}).Error("invalid request data")

		return c.JSON(http.StatusBadRequest, types.MainResponse{
			Code:        fmt.Sprintf("%04d", http.StatusBadRequest),
			Description: strings.ToUpper(strings.ReplaceAll(http.StatusText(http.StatusBadRequest), " ", "_")),
			Data:        err.(*echo.HTTPError).Message,
		})
	}

	userResult := uh.Application.Database.First(&user, "id = ?", request.ID)

	if userResult.RowsAffected == 0 {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "02",
			"error": "Failed to Get User Data",
		}).Error("failed to get user data")

		return c.JSON(http.StatusNotFound, types.MainResponse{
			Code:        fmt.Sprintf("%04d", http.StatusNotFound),
			Description: strings.ToUpper(strings.ReplaceAll(http.StatusText(http.StatusNotFound), " ", "_")),
		})
	}

	if request.Name != "" {
		user.Name = request.Name
	}

	user.UpdatedAt = time.Now().In(uh.Application.TimeLocation)

	editUser := uh.Application.Database.Save(&user)

	if editUser.Error != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "03",
			"error": editUser.Error.Error(),
		}).Error("failed to update user data")

		return c.JSON(http.StatusInternalServerError, types.MainResponse{
			Code:        fmt.Sprintf("%04d", http.StatusInternalServerError),
			Description: strings.ToUpper(strings.ReplaceAll(http.StatusText(http.StatusInternalServerError), " ", "_")),
		})
	}

	if editUser.RowsAffected == 0 {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "04",
			"error": "Failed to Update User Data",
		}).Error("failed to update user data")

		return c.JSON(http.StatusInternalServerError, types.MainResponse{
			Code:        fmt.Sprintf("%04d", http.StatusInternalServerError),
			Description: strings.ToUpper(strings.ReplaceAll(http.StatusText(http.StatusInternalServerError), " ", "_")),
		})
	}

	return c.JSON(http.StatusOK, types.MainResponse{
		Code:        fmt.Sprintf("%04d", http.StatusOK),
		Description: "SUCCESS",
		Data:        user,
	})
}

func (uh *UserHandler) Delete(c echo.Context) error {
	var (
		tag     string = "Applications.Handlers.User.Delete."
		request types.DeleteUserRequest
		user    models.User
	)

	if err := uh.Application.BindRequest(c, &request); err != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "01",
			"error": err.(*echo.HTTPError).Message,
		}).Error("invalid request data")

		return c.JSON(http.StatusBadRequest, types.MainResponse{
			Code:        fmt.Sprintf("%04d", http.StatusBadRequest),
			Description: strings.ToUpper(strings.ReplaceAll(http.StatusText(http.StatusBadRequest), " ", "_")),
			Data:        err.(*echo.HTTPError).Message,
		})
	}

	userResult := uh.Application.Database.First(&user, "id = ?", request.ID)

	if userResult.RowsAffected == 0 {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "02",
			"error": "Failed To Get User Data",
		}).Error("failed to get user data")

		return c.JSON(http.StatusNotFound, types.MainResponse{
			Code:        fmt.Sprintf("%04d", http.StatusNotFound),
			Description: strings.ToUpper(strings.ReplaceAll(http.StatusText(http.StatusNotFound), " ", "_")),
		})
	}

	deleteUser := uh.Application.Database.Delete(&user, "id = ?", request.ID)

	if deleteUser.Error != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "03",
			"error": deleteUser.Error.Error(),
		}).Error("failed to delete user data")

		return c.JSON(http.StatusInternalServerError, types.MainResponse{
			Code:        fmt.Sprintf("%04d", http.StatusInternalServerError),
			Description: strings.ToUpper(strings.ReplaceAll(http.StatusText(http.StatusInternalServerError), " ", "_")),
		})
	}

	if deleteUser.RowsAffected == 0 {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "04",
			"error": "Failed To Delete User Data",
		}).Error("failed to delete user data")

		return c.JSON(http.StatusInternalServerError, types.MainResponse{
			Code:        fmt.Sprintf("%04d", http.StatusInternalServerError),
			Description: strings.ToUpper(strings.ReplaceAll(http.StatusText(http.StatusInternalServerError), " ", "_")),
		})
	}

	return c.JSON(http.StatusOK, types.MainResponse{
		Code:        fmt.Sprintf("%04d", http.StatusOK),
		Description: "SUCCESS",
	})
}
