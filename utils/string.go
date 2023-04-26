package utils

import "net/url"

func MapToQueryString(m map[string]string) string {
	values := make(url.Values)
	for k, v := range m {
		values.Add(k, v)
	}
	return values.Encode()
}

func GetQueryStringFromURL(input string) (map[string]string, error) {
	u, err := url.Parse(input)
	if err != nil {
		return nil, err
	}
	m := make(map[string]string)
	query := u.RawQuery
	values, err := url.ParseQuery(query)
	if err != nil {
		return nil, err
	}
	for k, v := range values {
		m[k] = v[0]
	}
	return m, nil
}
