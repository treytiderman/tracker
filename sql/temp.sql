-- sqlite3 ./data/sql.db < ./sql/temp.sql
SELECT * 
FROM tracker 
WHERE tracker_notes LIKE '%sp%' ESCAPE '\';
