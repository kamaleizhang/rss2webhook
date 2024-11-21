

# workflow 

config rss source via http endpoint
API will save config to db (sqlite)

rss updater fetch rss feed from source
diff new/update/remove/unchanged items by query db 
save new/update items to db

call webhook take new/update items 