# Rss Aggregator Service

<img src="https://skillicons.dev/icons?i=go,postgresql,rabbitmq,kubernetes,docker" alt="https://skillicons.dev/icons?i=go,postgresql,rabbitmq,kubernetes,docker" /> 

</br>
</br>

- A microservices app to get your favourite blogposts directly to your email
- Backend is written in golang.
- Uses [Postgres](https://www.postgresql.org/) for data storage,
and [RabbitMQ](https://www.rabbitmq.com/) as a queue between different services

## Overview
- gateway 
    Responsible for user management, feed following, and acting as an entry point for external interactions.
- scheduler
    Handles the scheduling of feed scraping tasks at regular intervals.
- worker
    Performs the actual scraping of RSS feeds and stores the data in a database.
