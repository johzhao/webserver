package server

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
	"strconv"
)

type Handler func(ctx context.Context, req interface{}) (interface{}, error)

type Decoder func(ctx *gin.Context) (interface{}, error)

type Encoder func(ctx *gin.Context, resp interface{}) error

type RouteConfig struct {
	Method          string
	Path            string
	RequestDecoder  Decoder
	RequestObject   interface{}
	ResponseEncoder Encoder
	Handler         Handler
}

func (r RouteConfig) handle(ctx *gin.Context) {
	req, err := r.RequestDecoder(ctx)
	if err != nil {
		r.handleError(ctx, err)
		return
	}

	resp, err := r.Handler(ctx, req)
	if err != nil {
		r.handleError(ctx, err)
		return
	}

	if err := r.ResponseEncoder(ctx, resp); err != nil {
		r.handleError(ctx, err)
		return
	}
}

func (r RouteConfig) handleError(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusOK, gin.H{
		"code":    -1,
		"message": err.Error(),
	})
}

func setFieldValue(field reflect.Value, value string) error {
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
		if value == "true" {
			field.SetBool(true)
		} else if value == "false" {
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

func fillPathParameter(ctx *gin.Context, toFill reflect.Value) error {
	const pathTagKey = "path"
	var obj reflect.Value
	if toFill.Kind() == reflect.Ptr {
		obj = toFill.Elem()
	} else {
		obj = toFill
	}

	objType := obj.Type()
	for i := 0; i < obj.NumField(); i++ {
		f := objType.Field(i)
		if tagValue, ok := f.Tag.Lookup(pathTagKey); ok {
			fieldValue := ctx.Param(tagValue)
			if err := setFieldValue(obj.Field(i), fieldValue); err != nil {
				return err
			}
		}
	}

	return nil
}

func fillQueryParameter(ctx *gin.Context, toFill reflect.Value) error {
	const queryTagKey = "query"
	var obj reflect.Value
	if toFill.Kind() == reflect.Ptr {
		obj = toFill.Elem()
	} else {
		obj = toFill
	}

	objType := obj.Type()
	for i := 0; i < obj.NumField(); i++ {
		f := objType.Field(i)
		if tagValue, ok := f.Tag.Lookup(queryTagKey); ok {
			field := obj.Field(i)
			fieldKind := obj.Field(i).Kind()

			if fieldKind == reflect.Array || fieldKind == reflect.Slice {
				fieldValues := ctx.QueryArray(tagValue)
				for _, fv := range fieldValues {
					elem := reflect.New(field.Type().Elem()).Elem()
					if err := setFieldValue(elem, fv); err != nil {
						break
					}
					field.Set(reflect.Append(field, elem))
				}
			} else {
				fieldValue := ctx.Query(tagValue)
				if err := setFieldValue(obj.Field(i), fieldValue); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func makeDefaultJsonRequestDecoder(conf *RouteConfig) Decoder {
	return func(ctx *gin.Context) (interface{}, error) {
		var err error
		if conf.RequestObject != nil {
			requestType := reflect.TypeOf(conf.RequestObject)
			var requestValue reflect.Value
			if requestType.Kind() == reflect.Ptr {
				requestValue = reflect.New(requestType.Elem())
			} else {
				requestValue = reflect.New(requestType)
			}

			// get request parameters from request body
			if err := ctx.ShouldBind(requestValue.Interface()); err != nil {
				return nil, err
			}

			// get request parameters and queries from request query
			if err == nil {
				err = fillPathParameter(ctx, requestValue)
				if err != nil {
					err = fmt.Errorf("request decoder: unmarshal request path failed,  error (%v)", err)
				}

				//queries := r.URL.Query()
				err = fillQueryParameter(ctx, requestValue)
				if err != nil {
					err = fmt.Errorf("request decoder: unmarshal request query failed,  error (%v)", err)
				}
			}

			return requestValue.Elem().Interface(), nil
		}

		if err != nil {
			return nil, err
		}
		return nil, nil
	}
}

func defaultJsonResponseEncoder(ctx *gin.Context, resp interface{}) error {
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    resp,
	})
	return nil
}

func MakeRouteHandler(conf *RouteConfig) gin.HandlerFunc {
	if conf.ResponseEncoder == nil {
		conf.ResponseEncoder = defaultJsonResponseEncoder
	}

	if conf.RequestDecoder == nil {
		conf.RequestDecoder = makeDefaultJsonRequestDecoder(conf)
	}

	if len(conf.Method) == 0 {
		conf.Method = "POST"
	}

	return conf.handle
}
