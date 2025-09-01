package shared

import (
	"net/http"
	urlpkg "net/url"

	"github.com/dracory/req"
	"github.com/samber/lo"
)

func URL(endpoint string, path string, params map[string]string) string {
	if params == nil {
		params = map[string]string{}
	}
	params["path"] = path

	url := endpoint + query(params)
	return url
}

func URLR(r *http.Request, path string, params map[string]string) string {
	endpoint := Endpoint(r)
	filterSiteID := req.GetStringTrimmed(r, "filter_site_id")

	if params == nil {
		params = map[string]string{}
	}

	params["path"] = path

	if !lo.HasKey(params, "filter_site_id") {
		params["filter_site_id"] = filterSiteID
	}

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
