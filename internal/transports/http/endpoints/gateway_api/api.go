package gateway_api

import (
	"github.com/gin-gonic/gin"
	GatewayApiProto "gitlab.com/wbwapis/go-genproto/wbw/gateway/gateway_api/v1"
	"gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/errors"
	"gitlab.com/wordbyword.io/microservices/pkg/constants"
	"gitlab.com/wordbyword.io/microservices/pkg/errors/outer"
	"gitlab.com/wordbyword.io/microservices/pkg/utils"
	"google.golang.org/protobuf/encoding/protojson"
	"io"
	"net/http"
	"strconv"
	"time"
)

func (e *GatewayApiHttpEndpoint) SignUp(c *gin.Context) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		e.prom.HttpReqDuration.WithLabelValues("SignUp").Observe(duration.Seconds())
	}()

	requestID := utils.AnyToString(c.Value(constants.RequestIdKey))
	lgr := e.lgr.With().
		Str(constants.RequestIdKey, requestID).
		Str("handler", "SignUp").Logger()

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToGetRequestBody)
		if c.Writer.Status() == http.StatusRequestEntityTooLarge {
			return
		}
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	defer c.Request.Body.Close()

	req := new(GatewayApiProto.SignUpRequest)
	if err = protojson.Unmarshal(body, req); err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToUnmarshalRequestBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	lgr = lgr.With().Interface("request", req).Logger()

	if err = req.Validate(); err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToValidateRequestBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}

	resp, err := e.gatewayApiService.SignUp(c, req)
	if err != nil {
		lgr.Error().Err(err).Msg("failed")
		code, obj := outer.GetHTTPError(err)
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}

	lgr.Debug().Interface("response", resp).Msg("executed")

	buf, err := protojson.MarshalOptions{
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
		UseProtoNames:   true,
	}.Marshal(resp)
	if err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToMarshalResponseBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	c.Data(http.StatusOK, "application/json", buf)
}

func (e *GatewayApiHttpEndpoint) RefreshTokens(c *gin.Context) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		e.prom.HttpReqDuration.WithLabelValues("RefreshTokens").Observe(duration.Seconds())
	}()

	requestID := utils.AnyToString(c.Value(constants.RequestIdKey))
	lgr := e.lgr.With().
		Str(constants.RequestIdKey, requestID).
		Str("handler", "RefreshTokens").Logger()

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToGetRequestBody)
		if c.Writer.Status() == http.StatusRequestEntityTooLarge {
			return
		}
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	defer c.Request.Body.Close()

	req := new(GatewayApiProto.RefreshTokensRequest)
	if err = protojson.Unmarshal(body, req); err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToUnmarshalRequestBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	lgr = lgr.With().Interface("request", req).Logger()

	if err = req.Validate(); err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToValidateRequestBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}

	resp, err := e.gatewayApiService.RefreshTokens(c, req)
	if err != nil {
		lgr.Error().Err(err).Msg("failed")
		code, obj := outer.GetHTTPError(err)
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}

	lgr.Debug().Interface("response", resp).Msg("executed")

	buf, err := protojson.MarshalOptions{
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
		UseProtoNames:   true,
	}.Marshal(resp)
	if err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToMarshalResponseBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	c.Data(http.StatusOK, "application/json", buf)
}

func (e *GatewayApiHttpEndpoint) SignIn(c *gin.Context) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		e.prom.HttpReqDuration.WithLabelValues("SignIn").Observe(duration.Seconds())
	}()

	requestID := utils.AnyToString(c.Value(constants.RequestIdKey))
	lgr := e.lgr.With().
		Str(constants.RequestIdKey, requestID).
		Str("handler", "SignIn").Logger()

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToGetRequestBody)
		if c.Writer.Status() == http.StatusRequestEntityTooLarge {
			return
		}
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	defer c.Request.Body.Close()

	req := new(GatewayApiProto.SignInRequest)
	if err = protojson.Unmarshal(body, req); err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToUnmarshalRequestBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	lgr = lgr.With().Interface("request", req).Logger()

	if err = req.Validate(); err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToValidateRequestBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}

	resp, err := e.gatewayApiService.SignIn(c, req)
	if err != nil {
		lgr.Error().Err(err).Msg("failed")
		code, obj := outer.GetHTTPError(err)
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}

	lgr.Debug().Interface("response", resp).Msg("executed")

	buf, err := protojson.MarshalOptions{
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
		UseProtoNames:   true,
	}.Marshal(resp)
	if err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToMarshalResponseBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	c.Data(http.StatusOK, "application/json", buf)
}

