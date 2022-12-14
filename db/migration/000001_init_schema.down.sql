/* drop entries and transfers tables first since there are foreign key referencing to accounts table */
DROP TABLE IF EXISTS entries;
DROP TABLE IF EXISTS transfers;
DROP TABLE IF EXISTS accounts;