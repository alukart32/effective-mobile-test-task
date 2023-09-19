// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package generated

import (
	"bytes"
	"context"
	"errors"
	"sync/atomic"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/introspection"
	"github.com/alukart32/effective-mobile-test-task/internal/person/ports/graph/model"
	gqlparser "github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
)

// NewExecutableSchema creates an ExecutableSchema from the ResolverRoot interface.
func NewExecutableSchema(cfg Config) graphql.ExecutableSchema {
	return &executableSchema{
		resolvers:  cfg.Resolvers,
		directives: cfg.Directives,
		complexity: cfg.Complexity,
	}
}

type Config struct {
	Resolvers  ResolverRoot
	Directives DirectiveRoot
	Complexity ComplexityRoot
}

type ResolverRoot interface {
	Mutation() MutationResolver
	Query() QueryResolver
}

type DirectiveRoot struct {
}

type ComplexityRoot struct {
	CreatePersonResponse struct {
		PersonID func(childComplexity int) int
		Success  func(childComplexity int) int
	}

	DeletePersonResponse struct {
		Success func(childComplexity int) int
	}

	Mutation struct {
		CreatePerson func(childComplexity int, input model.CreatePersonInput) int
		DeletePerson func(childComplexity int, input model.DeletePersonInput) int
		UpdatePerson func(childComplexity int, input model.UpdatePersonInput) int
	}

	Person struct {
		Age        func(childComplexity int) int
		Gender     func(childComplexity int) int
		ID         func(childComplexity int) int
		Name       func(childComplexity int) int
		Nation     func(childComplexity int) int
		Patronymic func(childComplexity int) int
		Surname    func(childComplexity int) int
	}

	Query struct {
		CollectPersons func(childComplexity int, limit *int, offset *int, filter *model.CollectPersonsFilter) int
		FindByID       func(childComplexity int, personID string) int
		GetAllPersons  func(childComplexity int) int
	}

	UpdatePersonResponse struct {
		Success func(childComplexity int) int
	}
}

type executableSchema struct {
	resolvers  ResolverRoot
	directives DirectiveRoot
	complexity ComplexityRoot
}

func (e *executableSchema) Schema() *ast.Schema {
	return parsedSchema
}

func (e *executableSchema) Complexity(typeName, field string, childComplexity int, rawArgs map[string]interface{}) (int, bool) {
	ec := executionContext{nil, e, 0, 0, nil}
	_ = ec
	switch typeName + "." + field {

	case "CreatePersonResponse.personId":
		if e.complexity.CreatePersonResponse.PersonID == nil {
			break
		}

		return e.complexity.CreatePersonResponse.PersonID(childComplexity), true

	case "CreatePersonResponse.success":
		if e.complexity.CreatePersonResponse.Success == nil {
			break
		}

		return e.complexity.CreatePersonResponse.Success(childComplexity), true

	case "DeletePersonResponse.success":
		if e.complexity.DeletePersonResponse.Success == nil {
			break
		}

		return e.complexity.DeletePersonResponse.Success(childComplexity), true

	case "Mutation.CreatePerson":
		if e.complexity.Mutation.CreatePerson == nil {
			break
		}

		args, err := ec.field_Mutation_CreatePerson_args(context.TODO(), rawArgs)
		if err != nil {
			return 0, false
		}

		return e.complexity.Mutation.CreatePerson(childComplexity, args["input"].(model.CreatePersonInput)), true

	case "Mutation.DeletePerson":
		if e.complexity.Mutation.DeletePerson == nil {
			break
		}

		args, err := ec.field_Mutation_DeletePerson_args(context.TODO(), rawArgs)
		if err != nil {
			return 0, false
		}

		return e.complexity.Mutation.DeletePerson(childComplexity, args["input"].(model.DeletePersonInput)), true

	case "Mutation.UpdatePerson":
		if e.complexity.Mutation.UpdatePerson == nil {
			break
		}

		args, err := ec.field_Mutation_UpdatePerson_args(context.TODO(), rawArgs)
		if err != nil {
			return 0, false
		}

		return e.complexity.Mutation.UpdatePerson(childComplexity, args["input"].(model.UpdatePersonInput)), true

	case "Person.age":
		if e.complexity.Person.Age == nil {
			break
		}

		return e.complexity.Person.Age(childComplexity), true

	case "Person.gender":
		if e.complexity.Person.Gender == nil {
			break
		}

		return e.complexity.Person.Gender(childComplexity), true

	case "Person.id":
		if e.complexity.Person.ID == nil {
			break
		}

		return e.complexity.Person.ID(childComplexity), true

	case "Person.name":
		if e.complexity.Person.Name == nil {
			break
		}

		return e.complexity.Person.Name(childComplexity), true

	case "Person.nation":
		if e.complexity.Person.Nation == nil {
			break
		}

		return e.complexity.Person.Nation(childComplexity), true

	case "Person.patronymic":
		if e.complexity.Person.Patronymic == nil {
			break
		}

		return e.complexity.Person.Patronymic(childComplexity), true

	case "Person.surname":
		if e.complexity.Person.Surname == nil {
			break
		}

		return e.complexity.Person.Surname(childComplexity), true

	case "Query.CollectPersons":
		if e.complexity.Query.CollectPersons == nil {
			break
		}

		args, err := ec.field_Query_CollectPersons_args(context.TODO(), rawArgs)
		if err != nil {
			return 0, false
		}

		return e.complexity.Query.CollectPersons(childComplexity, args["limit"].(*int), args["offset"].(*int), args["filter"].(*model.CollectPersonsFilter)), true

	case "Query.FindById":
		if e.complexity.Query.FindByID == nil {
			break
		}

		args, err := ec.field_Query_FindById_args(context.TODO(), rawArgs)
		if err != nil {
			return 0, false
		}

		return e.complexity.Query.FindByID(childComplexity, args["PersonId"].(string)), true

	case "Query.GetAllPersons":
		if e.complexity.Query.GetAllPersons == nil {
			break
		}

		return e.complexity.Query.GetAllPersons(childComplexity), true

	case "UpdatePersonResponse.success":
		if e.complexity.UpdatePersonResponse.Success == nil {
			break
		}

		return e.complexity.UpdatePersonResponse.Success(childComplexity), true

	}
	return 0, false
}

