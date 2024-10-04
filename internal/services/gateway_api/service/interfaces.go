package service

import (
	"context"
	"github.com/rs/zerolog"
	ActionApiProto "gitlab.com/wbwapis/go-genproto/wbw/action/action_api/v1"
	AuthApiProto "gitlab.com/wbwapis/go-genproto/wbw/auth/auth_api/v1"
	GatewayApiProto "gitlab.com/wbwapis/go-genproto/wbw/gateway/gateway_api/v1"
	LanguageApiProto "gitlab.com/wbwapis/go-genproto/wbw/language/language_api/v1"
	NotificationApiProto "gitlab.com/wbwapis/go-genproto/wbw/notification/notification_api/v1"
	SpeakerApiProto "gitlab.com/wbwapis/go-genproto/wbw/speaker/speaker_api/v1"
	TranslationApiProto "gitlab.com/wbwapis/go-genproto/wbw/translation/translation_api/v1"
	UserApiProto "gitlab.com/wbwapis/go-genproto/wbw/user/user_api/v1"
	VocabularyApiProto "gitlab.com/wbwapis/go-genproto/wbw/vocabulary/vocabulary_api/v1"
	"gitlab.com/wordbyword.io/microservices/gateways/gateway-api/internal/entities"
)

// -----------------------
// SERVICE INTERFACES
// -----------------------

type IGatewayApiService interface {
	SignUp(context.Context, *GatewayApiProto.SignUpRequest) (*GatewayApiProto.SignUpResponse, error)
	SignIn(context.Context, *GatewayApiProto.SignInRequest) (*GatewayApiProto.SignInResponse, error)
	Logout(context.Context, *GatewayApiProto.LogoutRequest) (*GatewayApiProto.LogoutResponse, error)
	RefreshTokens(context.Context, *GatewayApiProto.RefreshTokensRequest) (*GatewayApiProto.RefreshTokensResponse, error)
	GetUser(context.Context, *GatewayApiProto.GetUserRequest) (*GatewayApiProto.GetUserResponse, error)
	UpdateUser(context.Context, *GatewayApiProto.UpdateUserRequest) (*GatewayApiProto.UpdateUserResponse, error)
	ConfirmEmail(context.Context, *GatewayApiProto.ConfirmEmailRequest) (*GatewayApiProto.ConfirmEmailResponse, error)
	ResetPassword(context.Context, *GatewayApiProto.ResetPasswordRequest) (*GatewayApiProto.ResetPasswordResponse, error)
	AskResetPassword(context.Context, *GatewayApiProto.AskResetPasswordRequest) (*GatewayApiProto.AskResetPasswordResponse, error)
	CreateCollection(context.Context, *GatewayApiProto.CreateCollectionRequest) (*GatewayApiProto.CreateCollectionResponse, error)
	UpdateCollection(context.Context, *GatewayApiProto.UpdateCollectionRequest) (*GatewayApiProto.UpdateCollectionResponse, error)
	GetCollections(context.Context, *GatewayApiProto.GetCollectionsRequest) (*GatewayApiProto.GetCollectionsResponse, error)
	GetCollection(context.Context, *GatewayApiProto.GetCollectionRequest) (*GatewayApiProto.GetCollectionResponse, error)
	DeleteCollection(context.Context, *GatewayApiProto.DeleteCollectionRequest) (*GatewayApiProto.DeleteCollectionResponse, error)
	CreateTerms(context.Context, *GatewayApiProto.CreateTermsRequest) (*GatewayApiProto.CreateTermsResponse, error)
	GetTerms(context.Context, *GatewayApiProto.GetTermsRequest) (*GatewayApiProto.GetTermsResponse, error)
	UpdateTerm(context.Context, *GatewayApiProto.UpdateTermRequest) (*GatewayApiProto.UpdateTermResponse, error)
	DeleteTerms(context.Context, *GatewayApiProto.DeleteTermsRequest) (*GatewayApiProto.DeleteTermsResponse, error)
	ChangeTermStatus(context.Context, *GatewayApiProto.ChangeTermStatusRequest) (*GatewayApiProto.ChangeTermStatusResponse, error)
	GetLanguages(context.Context, *GatewayApiProto.GetLanguagesRequest) (*GatewayApiProto.GetLanguagesResponse, error)
	GetVoiceover(context.Context, *GatewayApiProto.GetVoiceoverRequest) (*GatewayApiProto.GetVoiceoverResponse, error)
	GetTranslation(context.Context, *GatewayApiProto.GetTranslationRequest) (*GatewayApiProto.GetTranslationResponse, error)
}

