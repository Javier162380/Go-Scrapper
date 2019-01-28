package set

type Pages_to_Crawl struct {
	link map[string]struct{}
}

var exists = struct{}{}

func links_set() *Pages_to_Crawl {
	p := &Pages_to_Crawl{}
	p.link = make(map[string]struct{})
	return p
}

func (p *Pages_to_Crawl) Add_link(link string) {
	p.link[link] = exists
}

func (p *Pages_to_Crawl) Contains(link string) bool {
	_, c := p.link[link]
	return c
}
