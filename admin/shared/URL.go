package shared

import urlpkg "net/url"

func URL(endpoint string, path string, params map[string]string) string {
	if params == nil {
		params = map[string]string{}
	}
	params["path"] = path

	url := endpoint + query(params)
	return url
}

func query(queryData map[string]string) string {
	queryString := ""

	if len(queryData) > 0 {
		v := urlpkg.Values{}
		for key, value := range queryData {
			v.Set(key, value)
		}
		queryString += "?" + httpBuildQuery(v)
	}

	return queryString
}

func httpBuildQuery(queryData urlpkg.Values) string {
	return queryData.Encode()
}
