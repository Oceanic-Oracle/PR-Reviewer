// e2e_test.go
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"pr/internal/app"
	"pr/internal/config"
	"pr/internal/dto"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	DBURL    = "postgres://User:241265@localhost:5433/postgres?sslmode=disable"
	HTTPPort = ":8081"
	HTTPHost = "http://localhost" + HTTPPort
)

func TestMain(m *testing.M) {
	if err := exec.Command("docker-compose", "-f", "compose.test.yml", "up", "-d").Run(); err != nil {
		log.Fatal("Failed to start test DB:", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := waitForDB(ctx); err != nil {
		_ = exec.Command("docker-compose", "-f", "compose.test.yml", "down", "-v").Run()
		log.Fatal("DB not ready:", err)
	}

	cfg := config.MustLoad()
	cfg.Storage.URL = DBURL
	cfg.HTTP.Addr = HTTPPort

	go func() {
		app := app.NewBootstrap(cfg)
		app.Run()
	}()

	time.Sleep(10 * time.Second)

	code := m.Run()
	_ = exec.Command("docker-compose", "-f", "compose.test.yml", "down", "-v").Run()
	os.Exit(code)
}

func waitForDB(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			conn, err := pgxpool.New(ctx, DBURL)

			if err != nil {
				time.Sleep(200 * time.Millisecond)
				continue
			}

			if err := conn.Ping(ctx); err == nil {
				conn.Close()
				return nil
			}

			conn.Close()
			time.Sleep(200 * time.Millisecond)
		}
	}
}

// =================================================================
// =================================================================

func TestScenario_CreateTeamsAndPRs(t *testing.T) {
	teams := []dto.CreateTeamRequest{
		{TeamName: "Infrastructure", Members: generateUsers(4)},
		{TeamName: "Android", Members: generateUsers(4)},
		{TeamName: "IOS", Members: generateUsers(4)},
		{TeamName: "Payments", Members: generateUsers(4)},
		{TeamName: "Frontend", Members: generateUsers(4)},
		{TeamName: "ML", Members: generateUsers(4)},
	}

	wg := &sync.WaitGroup{}

	mtx := &sync.Mutex{}

	createdTeams := make([]dto.CreateTeamResponse, len(teams))

	for i, team := range teams {
		wg.Add(1)

		go func(idx int, team dto.CreateTeamRequest) {
			defer wg.Done()

			reqBody, _ := json.Marshal(team)
			respBody, statusCode, err := Request(reqBody, http.MethodPost, "/team/add")
			if err != nil || statusCode != http.StatusCreated {
				t.Errorf("Create team failed: status=%d, err=%v", statusCode, err)
				return
			}

			var resp dto.CreateTeamResponse
			if err := json.Unmarshal(respBody, &resp); err != nil {
				t.Errorf("Failed to unmarshal team response: %v", err)
				return
			}

			mtx.Lock()
			createdTeams[idx] = resp
			mtx.Unlock()
		}(i, team)
	}
	wg.Wait()

	var prWg sync.WaitGroup
	for _, teamResp := range createdTeams {
		if len(teamResp.Team.Members) < 2 {
			continue
		}

		prWg.Add(1)
		go func(team dto.Team) {
			defer prWg.Done()

			author := team.Members[0]
			prReq := dto.CreatePRRequest{
				ID:       uuid.New().String(),
				Name:     "E2E Test PR",
				AuthorID: author.UserID,
			}

			reqBody, _ := json.Marshal(prReq)
			respBody, statusCode, err := Request(reqBody, http.MethodPost, "/pullRequest/create")
			if err != nil || statusCode != http.StatusCreated {
				t.Errorf("Create PR failed for team %s: status=%d, err=%v", team.TeamName, statusCode, err)
				return
			}

			var prResp dto.CreatePRResponse
			if err := json.Unmarshal(respBody, &prResp); err != nil {
				t.Errorf("Failed to unmarshal PR response: %v", err)
				return
			}

			if len(prResp.PR.AssignedReviewers) == 0 {
				t.Errorf("PR must have at least 1 reviewer")
			}
			if len(prResp.PR.AssignedReviewers) > 2 {
				t.Errorf("PR must have at most 2 reviewers")
			}
		}(teamResp.Team)
	}
	prWg.Wait()
}

// =================================================================
// =================================================================

var counter int64

func generateUsers(n int) []dto.User {
	users := make([]dto.User, n)
	for i := 0; i < n; i++ {
		id := atomic.AddInt64(&counter, 1)
		users[i] = dto.User{
			UserID:   uuid.New().String(),
			Username: fmt.Sprintf("testuser-%d", id),
			IsActive: true,
		}
	}

	return users
}

func Request(reqBody []byte, method, path string) (resBody []byte, statusCode int, err error) {
	req, err := http.NewRequest(method, HTTPHost+path, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, resp.StatusCode, nil
}
