package frontend

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/signadot/hotrod/pkg/notifications"
)

func TestFrontendAPI(t *testing.T) {
	frontendAddr := os.Getenv("TEST_FRONTEND_ADDR")
	targetWorkload := os.Getenv("TEST_TARGET_WORKLOAD")
	sandboxName := os.Getenv("TEST_SANDBOX_NAME")
	if frontendAddr == "" || targetWorkload == "" {
		t.Skip()
		return
	}

	sessionID := rand.Intn(10000)
	reqID, err := sendBaselineDispatch(frontendAddr, sessionID)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("sent dispatch request, sessionID=%d", sessionID)
	err = waitNotification(t, frontendAddr, targetWorkload, sandboxName, sessionID, int(reqID), 20*time.Second)
	if err != nil {
		t.Error(err)
		return
	}
}

func sendBaselineDispatch(frontendAddr string, sessionID int) (uint, error) {
	dispatch := DispatchRequest{
		SessionID:         uint(sessionID),
		RequestID:         uint(rand.Intn(10000)),
		PickupLocationID:  123,
		DropoffLocationID: 567,
	}
	d, err := json.Marshal(dispatch)
	if err != nil {
		return 0, err
	}
	buf := new(bytes.Buffer)
	buf.Write(d)
	url := fmt.Sprintf("http://%s/dispatch", frontendAddr)
	req, err := http.NewRequest(http.MethodPost, url, buf)
	if err != nil {
		return 0, err
	}
	client := &http.Client{
		Transport: &http.Transport{},
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("dispatch failed with status %d", resp.StatusCode)
	}
	return dispatch.RequestID, nil
}

func getNotification(frontendAddr string, sessionID, cursor int, n *notifications.NotificationList) error {
	url := fmt.Sprintf("http://%s/notifications?cursor=%d&sessionID=%d", frontendAddr, cursor, sessionID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	client := &http.Client{Transport: &http.Transport{}}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("notififications GET failed with status %d", resp.StatusCode)
	}
	if err := json.NewDecoder(resp.Body).Decode(n); err != nil {
		return err
	}
	return nil
}

func waitNotification(t *testing.T, frontendAddr, targetWorkload, sandboxName string,
	sessionID, reqID int, dur time.Duration) error {
	end := time.After(dur)
	ticker := time.NewTicker(time.Second / 1)
	defer ticker.Stop()

	found := false
	cursor := -1
	for {
		n := notifications.NotificationList{}

		if err := getNotification(frontendAddr, sessionID, cursor, &n); err != nil {
			return err
		}
		for i := range n.Notifications {
			// wait until the driver has been dispatched
			nn := &n.Notifications[i]
			t.Logf("got notification: [%v] id=%s, body=%s, context=%+v",
				nn.Timestamp, nn.ID, nn.Body, nn.Context)

			// look for the target forked workload
			if nn.Context.BaselineWorkload == targetWorkload && nn.Context.SandboxName == sandboxName {
				found = true
			}
			if fmt.Sprintf("req-%d-dispatched-driver", reqID) == nn.ID {
				// driver has been dispatched
				if found {
					return nil
				} else {
					return errors.New("did not find target workload, participating in request flow")
				}
			}
		}
		cursor = n.Cursor
		select {
		case <-ticker.C:
		case <-end:
			return os.ErrDeadlineExceeded
		}
	}
}