func (e *GatewayApiHttpEndpoint) Logout(c *gin.Context) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		e.prom.HttpReqDuration.WithLabelValues("Logout").Observe(duration.Seconds())
	}()

	requestID := utils.AnyToString(c.Value(constants.RequestIdKey))
	lgr := e.lgr.With().
		Str(constants.RequestIdKey, requestID).
		Str("handler", "Logout").Logger()

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToGetRequestBody)
		if c.Writer.Status() == http.StatusRequestEntityTooLarge {
			return
		}
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	defer c.Request.Body.Close()

	req := new(GatewayApiProto.LogoutRequest)
	if err = protojson.Unmarshal(body, req); err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToUnmarshalRequestBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	lgr = lgr.With().Interface("request", req).Logger()

	if err = req.Validate(); err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToValidateRequestBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}

	resp, err := e.gatewayApiService.Logout(c, req)
	if err != nil {
		lgr.Error().Err(err).Msg("failed")
		code, obj := outer.GetHTTPError(err)
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}

	lgr.Debug().Interface("response", resp).Msg("executed")

	buf, err := protojson.MarshalOptions{
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
		UseProtoNames:   true,
	}.Marshal(resp)
	if err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToMarshalResponseBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	c.Data(http.StatusOK, "application/json", buf)
}

func (e *GatewayApiHttpEndpoint) GetUser(c *gin.Context) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		e.prom.HttpReqDuration.WithLabelValues("GetUser").Observe(duration.Seconds())
	}()

	requestID := utils.AnyToString(c.Value(constants.RequestIdKey))
	lgr := e.lgr.With().
		Str(constants.RequestIdKey, requestID).
		Str("handler", "GetUser").Logger()

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToGetRequestBody)
		if c.Writer.Status() == http.StatusRequestEntityTooLarge {
			return
		}
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	defer c.Request.Body.Close()

	req := new(GatewayApiProto.GetUserRequest)
	if err = protojson.Unmarshal(body, req); err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToUnmarshalRequestBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	lgr = lgr.With().Interface("request", req).Logger()

	if err = req.Validate(); err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToValidateRequestBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}

	resp, err := e.gatewayApiService.GetUser(c, req)
	if err != nil {
		lgr.Error().Err(err).Msg("failed")
		code, obj := outer.GetHTTPError(err)
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}

	lgr.Debug().Interface("response", resp).Msg("executed")

	buf, err := protojson.MarshalOptions{
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
		UseProtoNames:   true,
	}.Marshal(resp)
	if err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToMarshalResponseBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	c.Data(http.StatusOK, "application/json", buf)
}

func (e *GatewayApiHttpEndpoint) UpdateUser(c *gin.Context) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		e.prom.HttpReqDuration.WithLabelValues("UpdateUser").Observe(duration.Seconds())
	}()

	requestID := utils.AnyToString(c.Value(constants.RequestIdKey))
	lgr := e.lgr.With().
		Str(constants.RequestIdKey, requestID).
		Str("handler", "UpdateUser").Logger()

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToGetRequestBody)
		if c.Writer.Status() == http.StatusRequestEntityTooLarge {
			return
		}
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	defer c.Request.Body.Close()

	req := new(GatewayApiProto.UpdateUserRequest)
	if err = protojson.Unmarshal(body, req); err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToUnmarshalRequestBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	lgr = lgr.With().Interface("request", req).Logger()

	if err = req.Validate(); err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToValidateRequestBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}

	resp, err := e.gatewayApiService.UpdateUser(c, req)
	if err != nil {
		lgr.Error().Err(err).Msg("failed")
		code, obj := outer.GetHTTPError(err)
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}

	lgr.Debug().Interface("response", resp).Msg("executed")

	buf, err := protojson.MarshalOptions{
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
		UseProtoNames:   true,
	}.Marshal(resp)
	if err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToMarshalResponseBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	c.Data(http.StatusOK, "application/json", buf)
}

