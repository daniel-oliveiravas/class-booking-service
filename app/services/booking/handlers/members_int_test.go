package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/daniel-oliveiravas/class-booking-service/app/services/booking/handlers"
	"github.com/daniel-oliveiravas/class-booking-service/business/bookings"
	pgbookings "github.com/daniel-oliveiravas/class-booking-service/business/bookings/integration/postgres"
	"github.com/daniel-oliveiravas/class-booking-service/business/classes"
	pgclasses "github.com/daniel-oliveiravas/class-booking-service/business/classes/integration/postgres"
	"github.com/daniel-oliveiravas/class-booking-service/business/members"
	pgmembers "github.com/daniel-oliveiravas/class-booking-service/business/members/integration/postgres"
	"github.com/daniel-oliveiravas/class-booking-service/foundation/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func setupIntegration(t *testing.T) (string, *http.Client) {
	if os.Getenv("INTEGRATION") == "" {
		t.Skip("skipping integration tests: set INTEGRATION environment variable")
	}
	ctx := context.Background()
	schema := t.Name()
	pgCfg := postgres.Config{
		Host:             "localhost",
		Port:             5432,
		DatabaseUser:     "class_booking",
		DatabasePassword: "class_booking",
		DatabaseName:     "class_booking_qa",
		SSLMode:          "none",
		SearchPath:       schema,
	}
	db, err := postgres.Open(ctx, pgCfg)
	require.NoError(t, err)

	err = postgres.DropAndCreateSchema(ctx, db, schema)
	require.NoError(t, err)

	err = postgres.Migrate("file://../../../../scripts/db/migrations/", pgCfg)
	require.NoError(t, err)

	logger := zap.NewNop().Sugar()
	membersRepo := pgmembers.NewMembersRepository(logger, db)
	membersUsecase := members.NewUsecase(membersRepo)

	classesRepo := pgclasses.NewClassesRepository(logger, db)
	classesUsecase := classes.NewUsecase(classesRepo)

	bookingsRepo := pgbookings.NewBookingsRepository(logger, db)
	bookingsUsecase := bookings.NewUsecase(bookingsRepo, membersUsecase, classesUsecase)

	cfg := handlers.Config{
		MembersUsecase: membersUsecase,
		ClassesUsecase: classesUsecase,
		BookingUsecase: bookingsUsecase,
		Logger:         logger,
	}
	handlersAPI, err := handlers.NewHandler(cfg)
	require.NoError(t, err)

	server := httptest.NewServer(handlersAPI.API())
	httpClient := server.Client()

	return server.URL, httpClient
}

func TestHandler_AddMember(t *testing.T) {
	serverURL, httpClient := setupIntegration(t)
	url := fmt.Sprintf("%s/members", serverURL)
	newMember := members.NewMember{
		Name: uuid.NewString(),
	}
	member := CreateNewMember(t, httpClient, url, newMember)

	assert.NotEmpty(t, member.ID)
	assert.NotEmpty(t, member.CreatedAt)
	assert.NotEmpty(t, member.UpdatedAt)
	assert.Equal(t, newMember.Name, member.Name)
}

func TestHandler_AddMember_InvalidData(t *testing.T) {

	serverURL, httpClient := setupIntegration(t)
	url := fmt.Sprintf("%s/members", serverURL)

	newMember := members.NewMember{
		Name: "",
	}
	requestBytes, err := json.Marshal(newMember)
	require.NoError(t, err)
	body := bytes.NewBuffer(requestBytes)

	resp, err := httpClient.Post(url, "application/json", body)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var member members.Member
	err = json.Unmarshal(respBody, &member)
	require.NoError(t, err)

	assert.Empty(t, member.ID)
}

func TestHandler_GetMember(t *testing.T) {
	serverURL, httpClient := setupIntegration(t)
	url := fmt.Sprintf("%s/members", serverURL)

	newMember := members.NewMember{
		Name: uuid.NewString(),
	}
	memberCreated := CreateNewMember(t, httpClient, url, newMember)
	assert.NotEmpty(t, memberCreated.ID)

	getMemberURL := fmt.Sprintf("%s/members/%s", serverURL, memberCreated.ID)
	resp, err := httpClient.Get(getMemberURL)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var member members.Member
	err = json.Unmarshal(respBody, &member)
	require.NoError(t, err)

	assert.NotEmpty(t, member.ID)
}

