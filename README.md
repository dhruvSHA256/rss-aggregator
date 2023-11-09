gateway 
    create user
    add/follow feeds
authsvc
    used by gateway
scraper scheduler
scraper worker



user profile service:
    user signup
    auth
    add feeds
aggregator scheduler service
    after every 10 min add msg to kafka queue to start workers with list of feed ids
worker node
    listen to kafka and scrap feeds
    sends message to kafka about scrapping completed
processing service
    reads msg from kafka about scrapped services
    saves it into db
    add msg to kafka to nofify user if email enabled
nofity service
    reads msg from kafka about email notification and sends email
