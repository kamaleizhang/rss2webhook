package r2w

import (
	"crypto/sha256"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/mmcdole/gofeed"
	"github.com/samber/lo"
)

type RssUpdater interface {
	SyncRss(Config) ([]HookRequest, error)
}

type RssUpdaterImpl struct {
	repo *sqlx.DB
}

func NewRssUpdaterImpl(db *sqlx.DB) RssUpdater {
	return &RssUpdaterImpl{
		repo: db,
	}
}

func (r *RssUpdaterImpl) SyncRss(config Config) ([]HookRequest, error) {
	log.Printf("Syncing rss: %v", config.RssURL)
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(config.RssURL)
	if err != nil {
		return nil, err
	}
	log.Printf("Feed %s Updated: %v", config.RssURL, feed.Updated)

	ids := make([]string, 0, len(feed.Items))
	for _, item := range feed.Items {
		ids = append(ids, item.GUID)
	}
	fetchFromDB, err := r.fetchRssFromDB(ids, config.Domain, config.RssURL)
	if err != nil {
		return nil, err
	}
	rssFromDBMap := lo.KeyBy(fetchFromDB, func(rss RssFeed) string {
		return rss.EntryId
	})

	var result []HookRequest
	var updatedRss []RssFeed
	for _, item := range feed.Items {
		hashValue := hash(item.Content)
		var eventType string
		if dbRss, ok := rssFromDBMap[item.GUID]; ok && dbRss.ContentHash == hashValue {
			eventType = "update"
		} else {
			eventType = "new"
		}
		result = append(result, NewHookRequest(eventType, config.Domain, config.RssURL, item))
		updatedRss = append(updatedRss, NewRssFeed(config.Domain, config.RssURL, item.GUID, item.Updated, hashValue))
	}

	err = r.saveRssToDB(updatedRss)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// sql: select * from rss_feed where entry_id in (?,?,?) and domain =? and  rss_url = ?
func (r *RssUpdaterImpl) fetchRssFromDB(ids []string, domain string, rssURL string) ([]RssFeed, error) {
	query, args, err := sqlx.In("select * from rss_feed where entry_id in (?) and domain =? and rss_url =?",
		ids, domain, rssURL)
	if err != nil {
		log.Fatalf("Error preparing query: %v", err)
	}
	query = r.repo.Rebind(query)
	var rssFeeds []RssFeed
	err = r.repo.Select(&rssFeeds, query, args...)
	if err != nil {
		return nil, err
	}
	return rssFeeds, nil
}

func (r *RssUpdaterImpl) saveRssToDB(result []RssFeed) error {
	tx, err := r.repo.Beginx()
	if err != nil {
		return err
	}
	upsertSQL := `
	INSERT OR REPLACE INTO rss_feed (domain, rss_url, entry_id, last_updated, content_hash)
	VALUES (:domain, :rss_url, :entry_id，:last_updated, :content_hash)
	`
	for _, rss := range result {
		_, err := tx.NamedExec(upsertSQL, rss)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func hash(s string) string {
	bytes := []byte(s)

	hash := sha256.New()

	hash.Write(bytes)

	hashValue := hash.Sum(nil)

	hashString := fmt.Sprintf("%x", hashValue)

	return hashString
}
