package common

func ResultCommon(code int, data interface{}, message ...string) map[string]interface{} {
	result := map[string]interface{}{
		"code":    code,
		"data":    data,
		"message": message[0],
	}
	return result
}
