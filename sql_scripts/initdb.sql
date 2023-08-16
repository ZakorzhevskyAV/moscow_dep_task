CREATE TABLE IF NOT EXISTS UserData(
user_id varchar(100),
data jsonb,
time TIMESTAMP DEFAULT NOW()
);

