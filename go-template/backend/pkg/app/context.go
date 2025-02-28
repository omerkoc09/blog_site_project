package app

import (
	"mime/multipart"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-playground/locales/tr"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	turkish "github.com/go-playground/validator/v10/translations/tr"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"go.uber.org/zap/zapcore"

	"github.com/hayrat/go-template2/backend/pkg/errorsx"
	"github.com/hayrat/go-template2/backend/pkg/log"
	"github.com/hayrat/go-template2/backend/pkg/viewmodel"
)

type Ctx struct {
	*fiber.Ctx
}

func (c *Ctx) SuccessResponse(data interface{}, dataCount ...int) error {
	m := &viewmodel.ResponseModel{
		Data: data,
	}
	if len(dataCount) > 0 {
		m.DataCount = dataCount[0]
		if m.Data == nil {
			m.Data = []any{}
		}
	}
	err := c.JSON(m)
	if err != nil {
		return err
	}
	return nil
}

func (c *Ctx) SetLogFields(fields ...zapcore.Field) {
	l := c.Locals("logger").(*log.Logger)
	l.SetFields(fields...)
}

func (c *Ctx) ErrorLog(msg string, fields ...zapcore.Field) {
	l := c.Locals("logger").(*log.Logger)
	l.Error(msg, fields...)
}

func (c *Ctx) InfoLog(msg string, fields ...zapcore.Field) {
	l := c.Locals("logger").(*log.Logger)
	l.Info(msg, fields...)
}

func (c *Ctx) GetPaginationModel() (*viewmodel.PaginationModel, error) {
	m := new(viewmodel.PaginationModel)
	if err := c.QueryParser(m); err != nil {
		return nil, err
	}
	if m.Page == 0 {
		m.Page = 1
	}
	if m.PerPage == 0 {
		m.PerPage = 500
	}
	if len(m.SortColumns) == 0 {
		m.SortColumns = append(m.SortColumns, "created_at")
	}
	if len(m.SortColumnTypes) == 0 {
		m.SortColumnTypes = append(m.SortColumnTypes, viewmodel.SortColumnTypeNormal)
	}
	if len(m.SortOrders) == 0 {
		m.SortOrders = append(m.SortOrders, viewmodel.SortOrderASC)
	}
	m.Offset = (m.Page - 1) * m.PerPage

	return m, nil
}

func (c *Ctx) GetQueryModel() (*viewmodel.QueryModel, error) {
	m := new(viewmodel.QueryModel)
	if err := c.QueryParser(m); err != nil {
		return nil, err
	}
	if len(m.Query) == 0 {
		m.Query = make([]string, 0)
		m.Columns = make([]string, 0)
		m.ColumnTypes = make([]viewmodel.QueryColumnType, 0)
	}
	if strTarih := c.Query("updated_at", ""); strTarih != "" {
		m.Append(strTarih, "updated_at", viewmodel.QueryColumnTypeDateBuyuk)
	}
	return m, nil
}

func (c *Ctx) GetQueryPaginationModel() (*viewmodel.QueryModel, *viewmodel.PaginationModel, error) {
	q, err := c.GetQueryModel()
	if err != nil {
		return nil, nil, err
	}

	p, err := c.GetPaginationModel()
	if err != nil {
		return nil, nil, err
	}

	return q, p, nil
}

func (c *Ctx) BodyParseValidate(m interface{}) []error {

	if v := reflect.ValueOf(m); v.Kind() != reflect.Ptr {
		if err := c.BodyParser(&m); err != nil {
			return []error{err}
		}
	} else {
		if err := c.BodyParser(m); err != nil {
			return []error{err}
		}
	}

	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("label"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	var err error
	var trans ut.Translator

	turkce := tr.New()
	uni := ut.New(turkce)
	trans, _ = uni.GetTranslator("tr")
	err = turkish.RegisterDefaultTranslations(validate, trans)

	if err != nil {
		return []error{err}
	}

	var errs []error
	err = validate.Struct(m)
	if err != nil {
		vErrs := err.(validator.ValidationErrors)
		for _, e := range vErrs {
			errs = append(errs, errorsx.New(e.Translate(trans)))
		}
	}

	viewmodelValidation, ok := m.(viewmodel.Validation)
	if ok {
		errs = append(errs, viewmodelValidation.Validate()...)
	}

	return errs
}

func (c *Ctx) ParamsInt64(key string) int64 {
	value, err := strconv.ParseInt(c.Params(key), 10, 64)
	if err != nil {
		return 0
	}

	return value
}

func (c *Ctx) ParamsUUID(key string) uuid.UUID {
	value, err := uuid.Parse(c.Params(key))
	if err != nil {
		return uuid.Nil
	}

	return value
}

func (c *Ctx) RequestID() string {
	return RequestID(c.Ctx)
}

func (c *Ctx) GetUserID() int64 {
	return GetUserID(c.Ctx)
}

func (c *Ctx) GetUserUUID() uuid.UUID {
	return GetUserUUID(c.Ctx)
}

func (c *Ctx) getClaims() jwt.MapClaims {
	u := c.Locals("user")
	if u == nil {
		return nil
	}
	user, ok := u.(*jwt.Token)
	if !ok {
		return nil
	}
	return user.Claims.(jwt.MapClaims)
}

func (c *Ctx) GetFormFile(key string) (multipart.File, string, error) {
	file, err := c.FormFile(key)
	if err != nil {
		return nil, "", err
	}
	f, err := file.Open()
	if err != nil {
		return nil, "", err
	}
	return f, file.Filename, nil
}

func (c *Ctx) FileDownload(b *[]byte, fileName string, contentType ...string) error {
	cType := "application/octet-stream"
	if len(contentType) > 0 {
		cType = contentType[0]
	}
	c.Response().Header.Set("Content-Type", cType)
	c.Response().Header.Set("Content-Length", strconv.Itoa(len(*b)))
	c.Response().Header.Set("Content-Disposition", "attachment; filename="+fileName)
	c.Status(http.StatusOK)
	_, err := c.Write(*b)
	return err
}
