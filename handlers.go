package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	db *pgxpool.Pool
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

type formErrors struct {
	Name  string
	Email string
}

func validateContact(name, email string) (formErrors, bool) {
	var errs formErrors
	valid := true
	if strings.TrimSpace(name) == "" {
		errs.Name = "Name is required"
		valid = false
	}
	if strings.TrimSpace(email) == "" {
		errs.Email = "Email is required"
		valid = false
	} else if !emailRegex.MatchString(email) {
		errs.Email = "Email must be a valid format"
		valid = false
	}
	return errs, valid
}

// Page handlers

func (h *Handler) ListContacts(c echo.Context) error {
	contacts, err := GetAllContacts(c.Request().Context(), h.db)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to load contacts")
	}
	flash := c.QueryParam("flash")
	return c.Render(http.StatusOK, "list", map[string]interface{}{
		"Contacts": contacts,
		"Flash":    flash,
	})
}

func (h *Handler) NewContactForm(c echo.Context) error {
	return c.Render(http.StatusOK, "new", map[string]interface{}{
		"Contact": &Contact{},
		"Errors":  formErrors{},
	})
}

func (h *Handler) CreateContact(c echo.Context) error {
	contact := &Contact{
		Name:  c.FormValue("name"),
		Email: c.FormValue("email"),
		Phone: c.FormValue("phone"),
		Notes: c.FormValue("notes"),
	}

	errs, valid := validateContact(contact.Name, contact.Email)
	if !valid {
		return c.Render(http.StatusUnprocessableEntity, "new", map[string]interface{}{
			"Contact": contact,
			"Errors":  errs,
		})
	}

	if err := CreateContactDB(c.Request().Context(), h.db, contact); err != nil {
		return c.String(http.StatusInternalServerError, "Failed to create contact")
	}
	return c.Redirect(http.StatusSeeOther, "/?flash=Contact+created+successfully")
}

func (h *Handler) ViewContact(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid ID")
	}
	contact, err := GetContactByID(c.Request().Context(), h.db, id)
	if err != nil {
		return c.String(http.StatusNotFound, "Contact not found")
	}
	flash := c.QueryParam("flash")
	return c.Render(http.StatusOK, "view", map[string]interface{}{
		"Contact": contact,
		"Flash":   flash,
	})
}

func (h *Handler) EditContactForm(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid ID")
	}
	contact, err := GetContactByID(c.Request().Context(), h.db, id)
	if err != nil {
		return c.String(http.StatusNotFound, "Contact not found")
	}
	return c.Render(http.StatusOK, "edit", map[string]interface{}{
		"Contact": contact,
		"Errors":  formErrors{},
	})
}

func (h *Handler) UpdateContact(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid ID")
	}

	contact := &Contact{
		ID:    id,
		Name:  c.FormValue("name"),
		Email: c.FormValue("email"),
		Phone: c.FormValue("phone"),
		Notes: c.FormValue("notes"),
	}

	errs, valid := validateContact(contact.Name, contact.Email)
	if !valid {
		return c.Render(http.StatusUnprocessableEntity, "edit", map[string]interface{}{
			"Contact": contact,
			"Errors":  errs,
		})
	}

	if err := UpdateContactDB(c.Request().Context(), h.db, contact); err != nil {
		return c.String(http.StatusInternalServerError, "Failed to update contact")
	}
	return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/contacts/%d?flash=Contact+updated+successfully", id))
}

func (h *Handler) DeleteContact(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid ID")
	}
	if err := DeleteContactDB(c.Request().Context(), h.db, id); err != nil {
		return c.String(http.StatusInternalServerError, "Failed to delete contact")
	}
	return c.Redirect(http.StatusSeeOther, "/?flash=Contact+deleted+successfully")
}

// API handlers

func (h *Handler) APIListContacts(c echo.Context) error {
	contacts, err := GetAllContacts(c.Request().Context(), h.db)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to load contacts"})
	}
	if contacts == nil {
		contacts = []Contact{}
	}
	return c.JSON(http.StatusOK, contacts)
}

func (h *Handler) APICreateContact(c echo.Context) error {
	var contact Contact
	if err := c.Bind(&contact); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	errs, valid := validateContact(contact.Name, contact.Email)
	if !valid {
		validationErrors := map[string]string{}
		if errs.Name != "" {
			validationErrors["name"] = errs.Name
		}
		if errs.Email != "" {
			validationErrors["email"] = errs.Email
		}
		return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{"errors": validationErrors})
	}

	if err := CreateContactDB(c.Request().Context(), h.db, &contact); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create contact"})
	}
	return c.JSON(http.StatusCreated, contact)
}

func (h *Handler) APIGetContact(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}
	contact, err := GetContactByID(c.Request().Context(), h.db, id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Contact not found"})
	}
	return c.JSON(http.StatusOK, contact)
}

func (h *Handler) APIUpdateContact(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	var contact Contact
	if err := c.Bind(&contact); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	contact.ID = id

	errs, valid := validateContact(contact.Name, contact.Email)
	if !valid {
		validationErrors := map[string]string{}
		if errs.Name != "" {
			validationErrors["name"] = errs.Name
		}
		if errs.Email != "" {
			validationErrors["email"] = errs.Email
		}
		return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{"errors": validationErrors})
	}

	if err := UpdateContactDB(c.Request().Context(), h.db, &contact); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update contact"})
	}

	updated, _ := GetContactByID(c.Request().Context(), h.db, id)
	if updated != nil {
		return c.JSON(http.StatusOK, updated)
	}
	return c.JSON(http.StatusOK, contact)
}

func (h *Handler) APIDeleteContact(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}
	if err := DeleteContactDB(c.Request().Context(), h.db, id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete contact"})
	}
	return c.NoContent(http.StatusNoContent)
}
