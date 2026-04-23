package form

import (
	"reflect"
	"strings"

	"github.com/nef734/ggit/internal/context"
	"github.com/nef734/ggit/internal/lazyregexp"
	"github.com/nef734/ggit/internal/strx"
)

const ErrAlphaDashDotSlash = "AlphaDashDotSlashError"

var AlphaDashDotSlashPattern = lazyregexp.New("[^\\d\\w-_\\./]")

// Validation Error Classifications
const (
	ERR_REQUIRED       = "Required"
	ERR_ALPHA_DASH     = "AlphaDash"
	ERR_ALPHA_DASH_DOT = "AlphaDashDot"
	ERR_SIZE           = "Size"
	ERR_MIN_SIZE       = "MinSize"
	ERR_MAX_SIZE       = "MaxSize"
	ERR_EMAIL          = "Email"
	ERR_URL            = "Url"
	ERR_INCLUDE        = "Include"
)

type Error struct {
	Classification string
	FieldNames     []string
}

type Errors []Error

func (e Errors) Len() int {
	return len(e)
}

func (e *Errors) Add(fieldNames []string, classification string, message string) {
	*e = append(*e, Error{
		Classification: classification,
		FieldNames:     fieldNames,
	})
}

type Form interface {
	Validate(ctx ValidationContext, errs Errors) Errors
}

type ValidationContext interface {
	GetData() map[string]any
	GetLocale() context.Locale
}

// Assign assign form values back to the template data.
func Assign(form any, data map[string]any) {
	typ := reflect.TypeOf(form)
	val := reflect.ValueOf(form)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		fieldName := field.Tag.Get("form")
		// Allow ignored fields in the struct
		if fieldName == "-" {
			continue
		} else if fieldName == "" {
			fieldName = strx.ToSnakeCase(field.Name)
		}

		data[fieldName] = val.Field(i).Interface()
	}
}

func getRuleBody(field reflect.StructField, prefix string) string {
	for rule := range strings.SplitSeq(field.Tag.Get("valid"), ";") {
		if strings.HasPrefix(rule, prefix) {
			return rule[len(prefix) : len(rule)-1]
		}
	}
	return ""
}

func getSize(field reflect.StructField) string {
	return getRuleBody(field, "Size(")
}

func getMinSize(field reflect.StructField) string {
	return getRuleBody(field, "MinSize(")
}

func getMaxSize(field reflect.StructField) string {
	return getRuleBody(field, "MaxSize(")
}

func getInclude(field reflect.StructField) string {
	return getRuleBody(field, "Include(")
}

func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return len(strings.TrimSpace(v.String())) == 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	}
	return v.IsZero()
}

func Check(f any) Errors {
	errs := Errors{}
	typ := reflect.TypeOf(f)
	val := reflect.ValueOf(f)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("valid")
		if tag == "" {
			continue
		}

		v := val.Field(i)
		rules := strings.Split(tag, ";")
		for _, rule := range rules {
			if rule == "Required" {
				if isZero(v) {
					errs = append(errs, Error{Classification: ERR_REQUIRED, FieldNames: []string{field.Name}})
					break
				}
			}
			// NOTE: Add more rule implementations as needed
		}
	}
	return errs
}

func validate(errs Errors, data map[string]any, f Form, l context.Locale) Errors {
	if errs.Len() == 0 {
		errs = Check(f)
	}

	if errs.Len() == 0 {
		return errs
	}

	data["HasError"] = true
	Assign(f, data)

	typ := reflect.TypeOf(f)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		fieldName := field.Tag.Get("form")
		// Allow ignored fields in the struct
		if fieldName == "-" {
			continue
		}

		if errs[0].FieldNames[0] == field.Name {
			data["Err_"+field.Name] = true

			trName := field.Tag.Get("locale")
			if trName == "" {
				trName = l.Tr("form." + field.Name)
			} else {
				trName = l.Tr(trName)
			}

			switch errs[0].Classification {
			case ERR_REQUIRED:
				data["ErrorMsg"] = trName + l.Tr("form.require_error")
			case ERR_ALPHA_DASH:
				data["ErrorMsg"] = trName + l.Tr("form.alpha_dash_error")
			case ERR_ALPHA_DASH_DOT:
				data["ErrorMsg"] = trName + l.Tr("form.alpha_dash_dot_error")
			case ErrAlphaDashDotSlash:
				data["ErrorMsg"] = trName + l.Tr("form.alpha_dash_dot_slash_error")
			case ERR_SIZE:
				data["ErrorMsg"] = trName + l.Tr("form.size_error", getSize(field))
			case ERR_MIN_SIZE:
				data["ErrorMsg"] = trName + l.Tr("form.min_size_error", getMinSize(field))
			case ERR_MAX_SIZE:
				data["ErrorMsg"] = trName + l.Tr("form.max_size_error", getMaxSize(field))
			case ERR_EMAIL:
				data["ErrorMsg"] = trName + l.Tr("form.email_error")
			case ERR_URL:
				data["ErrorMsg"] = trName + l.Tr("form.url_error")
			case ERR_INCLUDE:
				data["ErrorMsg"] = trName + l.Tr("form.include_error", getInclude(field))
			default:
				data["ErrorMsg"] = l.Tr("form.unknown_error") + " " + errs[0].Classification
			}
			return errs
		}
	}
	return errs
}
