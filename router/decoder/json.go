package decoder

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"reflect"
	"strconv"
	"strings"
)

func NewJsonRequestDecoder(requestObject interface{}) RequestDecoder {
	if requestObject == nil {
		return nil
	}

	return &JsonRequestDecoder{
		requestType: reflect.TypeOf(requestObject),
	}
}

type JsonRequestDecoder struct {
	requestType reflect.Type
}

func (d *JsonRequestDecoder) DecodeRequest(ctx *gin.Context) (interface{}, error) {
	var requestValue reflect.Value
	if d.requestType.Kind() == reflect.Ptr {
		requestValue = reflect.New(d.requestType.Elem())
	} else {
		requestValue = reflect.New(d.requestType)
	}

	// get request parameters from request body
	if err := ctx.ShouldBind(requestValue.Interface()); err != nil {
		return nil, err
	}

	// get request parameters from path
	if err := d.fillPathParameter(ctx, requestValue); err != nil {
		return nil, fmt.Errorf("request decoder: unmarshal request path failed,  error (%v)", err)
	}

	// get request parameters from query parameters
	if err := d.fillQueryParameter(ctx, requestValue); err != nil {
		return nil, fmt.Errorf("request decoder: unmarshal request query failed,  error (%v)", err)
	}

	return requestValue.Elem().Interface(), nil
}

func (d *JsonRequestDecoder) fillPathParameter(ctx *gin.Context, value reflect.Value) error {
	const pathTagKey = "path"
	var objectValue reflect.Value
	if value.Kind() == reflect.Ptr {
		objectValue = value.Elem()
	} else {
		objectValue = value
	}

	objectType := objectValue.Type()
	for i := 0; i < objectValue.NumField(); i++ {
		f := objectType.Field(i)
		tagValue, ok := f.Tag.Lookup(pathTagKey)
		if !ok {
			continue
		}

		fieldValue := ctx.Param(tagValue)
		if err := d.setFieldValue(objectValue.Field(i), fieldValue); err != nil {
			return err
		}
	}

	return nil
}

func (d *JsonRequestDecoder) fillQueryParameter(ctx *gin.Context, value reflect.Value) error {
	const queryTagKey = "query"
	var objectValue reflect.Value
	if value.Kind() == reflect.Ptr {
		objectValue = value.Elem()
	} else {
		objectValue = value
	}

	objectType := objectValue.Type()
	for i := 0; i < objectValue.NumField(); i++ {
		typeField := objectType.Field(i)
		tagValue, ok := typeField.Tag.Lookup(queryTagKey)
		if !ok {
			continue
		}

		valueField := objectValue.Field(i)
		valueFieldKind := valueField.Kind()
		if valueFieldKind == reflect.Array || valueFieldKind == reflect.Slice {
			fieldValues := ctx.QueryArray(tagValue)
			for _, fv := range fieldValues {
				elem := reflect.New(valueField.Type().Elem()).Elem()
				if err := d.setFieldValue(elem, fv); err != nil {
					break
				}
				valueField.Set(reflect.Append(valueField, elem))
			}
		} else {
			fieldValue := ctx.Query(tagValue)
			if err := d.setFieldValue(objectValue.Field(i), fieldValue); err != nil {
				return err
			}
		}
	}

	return nil
}

func (d *JsonRequestDecoder) setFieldValue(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
		break
	case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64:
		iv, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("failed to convert %v to int", value)
		}
		field.SetInt(int64(iv))
		break
	case reflect.Float32, reflect.Float64:
		fv, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("failed to convert %v to float", value)
		}
		field.SetFloat(fv)
		break
	case reflect.Bool:
		if strings.ToLower(value) == "true" {
			field.SetBool(true)
		} else if strings.ToLower(value) == "false" {
			field.SetBool(false)
		} else {
			return fmt.Errorf("boolean field should be one of [true, false], not %v", value)
		}
		break
	default:
		return fmt.Errorf("unsupported field type: %v", field.Kind())
	}
	return nil
}
