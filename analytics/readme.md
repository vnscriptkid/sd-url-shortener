# Analytics

## Setup
- Go kibana: `http://localhost:5601`
- Go to `Management` -> `Index Patterns` -> `Create index pattern`
- Go `Analytics` -> `Visualize` -> `Create visualization`

## Visualizations
- Total Visits:
    - Aggregation based -> Metric
    - Set Aggregation to Count.

- Unique Visitors:
    - Aggregation based -> Metric
    - Set Aggregation to `Unique Count`.
    - Select ip_address field.

- Top Referrers:
    - Lens -> Bar
    - Set Y-axis to Count.
    - Set X-axis to Terms aggregation on referrer.

- Geo-distribution:
    - Create a Coordinate Map.
    - Set Aggregation to Geo coordinates (if geo data is available).

- Visit Trends:
    - Create a Line chart.
    - Set X-axis to Date Histogram on visited_at.
    - Set Y-axis to Count.