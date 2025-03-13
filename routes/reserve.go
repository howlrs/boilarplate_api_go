package routes

import (
	"backend/models"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/labstack/echo/v4"
	"github.com/rs/xid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ReadReservation is Get reservation
func (p *Client) ReadReservation(c echo.Context) error {
	id := c.QueryParam("id")
	if id == "" {
		return responseHandler(c, http.StatusBadRequest, nil, nil, "Failed to get reservation")
	}

	// read database
	reserve := &models.ReservatedTime{}
	doc, err := p.firestore.Collection(reserve.ToCollection(p.IsTest())).Doc(id).Get(c.Request().Context())
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return responseHandler(c, http.StatusNotFound, nil, err, "Failed to get reservation")
		}

		return responseHandler(c, http.StatusInternalServerError, nil, err, "Failed to get reservation")
	}

	if err := doc.DataTo(reserve); err != nil {
		return responseHandler(c, http.StatusInternalServerError, nil, err, "Failed to get reservation")
	}

	return responseHandler(c, http.StatusOK, reserve, nil, "success, read reservation")
}

// CreateReservation is Create reservation
func (p *Client) CreateReservation(c echo.Context) error {
	reserve := &models.ReservatedTime{}
	if err := c.Bind(reserve); err != nil {
		return responseHandler(c, http.StatusBadRequest, nil, err, "Failed to bind request")
	}

	// set database
	reserve.ID = xid.New().String()
	if _, err := p.firestore.Collection(reserve.ToCollection(p.IsTest())).Doc(reserve.ID).Set(c.Request().Context(), reserve); err != nil {
		return responseHandler(c, http.StatusInternalServerError, nil, err, "Failed to set reservation")
	}

	return responseHandler(c, http.StatusOK, reserve, nil, "success, create reservation")
}

// UpdateReservation is Update reservation
func (p *Client) UpdateReservation(c echo.Context) error {
	reserve := &models.ReservatedTime{}
	if err := c.Bind(reserve); err != nil {
		return responseHandler(c, http.StatusBadRequest, nil, err, "Failed to bind request")
	}

	// update database
	updateFields := []firestore.Update{
		{
			FieldPath: []string{"content"},
			Value:     reserve.Content,
		},
	}
	if _, err := p.firestore.Collection(reserve.ToCollection(p.IsTest())).Doc(reserve.ID).Update(c.Request().Context(), updateFields); err != nil {
		return responseHandler(c, http.StatusInternalServerError, nil, err, "Failed to update reservation")
	}

	return responseHandler(c, http.StatusOK, reserve, nil, "success, update reservation")
}

// DeleteReservation is Cancel reservation
func (p *Client) DeleteReservation(c echo.Context) error {
	id := c.QueryParam("id")
	if id == "" {
		return responseHandler(c, http.StatusBadRequest, nil, nil, "Failed to get reservation")
	}

	// delete database
	reserve := &models.ReservatedTime{}
	updateStatus := []firestore.Update{
		{
			FieldPath: []string{"status"},
			Value:     models.ReserveStateCancel,
		},
	}
	if _, err := p.firestore.Collection(reserve.ToCollection(p.IsTest())).Doc(id).Update(c.Request().Context(), updateStatus); err != nil {
		return responseHandler(c, http.StatusInternalServerError, nil, err, "Failed to delete reservation")
	}

	return responseHandler(c, http.StatusOK, nil, nil, "success, canceled reservation")
}
