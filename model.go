package r2w

type Config struct {
	Domain        string
	RssURL        string
	TargetWebhook string
}

type RssFeed struct {
	Domain      string `db:"domain"`
	RssURL      string `db:"rss_url"`
	EntryId     string `db:"entry_id"`
	LastUpdated string `db:"last_updated"`
	ContentHash string `db:"content_hash"`
}

func NewRssFeed(domain, rssUrl, entryId, lastUpdated, contentHash string) RssFeed {
	return RssFeed{
		Domain:      domain,
		RssURL:      rssUrl,
		EntryId:     entryId,
		LastUpdated: lastUpdated,
		ContentHash: contentHash,
	}

}

type HookRequest struct {
	EventType string
	Domain    string
	RssURL    string
	Content   any
}

func NewHookRequest(eventType, domain, rssUrl string, content any) HookRequest {
	return HookRequest{
		EventType: eventType,
		Domain:    domain,
		RssURL:    rssUrl,
		Content:   content,
	}
}