// -----------------------
// CLIENTS INTERFACES
// -----------------------

type IVocabularyApiClient interface {
	ChangeTermStatus(context.Context, *VocabularyApiProto.ChangeTermStatusRequest) (*VocabularyApiProto.ChangeTermStatusResponse, error)
	CreateCollection(context.Context, *VocabularyApiProto.CreateCollectionRequest) (*VocabularyApiProto.CreateCollectionResponse, error)
	UpdateCollection(context.Context, *VocabularyApiProto.UpdateCollectionRequest) (*VocabularyApiProto.UpdateCollectionResponse, error)
	CreateTerms(context.Context, *VocabularyApiProto.CreateTermsRequest) (*VocabularyApiProto.CreateTermsResponse, error)
	DeleteCollection(context.Context, *VocabularyApiProto.DeleteCollectionRequest) (*VocabularyApiProto.DeleteCollectionResponse, error)
	DeleteTerms(context.Context, *VocabularyApiProto.DeleteTermsRequest) (*VocabularyApiProto.DeleteTermsResponse, error)
	GetCollection(context.Context, *VocabularyApiProto.GetCollectionRequest) (*VocabularyApiProto.GetCollectionResponse, error)
	GetCollections(context.Context, *VocabularyApiProto.GetCollectionsRequest) (*VocabularyApiProto.GetCollectionsResponse, error)
	GetTerms(context.Context, *VocabularyApiProto.GetTermsRequest) (*VocabularyApiProto.GetTermsResponse, error)
	UpdateTerm(context.Context, *VocabularyApiProto.UpdateTermRequest) (*VocabularyApiProto.UpdateTermResponse, error)
	Shutdown()
}

type IUserApiClient interface {
	GetUser(context.Context, *UserApiProto.GetUserRequest) (*UserApiProto.GetUserResponse, error)
	GetUserByCredentials(context.Context, *UserApiProto.GetUserByCredentialsRequest) (*UserApiProto.GetUserByCredentialsResponse, error)
	CreateUser(context.Context, *UserApiProto.CreateUserRequest) (*UserApiProto.CreateUserResponse, error)
	UpdateUser(context.Context, *UserApiProto.UpdateUserRequest) (*UserApiProto.UpdateUserResponse, error)
	Shutdown()
}

type ISpeakerApiClient interface {
	GetVoiceover(context.Context, *SpeakerApiProto.GetVoiceoverRequest) (*SpeakerApiProto.GetVoiceoverResponse, error)
	Shutdown()
}

type INotificationApiClient interface {
	SendEmail(context.Context, zerolog.Logger, *NotificationApiProto.SendEmailRequest) error
	Shutdown()
}

type ILanguageApiClient interface {
	GetLanguages(context.Context, *LanguageApiProto.GetLanguagesRequest) (*LanguageApiProto.GetLanguagesResponse, error)
	Shutdown()
}

type ITranslationApiClient interface {
	GetTranslation(context.Context, *TranslationApiProto.GetTranslationRequest) (*TranslationApiProto.GetTranslationResponse, error)
	Shutdown()
}

type IGoogleAuthClient interface {
	GetGoogleUser(cxt context.Context, accessToken string) (*entities.GoogleAuthUser, error)
	Shutdown()
}

type IAuthApiClient interface {
	DeleteTokens(context.Context, *AuthApiProto.DeleteTokensRequest) (*AuthApiProto.DeleteTokensResponse, error)
	GenerateTokens(context.Context, *AuthApiProto.GenerateTokensRequest) (*AuthApiProto.GenerateTokensResponse, error)
	RefreshTokens(context.Context, *AuthApiProto.RefreshTokensRequest) (*AuthApiProto.RefreshTokensResponse, error)
	Shutdown()
}

type IActionApiClient interface {
	CreateAction(context.Context, *ActionApiProto.CreateActionRequest) (*ActionApiProto.CreateActionResponse, error)
	ExecuteAction(context.Context, *ActionApiProto.ExecuteActionRequest) (*ActionApiProto.ExecuteActionResponse, error)
	Shutdown()
}
