package utils

func GetMapValue(key string, customMap map[string]interface{}) interface{} {
	if value, found := customMap[key]; found {
		return value
	}
	return nil
}
