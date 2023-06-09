package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/daniel-oliveiravas/class-booking-service/business/bookings"
	"github.com/daniel-oliveiravas/class-booking-service/business/classes"
	"github.com/daniel-oliveiravas/class-booking-service/business/members"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler_BookClass(t *testing.T) {
	serverURL, httpClient := setupIntegration(t)
	url := fmt.Sprintf("%s/bookings", serverURL)

	class, member := PrepareToBookClass(t, httpClient, serverURL)

	bookClass := bookings.BookClass{
		MemberID:  member.ID,
		ClassID:   class.ID,
		ClassDate: time.Now().UTC(),
	}
	booking := BookClass(t, httpClient, url, bookClass)

	assert.NotEmpty(t, booking.ID)
	assert.NotEmpty(t, booking.BookedAt)
	assert.NotEmpty(t, booking.UpdatedAt)
	assert.Equal(t, class.ID, booking.ClassID)
	assert.Equal(t, member.ID, booking.MemberID)
	assert.Equal(t, bookClass.ClassDate.Truncate(time.Hour*24), booking.ClassDate)
}

func TestHandler_BookClass_InvalidData(t *testing.T) {
	serverURL, httpClient := setupIntegration(t)
	url := fmt.Sprintf("%s/bookings", serverURL)

	class, member := PrepareToBookClass(t, httpClient, serverURL)

	bookClass := bookings.BookClass{
		MemberID:  member.ID,
		ClassID:   class.ID,
		ClassDate: time.Now().UTC().Add(time.Hour * 24 * 30 * -1),
	}
	requestBytes, err := json.Marshal(bookClass)
	require.NoError(t, err)
	body := bytes.NewBuffer(requestBytes)

	resp, err := httpClient.Post(url, "application/json", body)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var booking bookings.Booking
	err = json.Unmarshal(respBody, &booking)
	require.NoError(t, err)

	assert.Empty(t, booking.ID)
}

func TestHandler_GetBooking(t *testing.T) {
	serverURL, httpClient := setupIntegration(t)
	url := fmt.Sprintf("%s/bookings", serverURL)

	class, member := PrepareToBookClass(t, httpClient, serverURL)

	bookClass := bookings.BookClass{
		MemberID:  member.ID,
		ClassID:   class.ID,
		ClassDate: time.Now().UTC(),
	}
	classBooked := BookClass(t, httpClient, url, bookClass)
	assert.NotEmpty(t, classBooked.ID)

	getBookingURL := fmt.Sprintf("%s/bookings/%s", serverURL, classBooked.ID)
	resp, err := httpClient.Get(getBookingURL)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var booking bookings.Booking
	err = json.Unmarshal(respBody, &booking)
	require.NoError(t, err)
	assert.NotEmpty(t, booking.ID)
}

func TestHandler_GetBooking_NotFound(t *testing.T) {
	serverURL, httpClient := setupIntegration(t)

	getBookingURL := fmt.Sprintf("%s/bookings/%s", serverURL, uuid.NewString())
	resp, err := httpClient.Get(getBookingURL)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestHandler_DeleteBooking(t *testing.T) {
	serverURL, httpClient := setupIntegration(t)
	url := fmt.Sprintf("%s/bookings", serverURL)

	class, member := PrepareToBookClass(t, httpClient, serverURL)

	bookClass := bookings.BookClass{
		MemberID:  member.ID,
		ClassID:   class.ID,
		ClassDate: time.Now().UTC(),
	}
	classBooked := BookClass(t, httpClient, url, bookClass)
	assert.NotEmpty(t, classBooked.ID)

	bookingURL := fmt.Sprintf("%s/bookings/%s", serverURL, classBooked.ID)
	req, err := http.NewRequest(http.MethodDelete, bookingURL, nil)
	require.NoError(t, err)

	resp, err := httpClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestHandler_DeleteNonExistingBooking(t *testing.T) {
	serverURL, httpClient := setupIntegration(t)

	classURL := fmt.Sprintf("%s/bookings/%s", serverURL, uuid.NewString())
	req, err := http.NewRequest(http.MethodDelete, classURL, nil)
	require.NoError(t, err)

	resp, err := httpClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestHandler_ListBookings(t *testing.T) {
	serverURL, httpClient := setupIntegration(t)
	url := fmt.Sprintf("%s/bookings", serverURL)
	class, member := PrepareToBookClass(t, httpClient, serverURL)

	newClass := bookings.BookClass{
		MemberID:  member.ID,
		ClassID:   class.ID,
		ClassDate: time.Now().UTC(),
	}
	BookClass(t, httpClient, url, newClass)

	class, member = PrepareToBookClass(t, httpClient, serverURL)

	newClass = bookings.BookClass{
		MemberID:  member.ID,
		ClassID:   class.ID,
		ClassDate: time.Now().UTC(),
	}
	BookClass(t, httpClient, url, newClass)

	resp, err := httpClient.Get(url)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var allBookings []bookings.Booking
	err = json.Unmarshal(respBody, &allBookings)
	require.NoError(t, err)

	assert.NotEmpty(t, allBookings)
	assert.Equal(t, 2, len(allBookings))
}

func TestHandler_ListBookings_EmptyList(t *testing.T) {
	serverURL, httpClient := setupIntegration(t)
	url := fmt.Sprintf("%s/bookings", serverURL)

	resp, err := httpClient.Get(url)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var allBookings []bookings.Booking
	err = json.Unmarshal(respBody, &allBookings)
	require.NoError(t, err)

	assert.Empty(t, allBookings)
	assert.Equal(t, 0, len(allBookings))
}

func PrepareToBookClass(t *testing.T, httpClient *http.Client, url string) (classes.Class, members.Member) {
	newClass := classes.NewClass{
		Name:      uuid.NewString(),
		StartDate: time.Now().UTC(),
		EndDate:   time.Now().UTC().Add(time.Hour * 24 * 10),
		Capacity:  30,
	}

	class := CreateNewClass(t, httpClient, fmt.Sprintf("%s/classes", url), newClass)

	newMember := members.NewMember{
		Name: uuid.NewString(),
	}

	member := CreateNewMember(t, httpClient, fmt.Sprintf("%s/members", url), newMember)

	return class, member
}

func BookClass(t *testing.T, httpClient *http.Client, url string, bookClass bookings.BookClass) bookings.Booking {
	requestBytes, err := json.Marshal(bookClass)
	require.NoError(t, err)
	body := bytes.NewBuffer(requestBytes)

	resp, err := httpClient.Post(url, "application/json", body)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var class bookings.Booking
	err = json.Unmarshal(respBody, &class)
	require.NoError(t, err)
	return class
}
