package record

import (
	"testing"
	"time"

	"github.com/epalmerini/abitudini/internal/testhelpers"
)

func TestRecordStore_Record(t *testing.T) {
	db := testhelpers.NewTestDB(t)
	store := NewStore(db)

	habitID := 1
	date := time.Now()

	err := store.Record(habitID, date)
	if err != nil {
		t.Fatalf("failed to record: %v", err)
	}
}
