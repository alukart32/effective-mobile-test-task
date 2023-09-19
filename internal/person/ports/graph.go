package ports

import (
	"fmt"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/alukart32/effective-mobile-test-task/internal/person/ports/graph"
	gengraph "github.com/alukart32/effective-mobile-test-task/internal/person/ports/graph/generated"
	"github.com/gin-gonic/gin"
)

func Graph(router *gin.Engine, api string, personManager personManager) error {
	if personManager == nil {
		return fmt.Errorf("init GraphQL schema: personManager is nil")
	}

	srv := handler.New(gengraph.NewExecutableSchema(
		gengraph.Config{
			Resolvers: &graph.Resolver{
				PersonManager: personManager,
			},
		},
	))

	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	router.POST(api, func(c *gin.Context) {
		srv.ServeHTTP(c.Writer, c.Request)
	})
	return nil
}
