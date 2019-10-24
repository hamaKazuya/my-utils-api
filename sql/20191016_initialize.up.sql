CREATE TABLE todos (
    id INTEGER auto_increment primary key,
    title VARCHAR(50),
    is_done INTEGER NOT NULL DEFAULT 0,
    detail VARCHAR(40),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO todos (title, is_done, detail) VALUES ('コーヒーを買う', 0, 'アイスカフェオレで');
INSERT INTO todos (title, is_done, detail) VALUES ('犬の散歩', 0, '吠えないように気をつける');