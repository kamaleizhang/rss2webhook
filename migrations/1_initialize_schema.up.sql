create table rss_feeds (
    domain text not null,
    rss_url text not null,
    entry_id text not null,
    last_updated text not null,
    content_hash text not null,
);