package errors

import (
	"github.com/golang/protobuf/proto"
	"gitlab.com/wordbyword.io/microservices/pkg/errors/outer"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"net/http"
)

var (
	TooManyRequests = &outer.OuterError{
		ErrorMessage:   TooManyRequestsMsg,
		HttpStatusCode: http.StatusTooManyRequests,
		GrpcStatusCode: codes.Unavailable,
	}
	TokenClaimsDoesNotSet = &outer.OuterError{
		ErrorMessage:   TokenClaimsDoesNotSetMsg,
		HttpStatusCode: http.StatusUnauthorized,
		GrpcStatusCode: codes.PermissionDenied,
	}
	FailedToGetCollectionPrivateCollection = &outer.OuterError{
		ErrorMessage:   FailedToGetCollectionPrivateCollectionMsg,
		HttpStatusCode: http.StatusForbidden,
		GrpcStatusCode: codes.PermissionDenied,
	}

	FailedToGenerateTokens = &outer.OuterError{
		ErrorMessage:   FailedToGenerateTokensMsg,
		HttpStatusCode: http.StatusBadRequest,
		GrpcStatusCode: codes.InvalidArgument,
	}
	FailedToDeleteTokens = &outer.OuterError{
		ErrorMessage:   FailedToDeleteTokensMsg,
		HttpStatusCode: http.StatusBadRequest,
		GrpcStatusCode: codes.InvalidArgument,
	}
	FailedToRefreshTokens = &outer.OuterError{
		ErrorMessage:   FailedToRefreshTokensMsg,
		HttpStatusCode: http.StatusBadRequest,
		GrpcStatusCode: codes.InvalidArgument,
	}
	FailedToCreateUser = &outer.OuterError{
		ErrorMessage:   FailedToCreateUserMsg,
		HttpStatusCode: http.StatusBadRequest,
		GrpcStatusCode: codes.InvalidArgument,
	}
	FailedToUpdateUser = &outer.OuterError{
		ErrorMessage:   FailedToUpdateUserMsg,
		HttpStatusCode: http.StatusBadRequest,
		GrpcStatusCode: codes.InvalidArgument,
	}
	FailedToConfirmEmail = &outer.OuterError{
		ErrorMessage:   FailedToConfirmEmailMsg,
		HttpStatusCode: http.StatusBadRequest,
		GrpcStatusCode: codes.InvalidArgument,
	}
	FailedToResetPassword = &outer.OuterError{
		ErrorMessage:   FailedToResetPasswordMsg,
		HttpStatusCode: http.StatusBadRequest,
		GrpcStatusCode: codes.InvalidArgument,
	}
	FailedToGetUser = &outer.OuterError{
		ErrorMessage:   FailedToGetUserMsg,
		HttpStatusCode: http.StatusBadRequest,
		GrpcStatusCode: codes.InvalidArgument,
	}
	FailedToSignIn = &outer.OuterError{
		ErrorMessage:   FailedToSignInMsg,
		HttpStatusCode: http.StatusBadRequest,
		GrpcStatusCode: codes.InvalidArgument,
	}
	FailedToSignInWrongPassword = &outer.OuterError{
		ErrorMessage:   FailedToSignInWrongPasswordMsg,
		HttpStatusCode: http.StatusBadRequest,
		GrpcStatusCode: codes.InvalidArgument,
	}
	FailedToCreateAction = &outer.OuterError{
		ErrorMessage:   FailedToCreateActionMsg,
		HttpStatusCode: http.StatusBadRequest,
		GrpcStatusCode: codes.InvalidArgument,
	}
	FailedToSendConfirmationEmail = &outer.OuterError{
		ErrorMessage:   FailedToSendConfirmationEmailMsg,
		HttpStatusCode: http.StatusBadRequest,
		GrpcStatusCode: codes.InvalidArgument,
	}
	FailedToSendResetPasswordEmail = &outer.OuterError{
		ErrorMessage:   FailedToSendResetPasswordEmailMsg,
		HttpStatusCode: http.StatusBadRequest,
		GrpcStatusCode: codes.InvalidArgument,
	}
	FailedToGetLanguages = &outer.OuterError{
		ErrorMessage:   FailedToGetLanguagesMsg,
		HttpStatusCode: http.StatusBadRequest,
		GrpcStatusCode: codes.InvalidArgument,
	}
	FailedToGetCollections = &outer.OuterError{
		ErrorMessage:   FailedToGetCollectionsMsg,
		HttpStatusCode: http.StatusBadRequest,
		GrpcStatusCode: codes.InvalidArgument,
	}
	FailedToGetCollection = &outer.OuterError{
		ErrorMessage:   FailedToGetCollectionMsg,
		HttpStatusCode: http.StatusBadRequest,
		GrpcStatusCode: codes.InvalidArgument,
	}
	FailedToGetTranslation = &outer.OuterError{
		ErrorMessage:   FailedToGetTranslationMsg,
		HttpStatusCode: http.StatusBadRequest,
		GrpcStatusCode: codes.InvalidArgument,
	}
	FailedToGetTerms = &outer.OuterError{
		ErrorMessage:   FailedToGetTermsMsg,
		HttpStatusCode: http.StatusBadRequest,
		GrpcStatusCode: codes.InvalidArgument,
	}
	FailedToDeleteTerms = &outer.OuterError{
		ErrorMessage:   FailedToDeleteTermsMsg,
		HttpStatusCode: http.StatusBadRequest,
		GrpcStatusCode: codes.InvalidArgument,
	}
	FailedToDeleteCollection = &outer.OuterError{
		ErrorMessage:   FailedToDeleteCollectionMsg,
		HttpStatusCode: http.StatusBadRequest,
		GrpcStatusCode: codes.InvalidArgument,
	}
	FailedToCreateCollection = &outer.OuterError{
		ErrorMessage:   FailedToCreateCollectionMsg,
		HttpStatusCode: http.StatusBadRequest,
		GrpcStatusCode: codes.InvalidArgument,
	}
	FailedToUpdateCollection = &outer.OuterError{
		ErrorMessage:   FailedToUpdateCollectionMsg,
		HttpStatusCode: http.StatusBadRequest,
		GrpcStatusCode: codes.InvalidArgument,
	}
	FailedToUpdateTerm = &outer.OuterError{
		ErrorMessage:   FailedToUpdateTermMsg,
		HttpStatusCode: http.StatusBadRequest,
		GrpcStatusCode: codes.InvalidArgument,
	}
	FailedToCreateTerms = &outer.OuterError{
		ErrorMessage:   FailedToCreateTermsMsg,
		HttpStatusCode: http.StatusBadRequest,
		GrpcStatusCode: codes.InvalidArgument,
	}
	FailedToChangeTermStatus = &outer.OuterError{
		ErrorMessage:   FailedToChangeTermStatusMsg,
		HttpStatusCode: http.StatusBadRequest,
		GrpcStatusCode: codes.InvalidArgument,
	}
	FailedToCreateTermsCollection = &outer.OuterError{
		ErrorMessage:   FailedToCreateTermsCollectionMsg,
		HttpStatusCode: http.StatusBadRequest,
		GrpcStatusCode: codes.InvalidArgument,
	}
	FailedToGetVoiceover = &outer.OuterError{
		ErrorMessage:   FailedToGetVoiceoverMsg,
		HttpStatusCode: http.StatusBadRequest,
		GrpcStatusCode: codes.InvalidArgument,
	}
	FailedToGetGoogleUser = &outer.OuterError{
		ErrorMessage:   FailedToGetGoogleUserMsg,
		HttpStatusCode: http.StatusBadRequest,
		GrpcStatusCode: codes.InvalidArgument,
	}
	FailedToSignInEmailNotConfirmed = &outer.OuterError{
		ErrorMessage:   FailedToSignInEmailNotConfirmedMsg,
		HttpStatusCode: http.StatusBadRequest,
		GrpcStatusCode: codes.InvalidArgument,
	}
)

