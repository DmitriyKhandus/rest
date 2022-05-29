package user

import (
	"fmt"
	"net/http"

	"github.com/DmitriyKhandus/rest-api/internal/apperror"
	"github.com/DmitriyKhandus/rest-api/internal/handlers"
	"github.com/DmitriyKhandus/rest-api/pkg/logging"
	"github.com/julienschmidt/httprouter"
)

const (
	usersUrl = "/users"
	userUrl  = "/user/:uuid"
)

type handler struct {
	logger *logging.Logger
}

func NewHandler(logger *logging.Logger) handlers.Handler {
	return &handler{
		logger: logger,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, usersUrl, apperror.Middleware(h.GetList))
	router.HandlerFunc(http.MethodGet, userUrl, apperror.Middleware(h.GetUserByUUID))
	router.HandlerFunc(http.MethodPost, userUrl, apperror.Middleware(h.CreateUser))
	router.HandlerFunc(http.MethodPut, userUrl, apperror.Middleware(h.UpdateUser))
	router.HandlerFunc(http.MethodPatch, userUrl, apperror.Middleware(h.PartiallyUpdateUser))
	router.HandlerFunc(http.MethodDelete, userUrl, apperror.Middleware(h.DeleteUser))
}

func (h *handler) GetList(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info("Use getList of users")
	return apperror.ErrNotFound

}

func (h *handler) GetUserByUUID(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info("Get one user by uuid")
	return apperror.NewAppError(nil, "test", "test", "t13")
}
func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info("Use create user")
	return fmt.Errorf("this is API error")
}
func (h *handler) UpdateUser(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info("Use update user")
	w.WriteHeader(200)
	w.Write([]byte("This is update user"))
	return nil
}
func (h *handler) PartiallyUpdateUser(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info("Use partially update user")
	w.WriteHeader(200)
	w.Write([]byte("This is partially updated"))
	return nil
}
func (h *handler) DeleteUser(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info("Use delete user")
	w.WriteHeader(200)
	w.Write([]byte("This is delete user"))
	return nil
}
