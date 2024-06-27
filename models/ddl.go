package models

const refreshTokensDDL = ` 
CREATE TABLE IF NOT EXISTS refresh_tokens (
	id serial4 NOT NULL,
	device_id int4 NOT NULL,
	user_id int4 NOT NULL,
	expires_at int4 NOT NULL,
    CONSTRAINT refresh_pk PRIMARY KEY (id)
);
`

const usersDDL = `
CREATE TABLE IF NOT EXISTS users (
	id serial4 NOT NULL,
	phone varchar(18) NOT NULL,
    roles varchar(16)[] NOT NULL,
    CONSTRAINT users_pk PRIMARY KEY (id),
    CONSTRAINT users_unique UNIQUE (phone)
);`

const smsCodesDDL = `
CREATE TABLE IF NOT EXISTS sms_codes (
    user_id serial4 NOT NULL,
    code int4 NOT NULL,
    sent_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    CONSTRAINT codes_uid_unique UNIQUE (user_id),
    CONSTRAINT fk_uid
        FOREIGN KEY(user_id) 
            REFERENCES users(id)
            ON UPDATE CASCADE
            ON DELETE CASCADE
);
`

var tableDDLS = []string{
	usersDDL,
	smsCodesDDL,
	refreshTokensDDL,
}
