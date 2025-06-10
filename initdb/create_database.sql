-- !SETUP-COMPOSE: make sure database name is following DATABSE_NAME on .env file
SELECT 'CREATE DATABASE video_rental_shop_db'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'video_rental_shop_db')\gexec