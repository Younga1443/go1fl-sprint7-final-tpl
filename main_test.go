package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)

		require.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(response, req)
		require.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, http.StatusOK, response.Code)
	}
}

func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		want    int
		request string
	}{
		{0, "/cafe?city=tula&count=0"},
		{1, "/cafe?city=tula&count=1"},
		{2, "/cafe?city=tula&count=2"},
		{min(len(cafeList["tula"]), 100), "/cafe?city=tula&count=100"},
	}

	for _, n := range requests {
		var r []string
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", n.request, nil)

		handler.ServeHTTP(response, req)
		body := strings.TrimSpace(response.Body.String())
		if body != "" {
			r = strings.Split(body, ",")
		}
		require.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, n.want, len(r))
	}
}

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		wantCount int
		search    string
		request   string
	}{
		{0, "фасоль", "/cafe?city=moscow&search=фасоль"},
		{2, "кофе", "/cafe?city=moscow&search=кофе"},
		{1, "вилка", "/cafe?city=moscow&search=вилка"},
	}

	for _, n := range requests {
		var r []string
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", n.request, nil)

		handler.ServeHTTP(response, req)
		body := strings.TrimSpace(response.Body.String())
		if body != "" {
			r = strings.Split(body, ",")
		}
		for _, cafeName := range r {
			assert.True(t, strings.Contains(strings.ToLower(cafeName), strings.ToLower(n.search)))
		}
		require.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, n.wantCount, len(r))
	}
}
