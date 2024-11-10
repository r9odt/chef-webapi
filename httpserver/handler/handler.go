package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/r9odt/chef-webapi/httpserver/errors"
	"github.com/r9odt/chef-webapi/httpserver/middleware"

	"github.com/go-chi/chi/v5"
	middle "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

// APIResponse base object for api responses.
type APIResponse interface{}

var (
	// ContentDir is a path to web interfdace.s files.
	ContentDir = filepath.Join("content", "webjs", "build")
)

// NewHandler set handler for requests and return router.
func NewHandler() http.Handler {
	router := chi.NewRouter()

	// Router middlewares and handlers.
	router.Use(middle.RequestID)
	router.Use(middle.RealIP)
	router.Use(middle.Logger)
	router.Use(middle.Recoverer)
	router.NotFound(notFoundHandler)
	router.MethodNotAllowed(methodNotAllowedHandler)

	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-CSRF-Token",
			middleware.AuthHeader,
			middleware.SessionHeader},
		ExposedHeaders: []string{"X-Total-Count"},
		// AllowCredentials: false,
		MaxAge: 300, // Maximum value not ignored by any of major browsers
	}))

	// Route table.
	router.Route("/", func(router chi.Router) {
		router.Route("/api", func(router chi.Router) {
			router.Route("/authentication", func(router chi.Router) {
				router.Use(middleware.IsAuth)
				router.Post("/auth", authenticationAPIHandler)
				router.Delete("/logout", authenticationLogoutAPIHandler)
				router.Get("/ping", authenticationPingAPIHandler)
				router.Get("/info", authenticationGetCurrentUserAPIHandler)
				router.Get("/perms", authenticationGetCurrentUserPermissionsAPIHandler)
			})
			router.Route("/profiles", func(router chi.Router) {
				router.Use(middleware.IsAuth)
				router.Get("/", profileGetUserProfilesAPIHandler)
				router.Route("/edit", func(router chi.Router) {
					router.Get("/", profileGetCurrentUserAPIHandler)
					router.Put("/", profileUpdateCurrentUserAPIHandler)
				})
				router.Route("/{id}", func(router chi.Router) {
					router.Use(middleware.IDContext)
					router.Get("/", profileGetUserProfileByIDAPIHandler)
				})
			})
			router.Route("/modules", func(router chi.Router) {
				router.Use(middleware.IsAuth)
				router.Use(middleware.HasAccessToModules)
				router.Get("/", appGetModulesAPIHandler)
				router.Route("/{id}", func(router chi.Router) {
					router.Use(middleware.IDContext)
					router.Get("/", appGetModuleByIDAPIHandler)
					router.Put("/", appUpdateModuleByIDAPIHandler)
				})
			})
			// REST
			router.Route("/users", func(router chi.Router) {
				router.Use(middleware.IsAuth)
				router.Use(middleware.HasAccessToUsers)
				router.Get("/", usersGetUsersAPIHandler)
				router.Post("/", usersCreateUserAPIHandler)
				router.Route("/{id}", func(router chi.Router) {
					router.Use(middleware.IDContext)
					router.Get("/", usersGetUserByIDAPIHandler)
					router.Put("/", usersUpdateUserByIDAPIHandler)
					router.Delete("/", usersDeleteUserByIDAPIHandler)
				})
			})
			router.Route("/nodes", func(router chi.Router) {
				router.Use(middleware.IsAuth)
				router.Get("/", nodesAPIHandler)
				router.Route("/{id}", func(router chi.Router) {
					router.Use(middleware.IDContext)
					router.Get("/", nodesGetNodesByNodeNameAPIHandler)
					router.Put("/", nodesCreateTaskForResourceAPIHandler)
				})
			})
			router.Route("/roles", func(router chi.Router) {
				router.Use(middleware.IsAuth)
				router.Get("/", rolesAPIHandler)
				router.Route("/{id}", func(router chi.Router) {
					router.Use(middleware.IDContext)
					router.Get("/", rolesGetNodesByRoleNameAPIHandler)
					router.Put("/", rolesCreateTaskForResourceAPIHandler)
				})
			})
			// REST
			router.Route("/deployers", func(router chi.Router) {
				router.Use(middleware.IsAuth)
				router.Get("/", deployersAPIHandler)
				router.Route("/{id}", func(router chi.Router) {
					router.Use(middleware.IDContext)
					router.Put("/", deployersCreateDeployerAPIHandler)
					router.Get("/", deployersGetDeployerByIDAPIHandler)
				})
			})
			router.With(middleware.IsAuth).Get("/cookbooks", cookbooksAPIHandler)
			router.Route("/search", func(router chi.Router) {
				router.Use(middleware.IsAuth)
				router.Route("/{searchIndex}", func(router chi.Router) {
					router.Use(middleware.SearchIndexContext)
					router.Route("/{searchQuery}", func(router chi.Router) {
						router.Use(middleware.SearchQueryContext)
						router.Get("/", searchAPIHandler)
					})
				})
			})
			router.Route("/keys", func(router chi.Router) {
				router.Use(middleware.IsAuth)
				router.Use(middleware.HasAccessToKeys)
				router.Get("/", keysGetKeysAPIHandler)
				router.Route("/{id}", func(router chi.Router) {
					router.Use(middleware.IDContext)
					router.Get("/", keysGetKeyByIDAPIHandler)
					router.Put("/", keysUpdateKeyByIDAPIHandler)
				})
			})
			router.Route("/monitor", func(router chi.Router) {
				router.Get("/ready", monitorReadyAPIHandler)
				router.Get("/health", monitorHealthAPIHandler)
			})
		})
		workDir, _ := os.Getwd()
		dir := http.Dir(
			filepath.Join(workDir, ContentDir))
		content(router, "/", dir)
	})

	return router
}

// content conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func content(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/",
			http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}

func notFoundHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("X-Content-Type-Options", "nosniff")
	_ = render.Render(writer, request,
		errors.ErrorNotFound(fmt.Errorf("not found")))
}

func methodNotAllowedHandler(
	writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(405)
}
