package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var tplApi = `package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
	"go.opencensus.io/trace"

	"gitlab.services.mts.ru/briefcase/brief-api/dto"
	httpEntry "gitlab.services.mts.ru/briefcase/brief-api/internal/utils/log/http/log"
	"gitlab.services.mts.ru/briefcase/brief-api/pkg/bmodels/db/orm"
	"gitlab.services.mts.ru/briefcase/libs/respond"
)

// {{.EntityName}} list
// {{.EntityName}} swagger:route POST /{{.EndpointName}} {{.EntityName}} Create{{.EntityName}}Request
//
// Responses:
//   201: Create{{.EntityName}}Response
//   400: badRequestResponse
//   500: errorResponse
func (ac *apiController) create{{.EntityName}}(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), timeOut)
	defer cancel()

	ctx, span := trace.StartSpan(ctx, "create{{.EntityName}}")
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.AddAttributes(trace.BoolAttribute("error", true))
			span.Annotate(nil, err.Error())
		}
	}()

	create := &dto.Create{{.EntityName}}{}
	if err = json.NewDecoder(r.Body).Decode(create); err != nil {
		respond.WithErrWrapJSON(ctx, w, http.StatusBadRequest, err)
		return
	}

	if err = ac.validator.Struct(create); err != nil {
		respond.WithErrWrapJSON(ctx, w, http.StatusBadRequest, err)
		return
	}

	if err = orm.Create(ctx, create); err != nil {
		respond.WithErrWrapJSON(ctx, w, http.StatusInternalServerError, err)
		return
	}

	respond.WithJSON(ctx, w, http.StatusCreated, create)
}

// {{.EntityPluralName}} list
// {{.EntityPluralName}} swagger:route GET /{{.EndpointName}}/{id} {{.EntityName}} Get{{.EntityName}}Request
//
// Responses:
//   200: Get{{.EntityName}}Response
//   400: badRequestResponse
//   500: errorResponse
func (ac *apiController) get{{.EntityName}}(w http.ResponseWriter, r *http.Request) {
	const funcName = "get{{.EntityName}}"
	ctx, cancel := context.WithTimeout(r.Context(), timeOut)
	defer cancel()

	ctx, span := trace.StartSpan(ctx, funcName)
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.AddAttributes(trace.BoolAttribute("error", true))
			span.Annotate(nil, err.Error())
		}
	}()

	filter := &dto.Get{{.EntityName}}Request{}
	filter.ID = chi.URLParam(r, "id")

	if err = ac.validator.Struct(filter); err != nil {
		respond.WithErrWrapJSON(ctx, w, http.StatusBadRequest, err)
		return
	}

	logParams := map[string]interface{}{"id": filter.ID}
	httpEntry.InitEntry(ctx).LogReqNonEmptyParam(funcName, logParams)

	entity, err := orm.Get{{.EntityName}}(ctx, filter)
	if err != nil {
		respond.WithErrWrapJSON(ctx, w, http.StatusInternalServerError, err)
		return
	}

	respond.WithJSON(ctx, w, http.StatusOK, entity)
}

// {{.EntityPluralName}} list
//
// {{.EntityPluralName}} swagger:route GET /{{.EndpointName}} {{.EntityPluralName}} Get{{.EntityPluralName}}Request
//
// Responses:
//   200: Get{{.EntityPluralName}}Response
//   400: badRequestResponse
//   500: errorResponse
func (ac *apiController) get{{.EntityPluralName}}(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), timeOut)
	defer cancel()

	ctx, span := trace.StartSpan(ctx, "get{{.EntityPluralName}}")
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.AddAttributes(trace.BoolAttribute("error", true))
			span.Annotate(nil, err.Error())
		}
	}()

	decoder := schema.NewDecoder()
	filter := &dto.Get{{.EntityPluralName}}Request{}
	if err = decoder.Decode(filter, r.URL.Query()); err != nil {
		respond.WithErrWrapJSON(ctx, w, http.StatusBadRequest, err)
		return
	}

	if err = ac.validator.Struct(filter); err != nil {
		respond.WithErrWrapJSON(ctx, w, http.StatusBadRequest, err)
		return
	}

	entities, total, err := orm.Get{{.EntityPluralName}}(ctx, filter)
	if err != nil {
		respond.WithErrWrapJSON(ctx, w, http.StatusInternalServerError, err)
		return
	}

	respond.WithJSON(ctx, w, http.StatusOK, dto.Get{{.EntityPluralName}}Response { 
		Items: entities,
		Total: total,
	})
}

// {{.EntityPluralName}} list
// {{.EntityPluralName}} swagger:route PUT /{{.EndpointName}} {{.EntityPluralName}} Update{{.EntityName}}Request
//
// Responses:
//   200: statusOkResponse
//   400: badRequestResponse
//   500: errorResponse
func (ac *apiController) update{{.EntityName}}(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), timeOut)
	defer cancel()

	ctx, span := trace.StartSpan(ctx, "update{{.EntityName}}")
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.AddAttributes(trace.BoolAttribute("error", true))
			span.Annotate(nil, err.Error())
		}
	}()

	update := &dto.Update{{.EntityName}}{}
	if err = json.NewDecoder(r.Body).Decode(update); err != nil {
		respond.WithErrWrapJSON(ctx, w, http.StatusBadRequest, err)
		return
	}

	if err = ac.validator.Struct(update); err != nil {
		respond.WithErrWrapJSON(ctx, w, http.StatusBadRequest, err)
		return
	}

	if err = orm.Update(ctx, update); err != nil {
		respond.WithErrWrapJSON(ctx, w, http.StatusInternalServerError, err)
		return
	}

	respond.WithJSON(ctx, w, http.StatusOK, nil)
}

// {{.EntityPluralName}} list
// {{.EntityPluralName}} swagger:route DELETE /{{.EndpointName}} {{.EntityPluralName}} Delete{{.EntityName}}Request
//
// Responses:
//   200: statusOkResponse
//   400: badRequestResponse
//   500: errorResponse
func (ac *apiController) delete{{.EntityName}}(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), timeOut)
	defer cancel()

	ctx, span := trace.StartSpan(ctx, "delete{{.EntityName}}")
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			span.AddAttributes(trace.BoolAttribute("error", true))
			span.Annotate(nil, err.Error())
		}
	}()

	req := &dto.Delete{{.EntityName}}{}
	if err = json.NewDecoder(r.Body).Decode(req); err != nil {
		respond.WithErrWrapJSON(ctx, w, http.StatusBadRequest, err)
		return
	}

	if err = ac.validator.Struct(req); err != nil {
		respond.WithErrWrapJSON(ctx, w, http.StatusBadRequest, err)
		return
	}

	if err = orm.Archive(ctx, req.TableName(), req.ID); err != nil {
		respond.WithErrWrapJSON(ctx, w, http.StatusInternalServerError, err)
		return
	}

	respond.WithJSON(ctx, w, http.StatusOK, nil)
}

func (s *Server) init{{.EntityName}}API() {
	controller := apiController{
		validator: validator.New(),
	}

	s.Router.Route("/{{.EndpointName}}", func(r chi.Router) {
		r.Get("/", controller.get{{.EntityPluralName}})
		r.Get("/{id}", controller.get{{.EntityName}})
		r.Post("/", controller.create{{.EntityName}})
		r.Put("/", controller.update{{.EntityName}})
		r.Delete("/", controller.delete{{.EntityName}})
	})
}
`
var tplDTO = `package dto

type {{.EntityName}}Table struct{}

func ({{.EntityName}}Table) TableName() string {
	return "{{.Schema}}.{{.TableName}}"
}

// {{.EntityName}} base type
type {{.EntityName}} struct {
	{{.EntityName}}Table
	{{range .Base}}{{.}}{{end}}
}

// swagger:parameters Create{{.EntityName}}Request
type Create{{.EntityName}}Request struct {
	// in: body
	Body Create{{.EntityName}}
}

type Create{{.EntityName}} struct {
	{{.EntityName}}Table
	Create
	{{range .Create}}{{.}}{{end}}
}

// swagger:response Create{{.EntityName}}Response
type Create{{.EntityName}}Response *{{.EntityName}}

// swagger:response Get{{.EntityName}}Response
type Get{{.EntityName}}Response *{{.EntityName}}

// swagger:parameters Get{{.EntityName}}Request
type Get{{.EntityName}}Request struct {
	{{.EntityName}}Table
	SearchRequest
	{{range .GetOne}}{{.}}{{end}}}

// swagger:response Get{{.EntityPluralName}}Response
type Get{{.EntityPluralName}}Response struct {
	Items []{{.EntityName}}
	Total int
}

// swagger:parameters Get{{.EntityPluralName}}Request
type Get{{.EntityPluralName}}Request struct {
	{{.EntityName}}Table
	SearchRequest
	{{range .GetAll}}{{.}}{{end}}
}

// swagger:parameters Update{{.EntityName}}Request
type Update{{.EntityName}}Request struct {
	//in: body
	Body Update{{.EntityName}}
}

type Update{{.EntityName}} struct {
	{{.EntityName}}Table
	{{range .Update}}{{.}}{{end}}
}

// swagger:parameters Delete{{.EntityName}}Request
type Delete{{.EntityName}}Request struct {
	//in: body
	Body Delete{{.EntityName}}
}

type Delete{{.EntityName}} struct {
	{{.EntityName}}Table
	{{range .Delete}}{{.}}{{end}}}
`
var tplORM = `package orm

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"gitlab.services.mts.ru/briefcase/brief-api/dto"
	"go.opencensus.io/trace"
)

func Get{{.EntityName}}(ctx context.Context, filter *dto.Get{{.EntityName}}Request) (result dto.{{.EntityName}}, err error) {
	ctx, span := trace.StartSpan(ctx, "orm.Get{{.EntityName}}")
	defer span.End()

	if err = db.Where("is_deleted = false").Where(filter).First(&result).Error; err != nil {
		err = errors.Wrap(err, "orm.Get{{.EntityName}}: "+fmt.Sprintf("filter = %+v", *filter))
		span.AddAttributes(trace.BoolAttribute("error", true))
		span.Annotate(nil, err.Error())
	}

	return
}

func Get{{.EntityPluralName}}(ctx context.Context, filter *dto.Get{{.EntityPluralName}}Request) (result []dto.{{.EntityName}}, total int, err error) {
	ctx, span := trace.StartSpan(ctx, "orm.Get{{.EntityPluralName}}")
	defer span.End()

	q := db.Where("is_deleted = false").Where(filter)

	q = addSearch(q, filter)
	q = addFilter(q, filter)
	q.Find(&result).Count(&total)
	q = addSort(q, filter)
	q = addPager(q, filter)

	if err = q.Find(&result).Error; err != nil {
		err = errors.Wrap(err, "orm.Get{{.EntityPluralName}}: "+fmt.Sprintf("filter = %+v", *filter))
		span.AddAttributes(trace.BoolAttribute("error", true))
		span.Annotate(nil, err.Error())
	}

	return
}
`

