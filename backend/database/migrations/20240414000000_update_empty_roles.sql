UPDATE "members" 
SET "role" = 'UNSUBSCRIBER' 
WHERE "role" IS NULL OR "role" = ''; 