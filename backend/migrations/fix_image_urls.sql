-- Fix absolute URLs in raffle_images table
-- Replaces 'http://localhost:8080' with empty string to make URLs relative

UPDATE raffle_images 
SET 
    url_original = REPLACE(url_original, 'http://localhost:8080', ''),
    url_large = REPLACE(url_large, 'http://localhost:8080', ''),
    url_medium = REPLACE(url_medium, 'http://localhost:8080', ''),
    url_thumbnail = REPLACE(url_thumbnail, 'http://localhost:8080', '');

-- Verify changes
SELECT id, url_original, url_thumbnail FROM raffle_images LIMIT 5;
