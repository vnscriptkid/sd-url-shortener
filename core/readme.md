# Core

## Setup
- DynamoDB
    - Create table
        - Table name: `urlShortener`
        - Primary key: `shortCode` (String)
        - Default settings for the rest
        - Create table
    - Create IAM user
        - IAM > Users -> Create user
            - User name: `dynamodb-user`
            - Attach policies directly
                - AmazonDynamoDBFullAccess
            - Create user
        - Go user details
            - Security credentials
                - Create access key
                - Grab the access key and secret key