func (e *GatewayApiHttpEndpoint) ConfirmEmail(c *gin.Context) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		e.prom.HttpReqDuration.WithLabelValues("ConfirmEmail").Observe(duration.Seconds())
	}()

	requestID := utils.AnyToString(c.Value(constants.RequestIdKey))
	lgr := e.lgr.With().
		Str(constants.RequestIdKey, requestID).
		Str("handler", "ConfirmEmail").Logger()

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToGetRequestBody)
		if c.Writer.Status() == http.StatusRequestEntityTooLarge {
			return
		}
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	defer c.Request.Body.Close()

	req := new(GatewayApiProto.ConfirmEmailRequest)
	if err = protojson.Unmarshal(body, req); err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToUnmarshalRequestBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	lgr = lgr.With().Interface("request", req).Logger()

	if err = req.Validate(); err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToValidateRequestBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}

	resp, err := e.gatewayApiService.ConfirmEmail(c, req)
	if err != nil {
		lgr.Error().Err(err).Msg("failed")
		code, obj := outer.GetHTTPError(err)
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}

	lgr.Debug().Interface("response", resp).Msg("executed")

	buf, err := protojson.MarshalOptions{
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
		UseProtoNames:   true,
	}.Marshal(resp)
	if err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToMarshalResponseBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	c.Data(http.StatusOK, "application/json", buf)
}

func (e *GatewayApiHttpEndpoint) AskResetPassword(c *gin.Context) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		e.prom.HttpReqDuration.WithLabelValues("AskResetPassword").Observe(duration.Seconds())
	}()

	requestID := utils.AnyToString(c.Value(constants.RequestIdKey))
	lgr := e.lgr.With().
		Str(constants.RequestIdKey, requestID).
		Str("handler", "AskResetPassword").Logger()

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToGetRequestBody)
		if c.Writer.Status() == http.StatusRequestEntityTooLarge {
			return
		}
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	defer c.Request.Body.Close()

	req := new(GatewayApiProto.AskResetPasswordRequest)
	if err = protojson.Unmarshal(body, req); err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToUnmarshalRequestBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	lgr = lgr.With().Interface("request", req).Logger()

	if err = req.Validate(); err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToValidateRequestBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}

	resp, err := e.gatewayApiService.AskResetPassword(c, req)
	if err != nil {
		lgr.Error().Err(err).Msg("failed")
		code, obj := outer.GetHTTPError(err)
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}

	lgr.Debug().Interface("response", resp).Msg("executed")

	buf, err := protojson.MarshalOptions{
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
		UseProtoNames:   true,
	}.Marshal(resp)
	if err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToMarshalResponseBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	c.Data(http.StatusOK, "application/json", buf)
}

func (e *GatewayApiHttpEndpoint) ResetPassword(c *gin.Context) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		e.prom.HttpReqDuration.WithLabelValues("ResetPassword").Observe(duration.Seconds())
	}()

	requestID := utils.AnyToString(c.Value(constants.RequestIdKey))
	lgr := e.lgr.With().
		Str(constants.RequestIdKey, requestID).
		Str("handler", "ResetPassword").Logger()

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToGetRequestBody)
		if c.Writer.Status() == http.StatusRequestEntityTooLarge {
			return
		}
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	defer c.Request.Body.Close()

	req := new(GatewayApiProto.ResetPasswordRequest)
	if err = protojson.Unmarshal(body, req); err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToUnmarshalRequestBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	lgr = lgr.With().Interface("request", req).Logger()

	if err = req.Validate(); err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToValidateRequestBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}

	resp, err := e.gatewayApiService.ResetPassword(c, req)
	if err != nil {
		lgr.Error().Err(err).Msg("failed")
		code, obj := outer.GetHTTPError(err)
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}

	lgr.Debug().Interface("response", resp).Msg("executed")

	buf, err := protojson.MarshalOptions{
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
		UseProtoNames:   true,
	}.Marshal(resp)
	if err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToMarshalResponseBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	c.Data(http.StatusOK, "application/json", buf)
}

