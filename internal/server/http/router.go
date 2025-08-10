// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"

	"github.com/rasorp/attila/internal/server/nomad"
	"github.com/rasorp/attila/internal/server/state"
)

func NewRouter(logger zerolog.Logger, accessLevel string, stateStore state.State, nomadController nomad.Controller) *chi.Mux {

	r := chi.NewRouter()
	r.Use(loggerMiddleware(logger, accessLevel))

	r.Route("/v1alpha1", func(r chi.Router) {

		r.Mount(
			"/topologies",
			topologiesEndpoint{
				nomadController: nomadController,
			}.routes(),
		)

		r.Mount("/jobs", jobRouter(logger, stateStore, nomadController))

		r.Mount("/regions", regionsEndpoint{
			nomadController: nomadController,
			state:           stateStore,
		}.routes())
	})

	return r
}

func jobRouter(logger zerolog.Logger, stateStore state.State, nomadController nomad.Controller) http.Handler {
	r := chi.NewRouter()

	r.Mount("/register/methods", jobsRegisterMethodsEndpoint{
		state: stateStore,
	}.routes())

	r.Mount("/register/plans", jobsRegisterPlansEndpoint{
		logger:          logger,
		nomadController: nomadController,
		state:           stateStore,
	}.routes())

	r.Mount("/register/rules", jobsRegisterRulesEndpoint{
		state: stateStore,
	}.routes())

	return r
}