func TestHandler_GetMember_NotFound(t *testing.T) {
	serverURL, httpClient := setupIntegration(t)

	getMemberURL := fmt.Sprintf("%s/members/%s", serverURL, uuid.NewString())
	resp, err := httpClient.Get(getMemberURL)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestHandler_UpdateMember(t *testing.T) {
	serverURL, httpClient := setupIntegration(t)
	url := fmt.Sprintf("%s/members", serverURL)

	newMember := members.NewMember{
		Name: uuid.NewString(),
	}
	memberCreated := CreateNewMember(t, httpClient, url, newMember)
	assert.NotEmpty(t, memberCreated.ID)

	newName := uuid.NewString()
	updateMember := members.UpdateMember{
		Name: &newName,
	}
	requestBytes, err := json.Marshal(updateMember)
	require.NoError(t, err)
	body := bytes.NewBuffer(requestBytes)

	memberURL := fmt.Sprintf("%s/members/%s", serverURL, memberCreated.ID)
	req, err := http.NewRequest(http.MethodPatch, memberURL, body)
	require.NoError(t, err)

	resp, err := httpClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var member members.Member
	err = json.Unmarshal(respBody, &member)
	require.NoError(t, err)

	assert.NotEmpty(t, member.ID)
}

func TestHandler_UpdateMember_NotFound(t *testing.T) {
	serverURL, httpClient := setupIntegration(t)

	newName := uuid.NewString()
	updateMember := members.UpdateMember{
		Name: &newName,
	}
	requestBytes, err := json.Marshal(updateMember)
	require.NoError(t, err)
	body := bytes.NewBuffer(requestBytes)

	memberURL := fmt.Sprintf("%s/members/%s", serverURL, uuid.NewString())
	req, err := http.NewRequest(http.MethodPatch, memberURL, body)
	require.NoError(t, err)

	resp, err := httpClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestHandler_DeleteMember(t *testing.T) {
	serverURL, httpClient := setupIntegration(t)
	url := fmt.Sprintf("%s/members", serverURL)

	newMember := members.NewMember{
		Name: uuid.NewString(),
	}
	memberCreated := CreateNewMember(t, httpClient, url, newMember)
	assert.NotEmpty(t, memberCreated.ID)

	newName := uuid.NewString()
	updateMember := members.UpdateMember{
		Name: &newName,
	}
	requestBytes, err := json.Marshal(updateMember)
	require.NoError(t, err)
	body := bytes.NewBuffer(requestBytes)

	memberURL := fmt.Sprintf("%s/members/%s", serverURL, memberCreated.ID)
	req, err := http.NewRequest(http.MethodDelete, memberURL, body)
	require.NoError(t, err)

	resp, err := httpClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestHandler_DeleteNonExistingMember(t *testing.T) {
	serverURL, httpClient := setupIntegration(t)

	newName := uuid.NewString()
	updateMember := members.UpdateMember{
		Name: &newName,
	}
	requestBytes, err := json.Marshal(updateMember)
	require.NoError(t, err)
	body := bytes.NewBuffer(requestBytes)

	memberURL := fmt.Sprintf("%s/members/%s", serverURL, uuid.NewString())
	req, err := http.NewRequest(http.MethodDelete, memberURL, body)
	require.NoError(t, err)

	resp, err := httpClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestHandler_ListMember(t *testing.T) {
	serverURL, httpClient := setupIntegration(t)
	url := fmt.Sprintf("%s/members", serverURL)
	newMember := members.NewMember{
		Name: uuid.NewString(),
	}
	CreateNewMember(t, httpClient, url, newMember)
	CreateNewMember(t, httpClient, url, newMember)

	resp, err := httpClient.Get(url)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var allMembers []members.Member
	err = json.Unmarshal(respBody, &allMembers)
	require.NoError(t, err)

	assert.NotEmpty(t, allMembers)
	assert.Equal(t, 2, len(allMembers))
}

func TestHandler_ListMember_EmptyList(t *testing.T) {
	serverURL, httpClient := setupIntegration(t)
	url := fmt.Sprintf("%s/members", serverURL)

	resp, err := httpClient.Get(url)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var allMembers []members.Member
	err = json.Unmarshal(respBody, &allMembers)
	require.NoError(t, err)

	assert.Empty(t, allMembers)
	assert.Equal(t, 0, len(allMembers))
}

func CreateNewMember(t *testing.T, httpClient *http.Client, url string, newMember members.NewMember) members.Member {
	requestBytes, err := json.Marshal(newMember)
	require.NoError(t, err)
	body := bytes.NewBuffer(requestBytes)

	resp, err := httpClient.Post(url, "application/json", body)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var member members.Member
	err = json.Unmarshal(respBody, &member)
	require.NoError(t, err)
	return member
}