func (e *GatewayApiHttpEndpoint) CreateCollection(c *gin.Context) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		e.prom.HttpReqDuration.WithLabelValues("CreateCollection").Observe(duration.Seconds())
	}()

	requestID := utils.AnyToString(c.Value(constants.RequestIdKey))
	lgr := e.lgr.With().
		Str(constants.RequestIdKey, requestID).
		Str("handler", "CreateCollection").Logger()

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToGetRequestBody)
		if c.Writer.Status() == http.StatusRequestEntityTooLarge {
			return
		}
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	defer c.Request.Body.Close()

	req := new(GatewayApiProto.CreateCollectionRequest)
	if err = protojson.Unmarshal(body, req); err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToUnmarshalRequestBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	lgr = lgr.With().Interface("request", req).Logger()

	if err = req.Validate(); err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToValidateRequestBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}

	resp, err := e.gatewayApiService.CreateCollection(c, req)
	if err != nil {
		lgr.Error().Err(err).Msg("failed")
		code, obj := outer.GetHTTPError(err)
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}

	lgr.Debug().Interface("response", resp).Msg("executed")

	buf, err := protojson.MarshalOptions{
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
		UseProtoNames:   true,
	}.Marshal(resp)
	if err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToMarshalResponseBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	c.Data(http.StatusOK, "application/json", buf)
}

func (e *GatewayApiHttpEndpoint) UpdateCollection(c *gin.Context) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		e.prom.HttpReqDuration.WithLabelValues("UpdateCollection").Observe(duration.Seconds())
	}()

	requestID := utils.AnyToString(c.Value(constants.RequestIdKey))
	lgr := e.lgr.With().
		Str(constants.RequestIdKey, requestID).
		Str("handler", "UpdateCollection").Logger()

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToGetRequestBody)
		if c.Writer.Status() == http.StatusRequestEntityTooLarge {
			return
		}
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	defer c.Request.Body.Close()

	req := new(GatewayApiProto.UpdateCollectionRequest)
	if err = protojson.Unmarshal(body, req); err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToUnmarshalRequestBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	lgr = lgr.With().Interface("request", req).Logger()

	if err = req.Validate(); err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToValidateRequestBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}

	resp, err := e.gatewayApiService.UpdateCollection(c, req)
	if err != nil {
		lgr.Error().Err(err).Msg("failed")
		code, obj := outer.GetHTTPError(err)
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}

	lgr.Debug().Interface("response", resp).Msg("executed")

	buf, err := protojson.MarshalOptions{
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
		UseProtoNames:   true,
	}.Marshal(resp)
	if err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToMarshalResponseBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	c.Data(http.StatusOK, "application/json", buf)
}

func (e *GatewayApiHttpEndpoint) GetCollections(c *gin.Context) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		e.prom.HttpReqDuration.WithLabelValues("GetCollections").Observe(duration.Seconds())
	}()

	requestID := utils.AnyToString(c.Value(constants.RequestIdKey))
	lgr := e.lgr.With().
		Str(constants.RequestIdKey, requestID).
		Str("handler", "GetCollections").Logger()

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToGetRequestBody)
		if c.Writer.Status() == http.StatusRequestEntityTooLarge {
			return
		}
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	defer c.Request.Body.Close()

	req := new(GatewayApiProto.GetCollectionsRequest)
	if err = protojson.Unmarshal(body, req); err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToUnmarshalRequestBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	lgr = lgr.With().Interface("request", req).Logger()

	if err = req.Validate(); err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToValidateRequestBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}

	resp, err := e.gatewayApiService.GetCollections(c, req)
	if err != nil {
		lgr.Error().Err(err).Msg("failed")
		code, obj := outer.GetHTTPError(err)
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}

	lgr.Debug().Msg("executed")

	buf, err := protojson.MarshalOptions{
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
		UseProtoNames:   true,
	}.Marshal(resp)
	if err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToMarshalResponseBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	c.Data(http.StatusOK, "application/json", buf)
}

func (e *GatewayApiHttpEndpoint) GetCollection(c *gin.Context) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		e.prom.HttpReqDuration.WithLabelValues("GetCollection").Observe(duration.Seconds())
	}()

	requestID := utils.AnyToString(c.Value(constants.RequestIdKey))
	lgr := e.lgr.With().
		Str(constants.RequestIdKey, requestID).
		Str("handler", "GetCollection").Logger()

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToGetRequestBody)
		if c.Writer.Status() == http.StatusRequestEntityTooLarge {
			return
		}
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	defer c.Request.Body.Close()

	req := new(GatewayApiProto.GetCollectionRequest)
	if err = protojson.Unmarshal(body, req); err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToUnmarshalRequestBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	lgr = lgr.With().Interface("request", req).Logger()

	if err = req.Validate(); err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToValidateRequestBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}

	resp, err := e.gatewayApiService.GetCollection(c, req)
	if err != nil {
		lgr.Error().Err(err).Msg("failed")
		code, obj := outer.GetHTTPError(err)
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}

	lgr.Debug().Interface("response", resp).Msg("executed")

	buf, err := protojson.MarshalOptions{
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
		UseProtoNames:   true,
	}.Marshal(resp)
	if err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToMarshalResponseBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	c.Data(http.StatusOK, "application/json", buf)
}

