package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/daniel-oliveiravas/class-booking-service/business/classes"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler_AddClass(t *testing.T) {
	serverURL, httpClient := setupIntegration(t)
	url := fmt.Sprintf("%s/classes", serverURL)
	newClass := classes.NewClass{
		Name:      uuid.NewString(),
		StartDate: time.Now().UTC(),
		EndDate:   time.Now().UTC().Add(time.Hour * 24 * 10),
		Capacity:  30,
	}
	class := CreateNewClass(t, httpClient, url, newClass)

	assert.NotEmpty(t, class.ID)
	assert.NotEmpty(t, class.CreatedAt)
	assert.NotEmpty(t, class.UpdatedAt)
	assert.Equal(t, newClass.Name, class.Name)
}

func TestHandler_AddClass_InvalidData(t *testing.T) {
	serverURL, httpClient := setupIntegration(t)
	url := fmt.Sprintf("%s/classes", serverURL)

	newClass := classes.NewClass{
		Name:      uuid.NewString(),
		StartDate: time.Now(),
		EndDate:   time.Now().Add(time.Hour * -24),
		Capacity:  30,
	}
	requestBytes, err := json.Marshal(newClass)
	require.NoError(t, err)
	body := bytes.NewBuffer(requestBytes)

	resp, err := httpClient.Post(url, "application/json", body)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var class classes.Class
	err = json.Unmarshal(respBody, &class)
	require.NoError(t, err)

	assert.Empty(t, class.ID)
}

func TestHandler_GetClass(t *testing.T) {
	serverURL, httpClient := setupIntegration(t)
	url := fmt.Sprintf("%s/classes", serverURL)

	newClass := classes.NewClass{
		Name:      uuid.NewString(),
		StartDate: time.Now().UTC(),
		EndDate:   time.Now().UTC().Add(time.Hour * 24 * 10),
		Capacity:  30,
	}
	classCreated := CreateNewClass(t, httpClient, url, newClass)
	assert.NotEmpty(t, classCreated.ID)

	getClassURL := fmt.Sprintf("%s/classes/%s", serverURL, classCreated.ID)
	resp, err := httpClient.Get(getClassURL)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var class classes.Class
	err = json.Unmarshal(respBody, &class)
	require.NoError(t, err)

	assert.NotEmpty(t, class.ID)
}

func TestHandler_GetClass_NotFound(t *testing.T) {
	serverURL, httpClient := setupIntegration(t)

	getClassURL := fmt.Sprintf("%s/classes/%s", serverURL, uuid.NewString())
	resp, err := httpClient.Get(getClassURL)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestHandler_UpdateClass(t *testing.T) {
	serverURL, httpClient := setupIntegration(t)
	url := fmt.Sprintf("%s/classes", serverURL)

	newClass := classes.NewClass{
		Name:      uuid.NewString(),
		StartDate: time.Now().UTC(),
		EndDate:   time.Now().UTC().Add(time.Hour * 24 * 10),
		Capacity:  30,
	}
	classCreated := CreateNewClass(t, httpClient, url, newClass)
	assert.NotEmpty(t, classCreated.ID)

	newName := uuid.NewString()
	updateClass := classes.UpdateClass{
		Name: &newName,
	}
	requestBytes, err := json.Marshal(updateClass)
	require.NoError(t, err)
	body := bytes.NewBuffer(requestBytes)

	classURL := fmt.Sprintf("%s/classes/%s", serverURL, classCreated.ID)
	req, err := http.NewRequest(http.MethodPatch, classURL, body)
	require.NoError(t, err)

	resp, err := httpClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var class classes.Class
	err = json.Unmarshal(respBody, &class)
	require.NoError(t, err)

	assert.NotEmpty(t, class.ID)
}

func TestHandler_UpdateClass_NotFound(t *testing.T) {
	serverURL, httpClient := setupIntegration(t)

	newName := uuid.NewString()
	updateClass := classes.UpdateClass{
		Name: &newName,
	}
	requestBytes, err := json.Marshal(updateClass)
	require.NoError(t, err)
	body := bytes.NewBuffer(requestBytes)

	classURL := fmt.Sprintf("%s/classes/%s", serverURL, uuid.NewString())
	req, err := http.NewRequest(http.MethodPatch, classURL, body)
	require.NoError(t, err)

	resp, err := httpClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestHandler_DeleteClass(t *testing.T) {
	serverURL, httpClient := setupIntegration(t)
	url := fmt.Sprintf("%s/classes", serverURL)

	newClass := classes.NewClass{
		Name:      uuid.NewString(),
		StartDate: time.Now().UTC(),
		EndDate:   time.Now().UTC().Add(time.Hour * 24 * 10),
		Capacity:  30,
	}
	classCreated := CreateNewClass(t, httpClient, url, newClass)
	assert.NotEmpty(t, classCreated.ID)

	classURL := fmt.Sprintf("%s/classes/%s", serverURL, classCreated.ID)
	req, err := http.NewRequest(http.MethodDelete, classURL, nil)
	require.NoError(t, err)

	resp, err := httpClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestHandler_DeleteNonExistingClass(t *testing.T) {
	serverURL, httpClient := setupIntegration(t)

	classURL := fmt.Sprintf("%s/classes/%s", serverURL, uuid.NewString())
	req, err := http.NewRequest(http.MethodDelete, classURL, nil)
	require.NoError(t, err)

	resp, err := httpClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestHandler_ListClasses(t *testing.T) {
	serverURL, httpClient := setupIntegration(t)
	url := fmt.Sprintf("%s/classes", serverURL)
	newClass := classes.NewClass{
		Name:      uuid.NewString(),
		StartDate: time.Now().UTC(),
		EndDate:   time.Now().UTC().Add(time.Hour * 24 * 10),
		Capacity:  30,
	}
	CreateNewClass(t, httpClient, url, newClass)
	CreateNewClass(t, httpClient, url, newClass)

	resp, err := httpClient.Get(url)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var allClasses []classes.Class
	err = json.Unmarshal(respBody, &allClasses)
	require.NoError(t, err)

	assert.NotEmpty(t, allClasses)
	assert.Equal(t, 2, len(allClasses))
}

func TestHandler_ListClasses_EmptyList(t *testing.T) {
	serverURL, httpClient := setupIntegration(t)
	url := fmt.Sprintf("%s/classes", serverURL)

	resp, err := httpClient.Get(url)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var allClasses []classes.Class
	err = json.Unmarshal(respBody, &allClasses)
	require.NoError(t, err)

	assert.Empty(t, allClasses)
	assert.Equal(t, 0, len(allClasses))
}

func CreateNewClass(t *testing.T, httpClient *http.Client, url string, newClass classes.NewClass) classes.Class {
	requestBytes, err := json.Marshal(newClass)
	require.NoError(t, err)
	body := bytes.NewBuffer(requestBytes)

	resp, err := httpClient.Post(url, "application/json", body)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var class classes.Class
	err = json.Unmarshal(respBody, &class)
	require.NoError(t, err)
	return class
}
