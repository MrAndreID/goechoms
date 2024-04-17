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

	tx := uh.Application.Database.Begin()

	userUUID, err := uuid.NewRandom()

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "02",
			"error": err.Error(),
		}).Error("failed to generate uuid")

		tx.Rollback()

		return c.JSON(http.StatusInternalServerError, types.MainResponse{
			Code:        fmt.Sprintf("%04d", http.StatusInternalServerError),
			Description: strings.ToUpper(strings.ReplaceAll(http.StatusText(http.StatusInternalServerError), " ", "_")),
		})
	}

	user.ID = userUUID.String()
	user.CreatedAt = time.Now().In(uh.Application.TimeLocation)
	user.UpdatedAt = time.Now().In(uh.Application.TimeLocation)
	user.Name = request.Name

	createUser := tx.Save(&user)

	if createUser.Error != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "03",
			"error": createUser.Error.Error(),
		}).Error("failed to create user")

		tx.Rollback()

		return c.JSON(http.StatusInternalServerError, types.MainResponse{
			Code:        fmt.Sprintf("%04d", http.StatusInternalServerError),
			Description: strings.ToUpper(strings.ReplaceAll(http.StatusText(http.StatusInternalServerError), " ", "_")),
		})
	}

	if createUser.RowsAffected == 0 {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "04",
			"error": "Failed to Create User",
		}).Error("failed to create user")

		tx.Rollback()

		return c.JSON(http.StatusInternalServerError, types.MainResponse{
			Code:        fmt.Sprintf("%04d", http.StatusInternalServerError),
			Description: strings.ToUpper(strings.ReplaceAll(http.StatusText(http.StatusInternalServerError), " ", "_")),
		})
	}

	for i := 0; i < len(request.Emails); i++ {
		for j := i + 1; j < len(request.Emails); j++ {
			if request.Emails[i].Email == request.Emails[j].Email {
				logrus.WithFields(logrus.Fields{
					"tag":   tag + "05",
					"error": "Duplicate Email",
				}).Error("duplicate email")

				tx.Rollback()

				return c.JSON(http.StatusConflict, types.MainResponse{
					Code:        fmt.Sprintf("%04d", http.StatusConflict),
					Description: strings.ToUpper(strings.ReplaceAll(http.StatusText(http.StatusConflict), " ", "_")),
				})
			}
		}
	}

	for _, v := range request.Emails {
		var email models.Email

		emailUUID, err := uuid.NewRandom()

		if err != nil {
			logrus.WithFields(logrus.Fields{
				"tag":   tag + "06",
				"error": err.Error(),
			}).Error("failed to generate uuid")

			tx.Rollback()

			return c.JSON(http.StatusInternalServerError, types.MainResponse{
				Code:        fmt.Sprintf("%04d", http.StatusInternalServerError),
				Description: strings.ToUpper(strings.ReplaceAll(http.StatusText(http.StatusInternalServerError), " ", "_")),
			})
		}

		email.ID = emailUUID.String()
		email.CreatedAt = time.Now().In(uh.Application.TimeLocation)
		email.UpdatedAt = time.Now().In(uh.Application.TimeLocation)
		email.UserID = user.ID
		email.Email = v.Email

		createEmail := tx.Save(&email)

		if createEmail.Error != nil {
			logrus.WithFields(logrus.Fields{
				"tag":   tag + "07",
				"error": createEmail.Error.Error(),
			}).Error("failed to create email")

			tx.Rollback()

			return c.JSON(http.StatusInternalServerError, types.MainResponse{
				Code:        fmt.Sprintf("%04d", http.StatusInternalServerError),
				Description: strings.ToUpper(strings.ReplaceAll(http.StatusText(http.StatusInternalServerError), " ", "_")),
			})
		}

		if createEmail.RowsAffected == 0 {
			logrus.WithFields(logrus.Fields{
				"tag":   tag + "08",
				"error": "Failed to Create Email",
			}).Error("failed to create email")

			tx.Rollback()

			return c.JSON(http.StatusInternalServerError, types.MainResponse{
				Code:        fmt.Sprintf("%04d", http.StatusInternalServerError),
				Description: strings.ToUpper(strings.ReplaceAll(http.StatusText(http.StatusInternalServerError), " ", "_")),
			})
		}

		user.Emails = append(user.Emails, email)
	}

	tx.Commit()

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

	tx := uh.Application.Database.Begin()

	userResult := tx.First(&user, "id = ?", request.ID)

	if userResult.RowsAffected == 0 {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "02",
			"error": "Failed to Get User Data",
		}).Error("failed to get user data")

		tx.Rollback()

		return c.JSON(http.StatusNotFound, types.MainResponse{
			Code:        fmt.Sprintf("%04d", http.StatusNotFound),
			Description: strings.ToUpper(strings.ReplaceAll(http.StatusText(http.StatusNotFound), " ", "_")),
		})
	}

	if request.Name != "" {
		user.Name = request.Name
	}

	if len(request.Emails) > 0 {
		deleteEmail := tx.Where("user_id = ?", user.ID).Delete(&models.Email{})

		if deleteEmail.Error != nil {
			logrus.WithFields(logrus.Fields{
				"tag":   tag + "03",
				"error": deleteEmail.Error.Error(),
			}).Error("failed to delete email data")

			tx.Rollback()

			return c.JSON(http.StatusInternalServerError, types.MainResponse{
				Code:        fmt.Sprintf("%04d", http.StatusInternalServerError),
				Description: strings.ToUpper(strings.ReplaceAll(http.StatusText(http.StatusInternalServerError), " ", "_")),
			})
		}

		if deleteEmail.RowsAffected == 0 {
			logrus.WithFields(logrus.Fields{
				"tag":   tag + "04",
				"error": "Failed to Delete Email Data",
			}).Error("failed to delete email data")

			tx.Rollback()

			return c.JSON(http.StatusInternalServerError, types.MainResponse{
				Code:        fmt.Sprintf("%04d", http.StatusInternalServerError),
				Description: strings.ToUpper(strings.ReplaceAll(http.StatusText(http.StatusInternalServerError), " ", "_")),
			})
		}

		for _, v := range request.Emails {
			var email models.Email

			emailUUID, err := uuid.NewRandom()

			if err != nil {
				logrus.WithFields(logrus.Fields{
					"tag":   tag + "05",
					"error": err.Error(),
				}).Error("failed to generate uuid")

				tx.Rollback()

				return c.JSON(http.StatusInternalServerError, types.MainResponse{
					Code:        fmt.Sprintf("%04d", http.StatusInternalServerError),
					Description: strings.ToUpper(strings.ReplaceAll(http.StatusText(http.StatusInternalServerError), " ", "_")),
				})
			}

			email.ID = emailUUID.String()
			email.CreatedAt = time.Now().In(uh.Application.TimeLocation)
			email.UpdatedAt = time.Now().In(uh.Application.TimeLocation)
			email.UserID = user.ID
			email.Email = v.Email

			createEmail := tx.Save(&email)

			if createEmail.Error != nil {
				logrus.WithFields(logrus.Fields{
					"tag":   tag + "06",
					"error": createEmail.Error.Error(),
				}).Error("failed to create email")

				tx.Rollback()

				return c.JSON(http.StatusInternalServerError, types.MainResponse{
					Code:        fmt.Sprintf("%04d", http.StatusInternalServerError),
					Description: strings.ToUpper(strings.ReplaceAll(http.StatusText(http.StatusInternalServerError), " ", "_")),
				})
			}

			if createEmail.RowsAffected == 0 {
				logrus.WithFields(logrus.Fields{
					"tag":   tag + "07",
					"error": "Failed to Create Email",
				}).Error("failed to create email")

				tx.Rollback()

				return c.JSON(http.StatusInternalServerError, types.MainResponse{
					Code:        fmt.Sprintf("%04d", http.StatusInternalServerError),
					Description: strings.ToUpper(strings.ReplaceAll(http.StatusText(http.StatusInternalServerError), " ", "_")),
				})
			}

			user.Emails = append(user.Emails, email)
		}
	} else {
		var emails []models.Email

		emailResult := tx.Find(&emails, "user_id = ?", user.ID)

		if emailResult.RowsAffected == 0 {
			logrus.WithFields(logrus.Fields{
				"tag":   tag + "08",
				"error": "Failed to Get Email Data",
			}).Error("failed to get email data")

			tx.Rollback()

			return c.JSON(http.StatusNotFound, types.MainResponse{
				Code:        fmt.Sprintf("%04d", http.StatusNotFound),
				Description: strings.ToUpper(strings.ReplaceAll(http.StatusText(http.StatusNotFound), " ", "_")),
			})
		}

		user.Emails = emails
	}

	user.UpdatedAt = time.Now().In(uh.Application.TimeLocation)

	editUser := tx.Save(&user)

	if editUser.Error != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "09",
			"error": editUser.Error.Error(),
		}).Error("failed to edit user data")

		tx.Rollback()

		return c.JSON(http.StatusInternalServerError, types.MainResponse{
			Code:        fmt.Sprintf("%04d", http.StatusInternalServerError),
			Description: strings.ToUpper(strings.ReplaceAll(http.StatusText(http.StatusInternalServerError), " ", "_")),
		})
	}

	if editUser.RowsAffected == 0 {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "10",
			"error": "Failed to Edit User Data",
		}).Error("failed to edit user data")

		tx.Rollback()

		return c.JSON(http.StatusInternalServerError, types.MainResponse{
			Code:        fmt.Sprintf("%04d", http.StatusInternalServerError),
			Description: strings.ToUpper(strings.ReplaceAll(http.StatusText(http.StatusInternalServerError), " ", "_")),
		})
	}

	tx.Commit()

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

	tx := uh.Application.Database.Begin()

	userResult := tx.First(&user, "id = ?", request.ID)

	if userResult.RowsAffected == 0 {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "02",
			"error": "Failed To Get User Data",
		}).Error("failed to get user data")

		tx.Rollback()

		return c.JSON(http.StatusNotFound, types.MainResponse{
			Code:        fmt.Sprintf("%04d", http.StatusNotFound),
			Description: strings.ToUpper(strings.ReplaceAll(http.StatusText(http.StatusNotFound), " ", "_")),
		})
	}

	deleteUser := tx.Delete(&user, "id = ?", request.ID)

	if deleteUser.Error != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "03",
			"error": deleteUser.Error.Error(),
		}).Error("failed to delete user data")

		tx.Rollback()

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

		tx.Rollback()

		return c.JSON(http.StatusInternalServerError, types.MainResponse{
			Code:        fmt.Sprintf("%04d", http.StatusInternalServerError),
			Description: strings.ToUpper(strings.ReplaceAll(http.StatusText(http.StatusInternalServerError), " ", "_")),
		})
	}

	deleteEmail := tx.Where("user_id = ?", request.ID).Delete(&models.Email{})

	if deleteEmail.Error != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "05",
			"error": deleteEmail.Error.Error(),
		}).Error("failed to delete email data")

		tx.Rollback()

		return c.JSON(http.StatusInternalServerError, types.MainResponse{
			Code:        fmt.Sprintf("%04d", http.StatusInternalServerError),
			Description: strings.ToUpper(strings.ReplaceAll(http.StatusText(http.StatusInternalServerError), " ", "_")),
		})
	}

	if deleteEmail.RowsAffected == 0 {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "06",
			"error": "Failed To Delete Email Data",
		}).Error("failed to delete email data")

		tx.Rollback()

		return c.JSON(http.StatusInternalServerError, types.MainResponse{
			Code:        fmt.Sprintf("%04d", http.StatusInternalServerError),
			Description: strings.ToUpper(strings.ReplaceAll(http.StatusText(http.StatusInternalServerError), " ", "_")),
		})
	}

	tx.Commit()

	return c.JSON(http.StatusOK, types.MainResponse{
		Code:        fmt.Sprintf("%04d", http.StatusOK),
		Description: "SUCCESS",
	})
}
