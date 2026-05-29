package Message

func CreateErrorMessage(message string, statusCode int) map[string]interface{} {

	return map[string]interface{}{
		"statusCode": statusCode,
		"data": map[string]string{
			"message": message,
		},
	}
}
