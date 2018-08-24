package enrichers

func SetBody(output []byte, data map[string][]byte, params map[string]string) ([]byte, error) {
	output = data[`body`]
	return output, nil
}

func AddHeaders(output []byte, data map[string][]byte, params map[string]string) ([]byte, error) {
	output = append(output, data[`headers`]...)
	return output, nil
}