type columnInfo struct {
	TableCatalog string `db:"table_catalog"`
 	TableSchema  string `db:"table_schema"`
	TableName    string `db:"table_name"`
	ColumnName   string `db:"column_name"`
	IsNullable   string `db:"is_nullable"`
	DataType     string `db:"data_type"`
	MaxLength    *int   `db:"character_maximum_length"`
}

func main() {
	conn := flag.String("connection", "postgres://localhost:5432/brief?user=guest&password=guest&sslmode=disable", "database connection string")
	tableName := flag.String("table", "", "table name without schema")
	entityName := flag.String("en", "", "entity name")
	entityNamePlural := flag.String("enp", "", "entity name plural")
	flag.Parse()

	db, err := sqlx.Connect("postgres", *conn)
	if err != nil {
		log.Fatal(err)
	}

	colInf := make([]columnInfo, 0)

	err = db.Select(&colInf, fmt.Sprintf(
		"select " +
			"table_catalog," +
			" table_schema," +
			" table_name," +
			" column_name," +
			" is_nullable," +
			" data_type," +
			" character_maximum_length" +
			" from information_schema.columns " +
			"where table_name = '%s'", *tableName,
		))
	if err != nil {
		log.Fatal(err)
	}

	ta := template.Must(template.New("api").Parse(tplApi))

	file, err := os.Create("api_" + strings.Replace(*tableName, "_", "", -1) + ".go")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	ta.Execute(
		file,
		map[string]string {
			"EntityName": *entityName,
			"EndpointName": strings.Replace(*tableName, "_", "-", -1),
			"EntityPluralName": *entityNamePlural,
		},
	)

	Base   := make([]string, 0, len(colInf))
	Create := make([]string, 0, len(colInf))
	Update := make([]string, 0, len(colInf))
	GetOne := make([]string, 0, 1)
	GetAll := make([]string, 0, len(colInf))
	Delete := make([]string, 0, 1)

	excludeFromCreateAndUpdate := map[string]bool {
		"create_date": true,
		"created_at": true,
		"created": true,
		"create": true,
		"is_deleted": true,
	}

	for i := range colInf {
		Base = append(Base, getStructField(colInf[i]))

		if _, ok := excludeFromCreateAndUpdate[colInf[i].ColumnName]; !ok {
			f := getStructFieldCreate(colInf[i])
			if len(f) > 0 {
				Create = append(Create, f)
			}
		}

		if _, ok := excludeFromCreateAndUpdate[colInf[i].ColumnName]; !ok {
			f := getStructFieldUpdate(colInf[i])
			if len(f) > 0 {
				Update = append(Update, f)
			}
		}

		f := getStructFieldGetOne(colInf[i])
		if len(f) > 0 {
			GetOne = append(GetOne, f)
		}

		f = getStructFieldGetAll(colInf[i])
		if len(f) > 0 {
			GetAll = append(GetAll, f)
		}

		if _, ok := excludeFromCreateAndUpdate[colInf[i].ColumnName]; !ok {
			f := getStructFieldDelete(colInf[i])
			if len(f) > 0 {
				Delete = append(Delete, f)
			}
		}
	}

	td := template.Must(template.New("dto").Parse(tplDTO))
	fName := "dto_" + strings.Replace(*tableName, "_", "", -1) + ".go"
	file2, err := os.Create(fName)
	if err != nil {
		log.Fatal(err)
	}

	td.Execute(
		file2,
		map[string]interface{} {
			"EntityName": *entityName,
			"EntityPluralName": *entityNamePlural,
			"Schema": colInf[0].TableSchema,
			"TableName": colInf[0].TableName,
			"Base": Base,
			"Create": Create,
			"Update": Update,
			"GetOne": GetOne,
			"GetAll": GetAll,
			"Delete": Delete,
		},
	)

	file2.Close()
	txtBuf, _ := ioutil.ReadFile(fName)
	strT := string(txtBuf)
	strT = strings.Replace(strT, "&#34;", "\"", -1)
	ioutil.WriteFile(fName, []byte(strT), 0777)
	time.Sleep(2 * time.Second)

	to := template.Must(template.New("orm").Parse(tplORM))
	fName2 := "orm_" + strings.Replace(strings.ToLower(*entityName), "_", "", -1) + ".go"
	file3, err := os.Create(fName2)
	if err != nil {
		log.Fatal(err)
	}

	to.Execute(
		file3,
		map[string]interface{} {
			"EntityName": *entityName,
			"EntityPluralName": *entityNamePlural,
			"Schema": colInf[0].TableSchema,
			"TableName": colInf[0].TableName,
			"Base": Base,
			"Create": Create,
			"Update": Update,
			"GetOne": GetOne,
			"GetAll": GetAll,
			"Delete": Delete,
		},
	)

	file3.Close()


	exec.Command("go fmt " + fName)
	exec.Command("go fmt " + fName2)

	fmt.Println("Success!")
	fmt.Println("Do not forget add InitAPI() in to api.server.go")
}

