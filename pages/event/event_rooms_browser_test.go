package event

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/live"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/delaneyj/toolbelt/embeddednats"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	natsserver "github.com/nats-io/nats-server/v2/server"
)

func TestEventRoomBrowserFixture(t *testing.T) {
	if os.Getenv("EVENT_ROOM_BROWSER_FIXTURE") != "1" { t.Skip("temporary browser fixture") }
	db := createEventRoomTestDB(t)
	seedEventVisibilityPulje(t, db, models.PuljeLordagMorgen)
	seedEventVisibilityEventPulje(t, db, "room-event", models.PuljeLordagMorgen, true)
	testutil.MustExec(t, db, `UPDATE puljer SET name = 'Lørdag morgen', start_at = '2026-10-10T10:00:00+02:00', end_at = '2026-10-10T15:00:00+02:00' WHERE id = ?`, models.PuljeLordagMorgen)
	testutil.MustExec(t, db, `UPDATE relation_event_puljer SET room_id = 42`)
	testutil.MustExec(t, db, `INSERT INTO rooms(id,name,room_number,floor,max_concurrent_games) VALUES (43,'Lucie Wolf','710',7,1)`)
	ns, err := embeddednats.New(context.Background(), embeddednats.WithNATSServerOptions(&natsserver.Options{Host:"127.0.0.1",Port:-1,JetStream:true,StoreDir:t.TempDir(),NoSigs:true}))
	if err != nil { t.Fatal(err) }; defer ns.Close(); ns.WaitForServer()
	manager, err := live.NewManager(context.Background(),ns,sessions.NewCookieStore([]byte("room-group-fixture-only-secret")))
	if err != nil { t.Fatal(err) }
	router := chi.NewRouter()
	if err := SetupEventRoute(router,manager,db,testutil.NewTestLogger(),nil); err != nil { t.Fatal(err) }
	router.Handle("/static/*",http.StripPrefix("/static/",http.FileServer(http.Dir("../../static"))))
	router.Post("/fixture/{action}",func(w http.ResponseWriter,r *http.Request){
		bucket := live.BucketRooms
		var query string
		switch chi.URLParam(r,"action") {
		case "split": query = `UPDATE relation_event_puljer SET room_id = 43 WHERE pulje_id = 'FredagKveld'`
		case "merge": query = `UPDATE relation_event_puljer SET room_id = 42`
		case "rename": query = `UPDATE rooms SET name = 'Oppdatert romnavn' WHERE id = 42`
		case "unpublish": query = `UPDATE program_publishing_state SET is_published = 0`; bucket = live.BucketEvents
		case "publish": query = `UPDATE program_publishing_state SET is_published = 1`; bucket = live.BucketEvents
		case "remove": query = `UPDATE relation_event_puljer SET room_id = NULL WHERE pulje_id = 'LordagMorgen'`
		default: http.NotFound(w,r); return
		}
		if _,err := db.Exec(query); err != nil { http.Error(w,err.Error(),500);return }
		if err := manager.Broadcast(r.Context(),bucket); err != nil { http.Error(w,err.Error(),500);return }
		w.WriteHeader(http.StatusNoContent)
	})
	fmt.Println("Room grouping fixture: http://127.0.0.1:8098/event/room-event")
	t.Fatal(http.ListenAndServe("127.0.0.1:8098",router))
}
