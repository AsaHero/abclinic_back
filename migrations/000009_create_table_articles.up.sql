CREATE TABLE IF NOT EXISTS articles (
    guid uuid NOT NULL,
    chapter_id uuid NOT NULL,
    info character varying,
    img character varying,
    side character varying,
    created_at timestamp without time zone DEFAULT now(),
    CONSTRAINT articles_pkey PRIMARY KEY (guid)
);

ALTER TABLE articles ADD CONSTRAINT "articles_chapter_id_fkey" FOREIGN KEY(chapter_id) REFERENCES chapters("guid");