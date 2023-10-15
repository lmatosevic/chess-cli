package command

import "strconv"

const HomeDirName = ".chess-cli"

const AccessTokenFile = HomeDirName + "/.access_token"

const ServerHostFile = HomeDirName + "/.server_host"

func BuildQueryParams(page int, size int, sort string, filter string) map[string]string {
	params := make(map[string]string)
	if page > 0 {
		params["page"] = strconv.Itoa(page)
	}
	if size > 0 {
		params["size"] = strconv.Itoa(size)
	}
	if sort != "" {
		params["sort"] = sort
	}
	if filter != "" {
		params["filter"] = filter
	}
	return params
}