func (e *GatewayApiHttpEndpoint) DeleteCollection(c *gin.Context) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		e.prom.HttpReqDuration.WithLabelValues("DeleteCollection").Observe(duration.Seconds())
	}()

	requestID := utils.AnyToString(c.Value(constants.RequestIdKey))
	lgr := e.lgr.With().
		Str(constants.RequestIdKey, requestID).
		Str("handler", "DeleteCollection").Logger()

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToGetRequestBody)
		if c.Writer.Status() == http.StatusRequestEntityTooLarge {
			return
		}
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	defer c.Request.Body.Close()

	req := new(GatewayApiProto.DeleteCollectionRequest)
	if err = protojson.Unmarshal(body, req); err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToUnmarshalRequestBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	lgr = lgr.With().Interface("request", req).Logger()

	if err = req.Validate(); err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToValidateRequestBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}

	resp, err := e.gatewayApiService.DeleteCollection(c, req)
	if err != nil {
		lgr.Error().Err(err).Msg("failed")
		code, obj := outer.GetHTTPError(err)
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}

	lgr.Debug().Interface("response", resp).Msg("executed")

	buf, err := protojson.MarshalOptions{
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
		UseProtoNames:   true,
	}.Marshal(resp)
	if err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToMarshalResponseBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	c.Data(http.StatusOK, "application/json", buf)
}

func (e *GatewayApiHttpEndpoint) CreateTerms(c *gin.Context) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		e.prom.HttpReqDuration.WithLabelValues("CreateTerms").Observe(duration.Seconds())
	}()

	requestID := utils.AnyToString(c.Value(constants.RequestIdKey))
	lgr := e.lgr.With().
		Str(constants.RequestIdKey, requestID).
		Str("handler", "CreateTerms").Logger()

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToGetRequestBody)
		if c.Writer.Status() == http.StatusRequestEntityTooLarge {
			return
		}
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	defer c.Request.Body.Close()

	req := new(GatewayApiProto.CreateTermsRequest)
	if err = protojson.Unmarshal(body, req); err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToUnmarshalRequestBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	lgr = lgr.With().Interface("request", req).Logger()

	if err = req.Validate(); err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToValidateRequestBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}

	resp, err := e.gatewayApiService.CreateTerms(c, req)
	if err != nil {
		lgr.Error().Err(err).Msg("failed")
		code, obj := outer.GetHTTPError(err)
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}

	lgr.Debug().Interface("response", resp).Msg("executed")

	buf, err := protojson.MarshalOptions{
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
		UseProtoNames:   true,
	}.Marshal(resp)
	if err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToMarshalResponseBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	c.Data(http.StatusOK, "application/json", buf)
}

func (e *GatewayApiHttpEndpoint) GetTerms(c *gin.Context) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		e.prom.HttpReqDuration.WithLabelValues("GetTerms").Observe(duration.Seconds())
	}()

	requestID := utils.AnyToString(c.Value(constants.RequestIdKey))
	lgr := e.lgr.With().
		Str(constants.RequestIdKey, requestID).
		Str("handler", "GetTerms").Logger()

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToGetRequestBody)
		if c.Writer.Status() == http.StatusRequestEntityTooLarge {
			return
		}
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	defer c.Request.Body.Close()

	req := new(GatewayApiProto.GetTermsRequest)
	if err = protojson.Unmarshal(body, req); err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToUnmarshalRequestBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	lgr = lgr.With().Interface("request", req).Logger()

	if err = req.Validate(); err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToValidateRequestBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}

	resp, err := e.gatewayApiService.GetTerms(c, req)
	if err != nil {
		lgr.Error().Err(err).Msg("failed")
		code, obj := outer.GetHTTPError(err)
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}

	lgr.Debug().Msg("executed")

	buf, err := protojson.MarshalOptions{
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
		UseProtoNames:   true,
	}.Marshal(resp)
	if err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToMarshalResponseBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	c.Data(http.StatusOK, "application/json", buf)
}

