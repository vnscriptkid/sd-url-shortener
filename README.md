# sd-url-shortener

## Ref:
- https://www.sundog-education.com/lesson/url-shortener-application/

![Screenshot 2024-05-26 at 3 30 47 AM](https://github.com/vnscriptkid/sd-url-shortener/assets/28957748/596d254a-83be-49d3-8db0-89bff0a536dc)

- https://youtu.be/Zcv_899yqhI?si=j932uHISt-y-nu3I

https://app.excalidraw.com/l/56zGeHiLyKZ/5WmP4rHX6Hq

## Setup
- Links
    - https://upstash.com/blog/kafka-url-shortener
    - redis: https://console.upstash.com/redis
    - kafka: https://console.upstash.com/kafka
    - materialize: https://console.materialize.com/regions/aws-eu-west-1/connections
    - cloudflare: https://dash.cloudflare.com/8ae389d2c984641cf30148710f54abe2/workers/services/view/cf-url-shortener/production?versionFilter=all
    - admin: https://cf-url-shortener.vnscriptkid.workers.dev/admin
    - visit: https://cf-url-shortener.vnscriptkid.workers.dev/s/dt

- Deploy server
    - sudo npx wrangler login
    - sudo npx wrangler secret list
    - sudo npx wrangler secret put UPSTASH_KAFKA_REST_URL
    - sudo npx wrangler deploy
    - sudo npx wrangler tail

```sql
CREATE SECRET kafka_password AS 'xxx';

CREATE CONNECTION kafka_connection TO KAFKA (
    BROKER 'cool-civet-11505-us1-kafka.upstash.io:9092',
    SASL MECHANISMS = 'SCRAM-SHA-256',
    SASL USERNAME = 'xxx',
    SASL PASSWORD = SECRET kafka_password
);

CREATE SOURCE click_stats
  FROM KAFKA CONNECTION kafka_connection (TOPIC 'visits-log')
  FORMAT JSON;


CREATE VIEW click_stats_v AS
    SELECT
        (data->>'shortCode')::string AS short_code,
        (data->>'longUrl')::string AS long_url,
        (data->>'country')::string AS country,
        (data->>'city')::string AS city,
        (data->>'ip')::string AS ip
    FROM click_stats;

CREATE MATERIALIZED VIEW click_stats_m AS
    SELECT
        *
    FROM click_stats_v;

SELECT * FROM click_stats_m;

-- order by the number of clicks per short link
CREATE MATERIALIZED VIEW order_by_clicks AS
    SELECT
        short_code,
        COUNT(*) AS clicks
    FROM click_stats_m
    GROUP BY short_code;

-- stream updates from the materialized view as they happen
COPY ( SUBSCRIBE ( SELECT * FROM order_by_clicks ) ) TO STDOUT;

-- Clean up
DROP MATERIALIZED VIEW click_stats_m;
DROP VIEW click_stats_v;
```