func BadRequestError(err error) *outer.OuterError {
	return &outer.OuterError{
		ErrorMessage:   "incorrect request",
		HttpStatusCode: http.StatusBadRequest,
		GrpcStatusCode: codes.InvalidArgument,
		GrpcDetails: []proto.Message{
			&errdetails.ErrorInfo{
				Metadata: map[string]string{"error": err.Error()},
			},
		},
	}
}

func InternalError(err error) *outer.OuterError {
	return &outer.OuterError{
		ErrorMessage:   "internal error",
		HttpStatusCode: http.StatusInternalServerError,
		GrpcStatusCode: codes.Internal,
		GrpcDetails: []proto.Message{
			&errdetails.ErrorInfo{
				Metadata: map[string]string{"error": err.Error()},
			},
		},
	}
}

func BadAuthorizationTokenError(errMsg string) *outer.OuterError {
	return &outer.OuterError{
		ErrorMessage:   errMsg,
		HttpStatusCode: http.StatusUnauthorized,
		GrpcStatusCode: codes.PermissionDenied,
	}
}

func AlreadyExistError(errMsg string) *outer.OuterError {
	return &outer.OuterError{
		ErrorMessage:   errMsg,
		HttpStatusCode: http.StatusBadRequest,
		GrpcStatusCode: codes.AlreadyExists,
	}
}
