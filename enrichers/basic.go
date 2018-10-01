package enrichers

func SetBody(output *enrichedData, data map[string][]byte, params map[string]string) (*enrichedData, error) {
	output.Body = string(data[`data`])
	return output, nil
}

func AddHeaders(output *enrichedData, data map[string][]byte, params map[string]string) (*enrichedData, error) {
	output.Headers = string(data[`headers`])
	return output, nil
}