func (e *executableSchema) Exec(ctx context.Context) graphql.ResponseHandler {
	rc := graphql.GetOperationContext(ctx)
	ec := executionContext{rc, e, 0, 0, make(chan graphql.DeferredResult)}
	inputUnmarshalMap := graphql.BuildUnmarshalerMap(
		ec.unmarshalInputCollectPersonsFilter,
		ec.unmarshalInputCreatePersonInput,
		ec.unmarshalInputDeletePersonInput,
		ec.unmarshalInputUpdatePersonInput,
	)
	first := true

	switch rc.Operation.Operation {
	case ast.Query:
		return func(ctx context.Context) *graphql.Response {
			var response graphql.Response
			var data graphql.Marshaler
			if first {
				first = false
				ctx = graphql.WithUnmarshalerMap(ctx, inputUnmarshalMap)
				data = ec._Query(ctx, rc.Operation.SelectionSet)
			} else {
				if atomic.LoadInt32(&ec.pendingDeferred) > 0 {
					result := <-ec.deferredResults
					atomic.AddInt32(&ec.pendingDeferred, -1)
					data = result.Result
					response.Path = result.Path
					response.Label = result.Label
					response.Errors = result.Errors
				} else {
					return nil
				}
			}
			var buf bytes.Buffer
			data.MarshalGQL(&buf)
			response.Data = buf.Bytes()
			if atomic.LoadInt32(&ec.deferred) > 0 {
				hasNext := atomic.LoadInt32(&ec.pendingDeferred) > 0
				response.HasNext = &hasNext
			}

			return &response
		}
	case ast.Mutation:
		return func(ctx context.Context) *graphql.Response {
			if !first {
				return nil
			}
			first = false
			ctx = graphql.WithUnmarshalerMap(ctx, inputUnmarshalMap)
			data := ec._Mutation(ctx, rc.Operation.SelectionSet)
			var buf bytes.Buffer
			data.MarshalGQL(&buf)

			return &graphql.Response{
				Data: buf.Bytes(),
			}
		}

	default:
		return graphql.OneShot(graphql.ErrorResponse(ctx, "unsupported GraphQL operation"))
	}
}

type executionContext struct {
	*graphql.OperationContext
	*executableSchema
	deferred        int32
	pendingDeferred int32
	deferredResults chan graphql.DeferredResult
}

func (ec *executionContext) processDeferredGroup(dg graphql.DeferredGroup) {
	atomic.AddInt32(&ec.pendingDeferred, 1)
	go func() {
		ctx := graphql.WithFreshResponseContext(dg.Context)
		dg.FieldSet.Dispatch(ctx)
		ds := graphql.DeferredResult{
			Path:   dg.Path,
			Label:  dg.Label,
			Result: dg.FieldSet,
			Errors: graphql.GetErrors(ctx),
		}
		// null fields should bubble up
		if dg.FieldSet.Invalids > 0 {
			ds.Result = graphql.Null
		}
		ec.deferredResults <- ds
	}()
}

func (ec *executionContext) introspectSchema() (*introspection.Schema, error) {
	if ec.DisableIntrospection {
		return nil, errors.New("introspection disabled")
	}
	return introspection.WrapSchema(parsedSchema), nil
}

func (ec *executionContext) introspectType(name string) (*introspection.Type, error) {
	if ec.DisableIntrospection {
		return nil, errors.New("introspection disabled")
	}
	return introspection.WrapTypeFromDef(parsedSchema, parsedSchema.Types[name]), nil
}

var sources = []*ast.Source{
	{Name: "../../../../../api/graph/schema.graphqls", Input: `type Person {
    id: String!
    name: String!
    surname: String!
    patronymic: String!
    nation: String!
    gender: String!
    age: Int!
}

type Query {
  GetAllPersons: [Person!]!
  CollectPersons(limit: Int, offset: Int, filter: CollectPersonsFilter): [Person!]
  FindById(PersonId: String!): Person
}

input CreatePersonInput {
  name: String!
  surname: String!
  patronymic: String!
}

type CreatePersonResponse {
  success: Boolean!
  personId: String!
}

input CollectPersonsFilter {
  olderThan: Int
  youngerThan: Int
	gender: String
	nations: [String!]
}


input UpdatePersonInput {
  personId: String!
  newNation: String
  newGender: String
  newAge: Int
}

type UpdatePersonResponse {
  success: Boolean!
}

input DeletePersonInput {
  personId: String!
}

type DeletePersonResponse {
  success: Boolean!
}

type Mutation {
  CreatePerson(input: CreatePersonInput!): CreatePersonResponse!
  UpdatePerson(input: UpdatePersonInput!): UpdatePersonResponse!
  DeletePerson(input: DeletePersonInput!): DeletePersonResponse!
}
`, BuiltIn: false},
}
var parsedSchema = gqlparser.MustLoadSchema(sources...)
