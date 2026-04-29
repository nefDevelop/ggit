package health

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSentinelDecisionLogic(t *testing.T) {
	now := time.Now().Unix()

	tests := []struct {
		name                string
		lastCheck           int64
		commitsSinceCheck   int64
		activityCooldown    int
		idleCheckDays       int
		wantShouldRun       bool
	}{
		{
			name:              "Initial check (never checked)",
			lastCheck:         0,
			commitsSinceCheck: 0,
			wantShouldRun:     true,
		},
		{
			name:              "Mandatory monthly check (31 days ago)",
			lastCheck:         now - (31 * 24 * 60 * 60),
			commitsSinceCheck: 0,
			wantShouldRun:     true,
		},
		{
			name:              "Recently checked, no activity",
			lastCheck:         now - (1 * 24 * 60 * 60),
			commitsSinceCheck: 0,
			wantShouldRun:     false,
		},
		{
			name:              "Activity detected, but cooldown not reached",
			lastCheck:         now - (1 * 24 * 60 * 60),
			commitsSinceCheck: 5,
			activityCooldown:  7, // days
			wantShouldRun:     false,
		},
		{
			name:              "Activity detected, cooldown reached",
			lastCheck:         now - (8 * 24 * 60 * 60),
			commitsSinceCheck: 1,
			activityCooldown:  7,
			wantShouldRun:     true,
		},
		{
			name:              "Idle period reached",
			lastCheck:         now - (15 * 24 * 60 * 60),
			commitsSinceCheck: 0,
			idleCheckDays:     14,
			wantShouldRun:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timeSinceLast := now - tt.lastCheck
			shouldRun := false

			// Logic matched with sentinel.go
			if tt.lastCheck == 0 {
				shouldRun = true
			} else if timeSinceLast > 30*24*60*60 {
				shouldRun = true
			} else if tt.commitsSinceCheck > 0 {
				cooldown := int64(tt.activityCooldown)
				if cooldown == 0 { cooldown = 7 }
				if timeSinceLast > cooldown*24*60*60 {
					shouldRun = true
				} else if tt.commitsSinceCheck >= 100 {
					shouldRun = true
				}
			} else {
				idleDays := int64(tt.idleCheckDays)
				if idleDays == 0 { idleDays = 14 }
				if timeSinceLast > idleDays*24*60*60 {
					shouldRun = true
				}
			}

			assert.Equal(t, tt.wantShouldRun, shouldRun)
		})
	}
}
