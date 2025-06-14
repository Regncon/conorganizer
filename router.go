package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Regncon/conorganizer/pages/admin"
	billettholderadmin "github.com/Regncon/conorganizer/pages/admin/billettholder_admin"
	"github.com/Regncon/conorganizer/pages/event"
	"github.com/Regncon/conorganizer/pages/index"
	"github.com/Regncon/conorganizer/pages/login"
	"github.com/Regncon/conorganizer/pages/myprofile"
	"github.com/Regncon/conorganizer/pages/myprofile/myevents"
	"github.com/Regncon/conorganizer/service/authctx"
	"github.com/delaneyj/toolbelt"
	"github.com/delaneyj/toolbelt/embeddednats"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	natsserver "github.com/nats-io/nats-server/v2/server"
)

func setupRoutes(ctx context.Context, logger *slog.Logger, routerPublic chi.Router, db *sql.DB) (cleanup func() error, err error) {
	natsPort, err := toolbelt.FreePort()
	if err != nil {
		return nil, fmt.Errorf("error getting free port: %w", err)
	}

	ns, err := embeddednats.New(ctx, embeddednats.WithNATSServerOptions(&natsserver.Options{
		JetStream: true,
		Port:      natsPort,
	}))

	if err != nil {
		return nil, fmt.Errorf("error creating embedded nats server: %w", err)
	}

	ns.WaitForServer()

	cleanup = func() error {
		return errors.Join(
			ns.Close(),
		)
	}

	sessionStore := sessions.NewCookieStore([]byte("session-secret"))
	sessionStore.MaxAge(int(24 * time.Hour / time.Second))

	routerAuth := routerPublic.With(authctx.AuthMiddleware(logger))
	routerAdmin := routerAuth.With(authctx.RequireAdmin(logger))

	if err := errors.Join(
		index.SetupIndexRoute(routerPublic, sessionStore, ns, db),
		admin.SetupAdminRoute(routerAdmin, sessionStore, logger, ns, db),
		billettholderadmin.SetupBillettholderAdminRoute(routerAdmin, sessionStore, ns, logger, db),
		event.SetupEventRoute(routerPublic, sessionStore, ns, db, logger),
		myevents.SetupMyEventsRoute(routerAuth, sessionStore, ns, db, logger),
		login.SetupAuthRoute(routerPublic, db, logger),
		myprofile.SetupMyProfileRoute(routerAuth, sessionStore, ns, db, logger),
	); err != nil {
		return cleanup, fmt.Errorf("error setting up routes: %w", err)
	}

	return cleanup, nil
}
