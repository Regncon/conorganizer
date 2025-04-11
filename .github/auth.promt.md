Always add the AuthMiddleware to authRouter or similar names that has a protectedRoute name or similar using route(where u use the context name).Use(service.AuthMiddleware(logger))
