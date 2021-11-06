package query

type GetUserQuery struct {
	UserIDs []string `query:"user_id"`
}
