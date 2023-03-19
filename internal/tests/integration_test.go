//go:build integration
// +build integration

package tests

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	httpserver "github.com/ezhdanovskiy/companies/internal/http"
	"github.com/ezhdanovskiy/companies/internal/http/requests"
	"github.com/ezhdanovskiy/companies/internal/repository"
	"github.com/ezhdanovskiy/companies/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"go.uber.org/zap"
)

const logsEnabled = false

func TestCreateCompany(t *testing.T) {
	ts := newTestService(t)
	defer ts.Finish()

	uid := uuid.New()
	company, err := ts.repo.GetCompany(context.Background(), uid.String())
	require.NoError(t, err)
	require.Nil(t, company)

	req := requests.CreateCompany{
		ID:              uid.String(),
		Name:            "Name-" + uid.String()[:10],
		EmployeesAmount: 17,
		Registered:      true,
		Type:            "Cooperative",
	}

	code, body := ts.doRequest(http.MethodPost, "/companies", req)
	assert.Equal(t, http.StatusCreated, code)
	assert.Equal(t, "null", body)

	company, err = ts.repo.GetCompany(context.Background(), uid.String())
	require.NoError(t, err)
	require.NotNil(t, company)
	assert.EqualValues(t, req.Name, company.Name)

	ts.cleanCompanies(req.ID)
}

// TestServer ---------------------------------------------------------------------------------------------------------
type TestServer struct {
	t      *testing.T
	log    *zap.SugaredLogger
	db     *bun.DB
	repo   *repository.Repo
	svc    *service.Service
	router *gin.Engine
}

func newTestService(t *testing.T) TestServer {
	t.Parallel()

	var log *zap.SugaredLogger
	if logsEnabled {
		logger, _ := zap.NewDevelopment()
		log = logger.Sugar()
	} else {
		log = zap.NewNop().Sugar()
	}

	dsn := "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	pgdb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(pgdb, pgdialect.New())

	require.NoError(t, db.Ping())

	if logsEnabled {
		// Print all queries to stdout.
		db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.TestMode)
	}

	err := repository.MigrateUp(log, db.DB, "file://../../migrations")
	require.NoError(t, err)

	repo, err := repository.NewRepo(log, db)
	require.NoError(t, err)

	svc := service.NewService(log, repo)
	srv := httpserver.NewServer(log, 0, svc)
	//router := gin.Default()
	router := gin.New()

	srv.SetAPIV1Routes(router.Group("/"))

	ts := TestServer{
		t:      t,
		log:    log,
		db:     db,
		repo:   repo,
		svc:    service.NewService(log, repo),
		router: router,
	}

	return ts
}

func (ts *TestServer) doRequest(method, target string, body interface{}) (code int, respBody string) {
	b := new(bytes.Buffer)
	if str, ok := body.(string); ok {
		b.WriteString(str)
	} else {
		err := json.NewEncoder(b).Encode(body)
		require.NoError(ts.t, err)
	}

	req := httptest.NewRequest(method, target, b)

	recorder := httptest.NewRecorder()
	ts.router.ServeHTTP(recorder, req)

	return recorder.Code, recorder.Body.String()
}

func (ts *TestServer) cleanCompanies(uuids ...string) {
	_, err := ts.db.QueryContext(context.Background(), "DELETE FROM companies WHERE id IN (?)", bun.In(uuids))
	require.NoError(ts.t, err)
}

func (ts *TestServer) Finish() {
}
