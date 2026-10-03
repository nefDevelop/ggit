package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Webhook embeds *HookEvent, which carries the event flags, and every Has*Event
// helper reads through it. Rows written with an empty "events" column used to
// come back with a nil embedded pointer, and listing the webhooks then panicked
// with a nil dereference instead of returning them.
func TestWebhookAfterFindAlwaysInitialisesHookEvent(t *testing.T) {
	for _, tc := range []struct {
		name   string
		events string
		want   []string
	}{
		{name: "empty events column", events: "", want: []string{}},
		{name: "populated events column", events: `{"choose_events":true,"events":{"push":true}}`, want: []string{"push"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			w := &Webhook{Events: tc.events}
			require.NoError(t, w.AfterFind(nil))

			require.NotNil(t, w.HookEvent, "HookEvent must never be nil after a load")

			// Reading through the embedded pointer is what used to panic.
			assert.NotPanics(t, func() {
				_ = w.HasPushEvent()
				_ = w.EventsArray()
			})
			assert.Equal(t, tc.want, w.EventsArray())
		})
	}
}

// The events column is what survives a round trip, so it has to be written on
// the way in. It was not, which is why a webhook answered with its events and
// then failed to list.
func TestWebhookUpdateEventPersistsTheFlags(t *testing.T) {
	w := &Webhook{HookEvent: &HookEvent{
		ChooseEvents: true,
		HookEvents:   HookEvents{Push: true, Issues: true},
	}}

	require.NoError(t, w.UpdateEvent())
	assert.Contains(t, w.Events, `"push":true`)
	assert.Contains(t, w.Events, `"issues":true`)

	// And they must survive the trip back.
	reloaded := &Webhook{Events: w.Events}
	require.NoError(t, reloaded.AfterFind(nil))
	assert.Equal(t, []string{"push", "issues"}, reloaded.EventsArray())
}
