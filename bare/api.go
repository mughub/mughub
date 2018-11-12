package bare

// RegisterAPIEndpoint registers the GraphQL API endpoint with the provided router.
func RegisterAPIEndpoint(r Router) {
	r.Methods("GET", "POST").Path("/graphql")
	// TODO
}
