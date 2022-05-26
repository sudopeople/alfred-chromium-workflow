package src

import (
	"fmt"
	"github.com/deanishe/awgo"
)

var FetchSearchData = func (wf *aw.Workflow, query string) {
	var whereStmt string
	titleQuery, domainQuery, isDomainSearch, _, _:= HandleUserQuery(query)

	if (isDomainSearch) {
		whereStmt = fmt.Sprintf("WHERE urls.url LIKE '%%%s%%' AND keyword_search_terms.term LIKE '%%%s%%'", domainQuery, titleQuery)
	} else {
		whereStmt = fmt.Sprintf("WHERE keyword_search_terms.term LIKE '%%%s%%'", titleQuery)
	}

	var dbQuery = fmt.Sprintf(`
		SELECT urls.url, urls.last_visit_time, keyword_search_terms.term
			FROM keyword_search_terms
			JOIN urls ON urls.id = keyword_search_terms.url_id
			%s
		ORDER BY last_visit_time DESC
	`, whereStmt)

	historyDB := GetHistoryDB()

	rows, err := historyDB.Query(dbQuery)
	CheckError(err)

	var urlTitle string
	var url string
	var urlLastVisitTime int64

	var itemCount = 0
	var previousTitle = ""

	for rows.Next() {
		if itemCount >= int(Conf.ResultLimitCount) {
			break
		}

		err := rows.Scan(&url, &urlLastVisitTime, &urlTitle)
		CheckError(err)

		if previousTitle == urlTitle {
			continue
		}

		domainName := ExtractDomainName(url)
		unixTimestamp := ConvertChromeTimeToUnixTimestamp(urlLastVisitTime)
		localeTimeStr := GetLocaleString(unixTimestamp)

		item := wf.NewItem(urlTitle).
			Subtitle(fmt.Sprintf(`From '%s', In '%s'`, domainName, localeTimeStr)).
			Valid(true).
			Quicklook(url).
			Var("type", "url").
			Var("url", url).
			Copytext(url).
			Largetype(url)

		iconPath := fmt.Sprintf(`cache/%s`, domainName)

		if FileExist(iconPath) {
			item.Icon(&aw.Icon{iconPath, ""})
		}

		previousTitle = urlTitle
		itemCount += 1
	}
}
