CREATE TABLE IF NOT EXISTS publications (
    guid uuid NOT NULL,
    category_id uuid NOT NULL,
    author_id uuid NOT NULL,
    title character varying,
    description character varying,
    type character varying NOT NULL,
    created_at timestamp without time zone DEFAULT now(),
    CONSTRAINT publications_pkey PRIMARY KEY (guid)
);

ALTER TABLE publications ADD CONSTRAINT "publications_category_id_fkey" FOREIGN KEY(category_id) REFERENCES categories("guid"); 
ALTER TABLE publications ADD CONSTRAINT "publications_author_id_fkey" FOREIGN KEY(author_id) REFERENCES authors("guid"); 