func (e *GatewayApiHttpEndpoint) UpdateTerm(c *gin.Context) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		e.prom.HttpReqDuration.WithLabelValues("UpdateTerm").Observe(duration.Seconds())
	}()

	requestID := utils.AnyToString(c.Value(constants.RequestIdKey))
	lgr := e.lgr.With().
		Str(constants.RequestIdKey, requestID).
		Str("handler", "UpdateTerm").Logger()

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToGetRequestBody)
		if c.Writer.Status() == http.StatusRequestEntityTooLarge {
			return
		}
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	defer c.Request.Body.Close()

	req := new(GatewayApiProto.UpdateTermRequest)
	if err = protojson.Unmarshal(body, req); err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToUnmarshalRequestBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	lgr = lgr.With().Interface("request", req).Logger()

	if err = req.Validate(); err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToValidateRequestBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}

	resp, err := e.gatewayApiService.UpdateTerm(c, req)
	if err != nil {
		lgr.Error().Err(err).Msg("failed")
		code, obj := outer.GetHTTPError(err)
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}

	lgr.Debug().Msg("executed")

	buf, err := protojson.MarshalOptions{
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
		UseProtoNames:   true,
	}.Marshal(resp)
	if err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToMarshalResponseBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	c.Data(http.StatusOK, "application/json", buf)
}

func (e *GatewayApiHttpEndpoint) DeleteTerms(c *gin.Context) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		e.prom.HttpReqDuration.WithLabelValues("DeleteTerms").Observe(duration.Seconds())
	}()

	requestID := utils.AnyToString(c.Value(constants.RequestIdKey))
	lgr := e.lgr.With().
		Str(constants.RequestIdKey, requestID).
		Str("handler", "DeleteTerms").Logger()

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToGetRequestBody)
		if c.Writer.Status() == http.StatusRequestEntityTooLarge {
			return
		}
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	defer c.Request.Body.Close()

	req := new(GatewayApiProto.DeleteTermsRequest)
	if err = protojson.Unmarshal(body, req); err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToUnmarshalRequestBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	lgr = lgr.With().Interface("request", req).Logger()

	if err = req.Validate(); err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToValidateRequestBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}

	resp, err := e.gatewayApiService.DeleteTerms(c, req)
	if err != nil {
		lgr.Error().Err(err).Msg("failed")
		code, obj := outer.GetHTTPError(err)
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}

	lgr.Debug().Interface("response", resp).Msg("executed")

	buf, err := protojson.MarshalOptions{
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
		UseProtoNames:   true,
	}.Marshal(resp)
	if err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToMarshalResponseBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	c.Data(http.StatusOK, "application/json", buf)
}

func (e *GatewayApiHttpEndpoint) ChangeTermStatus(c *gin.Context) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		e.prom.HttpReqDuration.WithLabelValues("ChangeTermStatus").Observe(duration.Seconds())
	}()

	requestID := utils.AnyToString(c.Value(constants.RequestIdKey))
	lgr := e.lgr.With().
		Str(constants.RequestIdKey, requestID).
		Str("handler", "ChangeTermStatus").Logger()

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToGetRequestBody)
		if c.Writer.Status() == http.StatusRequestEntityTooLarge {
			return
		}
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	defer c.Request.Body.Close()

	req := new(GatewayApiProto.ChangeTermStatusRequest)
	if err = protojson.Unmarshal(body, req); err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToUnmarshalRequestBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	lgr = lgr.With().Interface("request", req).Logger()

	if err = req.Validate(); err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToValidateRequestBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}

	resp, err := e.gatewayApiService.ChangeTermStatus(c, req)
	if err != nil {
		lgr.Error().Err(err).Msg("failed")
		code, obj := outer.GetHTTPError(err)
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}

	lgr.Debug().Interface("response", resp).Msg("executed")

	buf, err := protojson.MarshalOptions{
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
		UseProtoNames:   true,
	}.Marshal(resp)
	if err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToMarshalResponseBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	c.Data(http.StatusOK, "application/json", buf)
}