func getStructField(col columnInfo) (field string) {
	field = fmt.Sprintf(
		"%s %s `json:\"%s\"`\n\t",
		toCamelCase(col.ColumnName),
		getType(col.DataType),
		toCamelCaseJson(col.ColumnName),
	)
	return
}

func getStructFieldUpdate(col columnInfo) (field string) {
	if col.ColumnName == "id" {
		field = fmt.Sprintf(
			"%s %s `json:\"%s\" validate:\"required,uuid\"`\n\t",
			toCamelCase(col.ColumnName),
			getType(col.DataType),
			toCamelCaseJson(col.ColumnName),
		)
		return
	}

	if col.IsNullable == "NO" {
		if col.MaxLength != nil {
			field = fmt.Sprintf(
				"%s %s `json:\"%s\" validate:\"omitempty,max=%d\"`\n\t",
				toCamelCase(col.ColumnName),
				getType(col.DataType),
				toCamelCaseJson(col.ColumnName),
				col.MaxLength,
			)
			return
		} else {
			field = fmt.Sprintf(
				"%s %s `json:\"%s\" validate:\"omitempty\"`\n\t",
				toCamelCase(col.ColumnName),
				getType(col.DataType),
				toCamelCaseJson(col.ColumnName),
			)
			return
		}
	}

	return fmt.Sprintf(
		"%s *%s `json:\"%s\"`\n",
		toCamelCase(col.ColumnName),
		getType(col.DataType),
		toCamelCaseJson(col.ColumnName),
	)
}

