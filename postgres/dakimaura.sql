--
-- PostgreSQL database dump
--

-- Dumped from database version 17.0
-- Dumped by pg_dump version 17.0

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: products; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.products (
    id integer NOT NULL,
    category character varying(20) NOT NULL,
    name character varying(50) NOT NULL,
    price integer,
    material character varying(20) NOT NULL,
    brand character varying(20) NOT NULL,
    produce_time character varying(10) NOT NULL,
    image character varying(20) NOT NULL
);


ALTER TABLE public.products OWNER TO admin;

--
-- Name: products_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.products_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.products_id_seq OWNER TO admin;

--
-- Name: products_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.products_id_seq OWNED BY public.products.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.users (
    id integer NOT NULL,
    username character varying(20) NOT NULL,
    password character(60) NOT NULL
);


ALTER TABLE public.users OWNER TO admin;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.users_id_seq OWNER TO admin;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: products id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.products ALTER COLUMN id SET DEFAULT nextval('public.products_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Data for Name: products; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.products (id, category, name, price, material, brand, produce_time, image) FROM stdin;
1	Новинки сезона	Доки Доки	6999	Полар-Болс	ВОЕНМЕМ	4-5 суток	doki.jpg
8	Кошкодевочки	Ahri	12799	Полар-Болс	ВОЕНМЕМ	5-6 суток	ahri.jpg
10	Кошкодевочки	Raphtalia	6000	Полар-Болс	ВОЕНМЕМ	2 недели	enot.jpg
13	Новинки сезона	Keqing	4599	Кожа	ВОЕНМЕМ	5 суток	daki2.jpg
14	Кошкодевочки	НЕшкольница	9599	Кожа	dakimakura17	3-4 суток	daki1.png
15	Мужские персонажи	Malchik	7000	Полар-Болс	dakimakura17	8-10 суток	malchik.png
16	Мужские персонажи	Вампир	5999	Полар-Болс	ВОЕНМЕХ	3 недели	vampire.jpeg
17	Мужские персонажи	Волейболист	8000	Полар-Болс	АНИМАНИЯ	7-8 суток	polo.jpg
6	Кошкодевочки	Arknight	8399	Кожа	ВОЕНМЕМ	2-3 суток	arknight.jpg
18	Новинки сезона	Kantai	5699	Полар-Болс	ВОЕНМЕМ	10 суток	kantai.jpg
11	Мужские персонажи	Redhead	8999	Полар-Болс	ВОЕНМЕМ	1-4 суток	cowgirl.jpg
19	Новинки сезона	Monika	5000	Полар-Болс	Doki Doki	7 лет	monika.jpg
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.users (id, username, password) FROM stdin;
1	12345	$2a$10$yjpWlZ5nhTc7iqxBmJqqkelq78yn3eINznHi4WNRsHcwg82P2t2Pm
3	admin	$2a$10$pt9/WP.uk6IV9EFdB4LXoefOEIWDnyts9gY10KoHKl/O7fcYUn7vq
\.


--
-- Name: products_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.products_id_seq', 19, true);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.users_id_seq', 5, true);


--
-- Name: products products_image_key; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_image_key UNIQUE (image);


--
-- Name: products products_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: users users_username_key; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_username_key UNIQUE (username);


--
-- PostgreSQL database dump complete
--

