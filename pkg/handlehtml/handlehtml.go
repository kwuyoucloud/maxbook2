package handlehtml

// PageInfo store pagenum and pagedownlink
type PageInfo struct {
	PageNum  string
	PageLink string
}

// GetPageNumAndPageSrclink can return pagenum pagelink list, if there is something wrong, return error
func GetPageNumAndPageSrclink(bodystr string) ([]PageInfo, error) {
	var pageInfo []PageInfo

	return pageInfo, nil
}