func getStructFieldCreate(col columnInfo) (field string) {
	if col.ColumnName == "id" {
		field = fmt.Sprintf(
			"%s %s `json:\"%s\"`\n\t",
			toCamelCase(col.ColumnName),
			getType(col.DataType),
			toCamelCaseJson(col.ColumnName),
		)
		return
	}

	if col.IsNullable == "NO" {
		if col.MaxLength != nil {
			field = fmt.Sprintf(
				"%s %s `json:\"%s\" validate:\"omitempty,max=%d\"`\n\t",
				toCamelCase(col.ColumnName),
				getType(col.DataType),
				toCamelCaseJson(col.ColumnName),
				col.MaxLength,
			)
			return
		} else {
			field = fmt.Sprintf(
				"%s %s `json:\"%s\" validate:\"omitempty\"`\n\t",
				toCamelCase(col.ColumnName),
				getType(col.DataType),
				toCamelCaseJson(col.ColumnName),
			)
			return
		}
	}

	return fmt.Sprintf(
		"%s *%s `json:\"%s\"`\n\t",
		toCamelCase(col.ColumnName),
		getType(col.DataType),
		toCamelCaseJson(col.ColumnName),
	)
}

func getStructFieldGetOne(col columnInfo) (field string) {
	if col.ColumnName == "id" {
		field = fmt.Sprintf(
			"// in: path\n\t%s %s `schema:\"%s\" validate:\"required,uuid\"`\n",
			toCamelCase(col.ColumnName),
			getType(col.DataType),
			toCamelCaseJson(col.ColumnName),
		)
		return
	}
	return
}

