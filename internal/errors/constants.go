package errors

const (
	TooManyRequestsMsg                        = "too many requests"
	FailedToGetRequestBody                    = "failed to get request body"
	FailedToUnmarshalRequestBody              = "failed to unmarshal request body"
	FailedToMarshalResponseBody               = "failed to marshal response body"
	FailedToValidateRequestBody               = "failed to validate request body"
	FailedToCreateUserMsg                     = "failed to create user"
	FailedToUpdateUserMsg                     = "failed to update user"
	FailedToGetUserMsg                        = "failed to get user"
	FailedToRefreshTokensMsg                  = "failed to refresh tokens"
	FailedToGenerateTokensMsg                 = "failed to generate tokens"
	FailedToDeleteTokensMsg                   = "failed to delete tokens"
	FailedToSignInMsg                         = "failed to sign in"
	FailedToSignInWrongPasswordMsg            = "failed to sign in, wrong password"
	TokenClaimsDoesNotSetMsg                  = "the token is not set"
	FailedToConfirmEmailMsg                   = "failed to confirm email"
	FailedToResetPasswordMsg                  = "failed to reset password"
	FailedToCreateActionMsg                   = "failed to create action"
	FailedToSendConfirmationEmailMsg          = "failed to send confirmation email"
	FailedToSendResetPasswordEmailMsg         = "failed to send reset password email"
	FailedToGetLanguagesMsg                   = "failed to get languages"
	FailedToGetCollectionsMsg                 = "failed to get collections"
	FailedToGetCollectionMsg                  = "failed to get collection"
	FailedToGetTranslationMsg                 = "failed to get translation"
	FailedToGetTermsMsg                       = "failed to get terms"
	FailedToGetCollectionPrivateCollectionMsg = "private collection"
	FailedToDeleteTermsMsg                    = "failed to delete terms"
	FailedToChangeTermStatusMsg               = "failed to change term status"
	FailedToDeleteCollectionMsg               = "failed to delete collection"
	FailedToCreateCollectionMsg               = "failed to create collection"
	FailedToUpdateCollectionMsg               = "failed to update collection"
	FailedToUpdateTermMsg                     = "failed to update term"
	FailedToCreateTermsMsg                    = "failed to create terms"
	FailedToCreateTermsCollectionMsg          = "collection_id not same"
	FailedToGetVoiceoverMsg                   = "failed to get voiceover"
	FailedToSignInEmailNotConfirmedMsg        = "failed to sign in, email did not confirmed"
	FailedToGetGoogleUserMsg                  = "failed to get google user"
)

// Authorization error messages
const (
	FailedToGetAuthorizationToken = "failed to get authorization token"
	AuthorizationTokenIsEmpty     = "authorization token is empty"
	AuthorizationTokenIsInvalid   = "authorization token is invalid"
)

const (
	FailedToGetUserFormDB = "rpc error: code = NotFound desc = failed to get user form DB"
)