func (e *GatewayApiHttpEndpoint) GetLanguages(c *gin.Context) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		e.prom.HttpReqDuration.WithLabelValues("GetLanguages").Observe(duration.Seconds())
	}()

	requestID := utils.AnyToString(c.Value(constants.RequestIdKey))
	lgr := e.lgr.With().
		Str(constants.RequestIdKey, requestID).
		Str("handler", "GetLanguages").Logger()

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToGetRequestBody)
		if c.Writer.Status() == http.StatusRequestEntityTooLarge {
			return
		}
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	defer c.Request.Body.Close()

	req := new(GatewayApiProto.GetLanguagesRequest)
	if err = protojson.Unmarshal(body, req); err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToUnmarshalRequestBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	lgr = lgr.With().Interface("request", req).Logger()

	if err = req.Validate(); err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToValidateRequestBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}

	resp, err := e.gatewayApiService.GetLanguages(c, req)
	if err != nil {
		lgr.Error().Err(err).Msg("failed")
		code, obj := outer.GetHTTPError(err)
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}

	lgr.Debug().Msg("executed")

	buf, err := protojson.MarshalOptions{
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
		UseProtoNames:   true,
	}.Marshal(resp)
	if err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToMarshalResponseBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	c.Data(http.StatusOK, "application/json", buf)
}

func (e *GatewayApiHttpEndpoint) GetVoiceover(c *gin.Context) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		e.prom.HttpReqDuration.WithLabelValues("GetVoiceover").Observe(duration.Seconds())
	}()

	requestID := utils.AnyToString(c.Value(constants.RequestIdKey))
	lgr := e.lgr.With().
		Str(constants.RequestIdKey, requestID).
		Str("handler", "GetVoiceover").Logger()

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToGetRequestBody)
		if c.Writer.Status() == http.StatusRequestEntityTooLarge {
			return
		}
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	defer c.Request.Body.Close()

	req := new(GatewayApiProto.GetVoiceoverRequest)
	if err = protojson.Unmarshal(body, req); err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToUnmarshalRequestBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	lgr = lgr.With().Interface("request", req).Logger()

	if err = req.Validate(); err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToValidateRequestBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}

	resp, err := e.gatewayApiService.GetVoiceover(c, req)
	if err != nil {
		lgr.Error().Err(err).Msg("failed")
		code, obj := outer.GetHTTPError(err)
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}

	lgr.Debug().Interface("response", resp).Msg("executed")

	buf, err := protojson.MarshalOptions{
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
		UseProtoNames:   true,
	}.Marshal(resp)
	if err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToMarshalResponseBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	c.Data(http.StatusOK, "application/json", buf)
}

func (e *GatewayApiHttpEndpoint) GetTranslation(c *gin.Context) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		e.prom.HttpReqDuration.WithLabelValues("GetTranslation").Observe(duration.Seconds())
	}()

	requestID := utils.AnyToString(c.Value(constants.RequestIdKey))
	lgr := e.lgr.With().
		Str(constants.RequestIdKey, requestID).
		Str("handler", "GetTranslation").Logger()

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToGetRequestBody)
		if c.Writer.Status() == http.StatusRequestEntityTooLarge {
			return
		}
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	defer c.Request.Body.Close()

	req := new(GatewayApiProto.GetTranslationRequest)
	if err = protojson.Unmarshal(body, req); err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToUnmarshalRequestBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	lgr = lgr.With().Interface("request", req).Logger()

	if err = req.Validate(); err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToValidateRequestBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}

	resp, err := e.gatewayApiService.GetTranslation(c, req)
	if err != nil {
		lgr.Error().Err(err).Msg("failed")
		code, obj := outer.GetHTTPError(err)
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}

	lgr.Debug().Interface("response", resp).Msg("executed")

	buf, err := protojson.MarshalOptions{
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
		UseProtoNames:   true,
	}.Marshal(resp)
	if err != nil {
		lgr.Error().Err(err).Msg(errors.FailedToMarshalResponseBody)
		code, obj := outer.GetHTTPError(errors.BadRequestError(err))
		e.prom.HttpRespCount.WithLabelValues(strconv.FormatInt(int64(code), 10)).Add(1)
		c.PureJSON(code, obj)
		return
	}
	c.Data(http.StatusOK, "application/json", buf)
}
