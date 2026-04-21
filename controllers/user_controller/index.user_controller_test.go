package user_controller_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"gin-gonic-gorm/controllers/user_controller"
	"gin-gonic-gorm/database"
	"gin-gonic-gorm/requests"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestDB() (sqlmock.Sqlmock, error) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		return nil, err
	}

	dialector := mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	})

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	database.DB = db
	return mock, nil
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	return r
}

func TestGetAllUser_Success(t *testing.T) {
	mock, err := setupTestDB()
	assert.NoError(t, err)

	r := setupRouter()
	r.GET("/users", user_controller.GetAllUser)

	rows := sqlmock.NewRows([]string{"id", "name", "email", "address", "born_date"}).
		AddRow(1, "John Doe", "john@test.com", "Address 1", time.Now())

	mock.ExpectQuery("SELECT \\* FROM `users`").WillReturnRows(rows)

	req, _ := http.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "John Doe")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByID_Success(t *testing.T) {
	mock, err := setupTestDB()
	assert.NoError(t, err)

	r := setupRouter()
	r.GET("/users/:id", user_controller.GetByID)

	rows := sqlmock.NewRows([]string{"id", "name", "email", "address"}).
		AddRow(1, "John Doe", "john@test.com", "Address 1")

	mock.ExpectQuery("SELECT \\* FROM `users` WHERE id = \\? ORDER BY `users`.`id` LIMIT \\?").
		WithArgs("1", 1).
		WillReturnRows(rows)

	req, _ := http.NewRequest("GET", "/users/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "data transmitted")
	assert.Contains(t, w.Body.String(), "John Doe")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByID_NotFound(t *testing.T) {
	mock, err := setupTestDB()
	assert.NoError(t, err)

	r := setupRouter()
	r.GET("/users/:id", user_controller.GetByID)

	mock.ExpectQuery("SELECT \\* FROM `users` WHERE id = \\?").
		WithArgs("999", 1).
		WillReturnError(gorm.ErrRecordNotFound)

	req, _ := http.NewRequest("GET", "/users/999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "data not found")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStore_Success(t *testing.T) {
	mock, err := setupTestDB()
	assert.NoError(t, err)

	r := setupRouter()
	r.POST("/users", user_controller.Store)

	bornDate := time.Now()
	userReq := requests.UserRequest{
		Name:     "Jane Doe",
		Email:    "jane@test.com",
		Address:  "Address 2",
		BornDate: bornDate,
	}
	body, _ := json.Marshal(userReq)

	// Expect email check
	mock.ExpectQuery(".*SELECT \\* FROM `users` WHERE email = \\?.*").
		WithArgs("jane@test.com", 1).
		WillReturnError(gorm.ErrRecordNotFound)

	// Expect insert
	mock.ExpectBegin()
	mock.ExpectExec(".*INSERT INTO `users`.*").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "data saved successfully")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateByID_Success(t *testing.T) {
	mock, err := setupTestDB()
	assert.NoError(t, err)

	r := setupRouter()
	r.PUT("/users/:id", user_controller.UpdateByID)

	bornDate := time.Now()
	userReq := requests.UserRequest{
		Name:     "Jane Updated",
		Email:    "jane_new@test.com",
		Address:  "Address Updated",
		BornDate: bornDate,
	}
	body, _ := json.Marshal(userReq)

	// Expect find user
	rows := sqlmock.NewRows([]string{"id", "name", "email", "address", "born_date"}).
		AddRow(1, "Jane Doe", "jane@test.com", "Address 2", bornDate)
	mock.ExpectQuery("SELECT \\* FROM `users` WHERE id = \\?").
		WithArgs("1").
		WillReturnRows(rows)

	// Expect find email
	emptyRows := sqlmock.NewRows([]string{"id", "name", "email", "address", "born_date"})
	mock.ExpectQuery(".*SELECT \\* FROM `users` WHERE email = \\?.*").
		WithArgs("jane_new@test.com").
		WillReturnRows(emptyRows)

	// Expect update
	mock.ExpectBegin()
	mock.ExpectExec(".*UPDATE `users` SET.*").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	req, _ := http.NewRequest("PUT", "/users/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "data updated successfully")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteByID_Success(t *testing.T) {
	mock, err := setupTestDB()
	assert.NoError(t, err)

	r := setupRouter()
	r.DELETE("/users/:id", user_controller.DeleteByID)

	// Expect find user
	rows := sqlmock.NewRows([]string{"id", "name", "email", "address", "born_date"}).
		AddRow(1, "Jane Doe", "jane@test.com", "Address 2", time.Now())
	mock.ExpectQuery("SELECT \\* FROM `users` WHERE id = \\?").
		WithArgs("1").
		WillReturnRows(rows)

	// Expect delete
	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM `users` WHERE id = \\?").
		WithArgs("1").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	req, _ := http.NewRequest("DELETE", "/users/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "data deleted successfully")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserPaginate_Success(t *testing.T) {
	mock, err := setupTestDB()
	assert.NoError(t, err)

	r := setupRouter()
	r.GET("/users-paginate", user_controller.GetUserPaginate)

	rows := sqlmock.NewRows([]string{"id", "name", "email", "address", "born_date"}).
		AddRow(1, "User 1", "u1@test.com", "A1", time.Now()).
		AddRow(2, "User 2", "u2@test.com", "A2", time.Now())

	mock.ExpectQuery("SELECT \\* FROM `users` LIMIT \\?").
		WithArgs(10).
		WillReturnRows(rows)

	req, _ := http.NewRequest("GET", "/users-paginate?page=1&perPage=10", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "User 1")
	assert.Contains(t, w.Body.String(), "User 2")
	assert.NoError(t, mock.ExpectationsWereMet())
}
