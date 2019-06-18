package utils

//ConvertToMapInterface []interface{} ==> []map[string]interface{}
func ConvertToMapInterface(datas []interface{}) []map[string]interface{} {
	var result []map[string]interface{}
	for _, data := range datas {
		result = append(result, data.(map[string]interface{}))
	}
	return result
}

//ConvertToSliceString []interface{} ==> []string
func ConvertToSliceString(datas []interface{}) []string {
	var result []string
	for _, data := range datas {
		result = append(result, data.(string))
	}
	return result
}
