package responses

func ErrorBadRequests(m string) (int, map[string]interface{}) {
	return 400, map[string]interface{}{
		"status":  "Error",
		"message": m,
	}
}

func ErrorNotFound(m string) (int, map[string]interface{}) {
	return 404, map[string]interface{}{
		"status":  "Error",
		"message": m,
	}
}

func ErrorConflict(m string) (int, map[string]interface{}) {
	return 409, map[string]interface{}{
		"status":  "Error",
		"message": m,
	}
}

func ErrorServer(m string) (int, map[string]interface{}) {
	return 500, map[string]interface{}{
		"status":  "Error",
		"message": m,
	}
}

func ErrorServers(m string) (int, map[string]interface{}) {
	return 500, map[string]interface{}{
		"status":  "Error",
		"message": m,
	}
}
