package synchronization_test

import (
	"testing"

	"nulink/core/internal/cltest"
	"nulink/core/services/synchronization"
	"nulink/core/store/models"
	"nulink/core/store/orm"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStatsPusher(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	wsserver, wscleanup := cltest.NewEventWebSocketServer(t)
	defer wscleanup()

	pusher := synchronization.NewStatsPusher(store.ORM, wsserver.URL, "", "")
	pusher.Start()
	defer pusher.Close()

	require.NoError(t, store.ORM.RawDB(func(db *gorm.DB) error { return db.Create(&models.SyncEvent{}).Error }))
	pusher.PushNow()

	assert.Equal(t, 1, lenSyncEvents(t, store.ORM), "jobrun sync event should be created")
	cltest.CallbackOrTimeout(t, "ws server receives jobrun creation", func() {
		<-wsserver.Received
		err := wsserver.Broadcast(`{"status": 201}`)
		assert.NoError(t, err)
	})
	cltest.WaitForSyncEventCount(t, store.ORM, 0)
}

func TestStatsPusher_ClockTrigger(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	wsserver, wscleanup := cltest.NewEventWebSocketServer(t)
	defer wscleanup()

	clock := cltest.NewTriggerClock(t)
	pusher := synchronization.NewStatsPusher(store.ORM, wsserver.URL, "", "", clock)
	pusher.Start()
	defer pusher.Close()

	err := store.ORM.RawDB(func(db *gorm.DB) error {
		return db.Save(&models.SyncEvent{Body: string("")}).Error
	})
	require.NoError(t, err)

	clock.Trigger()
	cltest.CallbackOrTimeout(t, "ws server receives jobrun update", func() {
		<-wsserver.Received
		err := wsserver.Broadcast(`{"status": 201}`)
		assert.NoError(t, err)
	})
	cltest.WaitForSyncEventCount(t, store.ORM, 0)
}

func TestStatsPusher_NoAckLeavesEvent(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	wsserver, wscleanup := cltest.NewEventWebSocketServer(t)
	defer wscleanup()

	pusher := synchronization.NewStatsPusher(store.ORM, wsserver.URL, "", "")
	pusher.Start()
	defer pusher.Close()

	require.NoError(t, store.ORM.RawDB(func(db *gorm.DB) error { return db.Create(&models.SyncEvent{}).Error }))
	pusher.PushNow()

	assert.Equal(t, 1, lenSyncEvents(t, store.ORM), "jobrun sync event should be created")
	cltest.CallbackOrTimeout(t, "ws server receives jobrun creation", func() {
		<-wsserver.Received
	})
	cltest.AssertSyncEventCountStays(t, store.ORM, 1)
}

func TestStatsPusher_BadSyncLeavesEvent(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	wsserver, wscleanup := cltest.NewEventWebSocketServer(t)
	defer wscleanup()

	clock := cltest.NewTriggerClock(t)
	pusher := synchronization.NewStatsPusher(store.ORM, wsserver.URL, "", "", clock)
	pusher.Start()
	defer pusher.Close()

	require.NoError(t, store.ORM.RawDB(func(db *gorm.DB) error { return db.Create(&models.SyncEvent{}).Error }))

	assert.Equal(t, 1, lenSyncEvents(t, store.ORM), "jobrun sync event should be created")
	clock.Trigger()
	cltest.CallbackOrTimeout(t, "ws server receives jobrun creation", func() {
		<-wsserver.Received
		err := wsserver.Broadcast(`{"status": 500}`)
		assert.NoError(t, err)
	})
	cltest.AssertSyncEventCountStays(t, store.ORM, 1)
}

func lenSyncEvents(t *testing.T, orm *orm.ORM) int {
	count, err := orm.CountOf(&models.SyncEvent{})
	require.NoError(t, err)
	return count
}