func getStructFieldGetAll(col columnInfo) (field string) {
	return fmt.Sprintf(
		"// in: path\n\t%s %s `schema:\"%s\" validate:\"omitempty\"`\n\t",
		toCamelCase(col.ColumnName),
		getType(col.DataType),
		toCamelCaseJson(col.ColumnName),
	)
}

func getStructFieldDelete(col columnInfo) (field string) {
	if col.ColumnName == "id" {
		field = fmt.Sprintf(
			"%s %s `json:\"%s\" validate:\"required,uuid\"`\n",
			toCamelCase(col.ColumnName),
			getType(col.DataType),
			toCamelCaseJson(col.ColumnName),
		)
		return
	}
	return
}

func getType(dbType string) string {
	t := map[string]string {
		"uuid": "string",
		"character varying": "string",
		"text": "string",
		"jsonb": "string",
		"boolean": "bool",
		"timestamp without time zone": "string",
		"integer": "int",
	}

	if val, ok := t[dbType]; ok {
		return val
	}

	return "string"
}

func toCamelCase(str string) (result string) {
	if str == "id" {
		return "ID"
	}

	if strings.HasSuffix(str, "_id") {
		str = str[:len(str)-3] + "ID"
	}

	skip := false
	for k, v := range str {
		if !skip {
			if k == 0 {
				result = strings.ToUpper(string(v))
				continue
			}

			if string(str[k]) == "_" {
				k++
				result += strings.ToUpper(string(str[k]))
				skip = true
				continue
			} else {
				result += string(v)
			}
		} else {
			skip = false
		}
	}
	return
}

func toCamelCaseJson(str string) (result string) {
	skip := false
	for k, v := range str {
		if !skip {
			if string(str[k]) == "_" {
				k++
				result += strings.ToUpper(string(str[k]))
				skip = true
				continue
			} else {
				result += string(v)
			}
		} else {
			skip = false
		}
	}
	return
}
