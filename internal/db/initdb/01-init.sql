-- Create extension
CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;
COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';

-- Create sequences
CREATE SEQUENCE public.goadmin_menu_myid_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    MAXVALUE 99999999
    CACHE 1;

CREATE SEQUENCE public.goadmin_operation_log_myid_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    MAXVALUE 99999999
    CACHE 1;

CREATE SEQUENCE public.goadmin_site_myid_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    MAXVALUE 99999999
    CACHE 1;

CREATE SEQUENCE public.goadmin_permissions_myid_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    MAXVALUE 99999999
    CACHE 1;

CREATE SEQUENCE public.goadmin_roles_myid_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    MAXVALUE 99999999
    CACHE 1;

CREATE SEQUENCE public.goadmin_session_myid_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    MAXVALUE 99999999
    CACHE 1;

CREATE SEQUENCE public.goadmin_users_myid_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    MAXVALUE 99999999
    CACHE 1;

-- Create tables
CREATE TABLE public.goadmin_menu (
                                     id integer DEFAULT nextval('public.goadmin_menu_myid_seq'::regclass) NOT NULL,
                                     parent_id integer DEFAULT 0 NOT NULL,
                                     type integer DEFAULT 0,
                                     "order" integer DEFAULT 0 NOT NULL,
                                     title character varying(50) NOT NULL,
                                     header character varying(100),
                                     plugin_name character varying(100) NOT NULL,
                                     icon character varying(50) NOT NULL,
                                     uri character varying(3000) NOT NULL,
                                     uuid character varying(100),
                                     created_at timestamp without time zone DEFAULT now(),
                                     updated_at timestamp without time zone DEFAULT now()
);

CREATE TABLE public.goadmin_operation_log (
                                              id integer DEFAULT nextval('public.goadmin_operation_log_myid_seq'::regclass) NOT NULL,
                                              user_id integer NOT NULL,
                                              path character varying(255) NOT NULL,
                                              method character varying(10) NOT NULL,
                                              ip character varying(15) NOT NULL,
                                              input text NOT NULL,
                                              created_at timestamp without time zone DEFAULT now(),
                                              updated_at timestamp without time zone DEFAULT now()
);

CREATE TABLE public.goadmin_site (
                                     id integer DEFAULT nextval('public.goadmin_site_myid_seq'::regclass) NOT NULL,
                                     key character varying(100) NOT NULL,
                                     value text NOT NULL,
                                     type integer DEFAULT 0,
                                     description character varying(3000),
                                     state integer DEFAULT 0,
                                     created_at timestamp without time zone DEFAULT now(),
                                     updated_at timestamp without time zone DEFAULT now()
);

CREATE TABLE public.goadmin_permissions (
                                            id integer DEFAULT nextval('public.goadmin_permissions_myid_seq'::regclass) NOT NULL,
                                            name character varying(50) NOT NULL,
                                            slug character varying(50) NOT NULL,
                                            http_method character varying(255),
                                            http_path text NOT NULL,
                                            created_at timestamp without time zone DEFAULT now(),
                                            updated_at timestamp without time zone DEFAULT now()
);

CREATE TABLE public.goadmin_role_menu (
                                          role_id integer NOT NULL,
                                          menu_id integer NOT NULL,
                                          created_at timestamp without time zone DEFAULT now(),
                                          updated_at timestamp without time zone DEFAULT now()
);

CREATE TABLE public.goadmin_role_permissions (
                                                 role_id integer NOT NULL,
                                                 permission_id integer NOT NULL,
                                                 created_at timestamp without time zone DEFAULT now(),
                                                 updated_at timestamp without time zone DEFAULT now()
);

CREATE TABLE public.goadmin_role_users (
                                           role_id integer NOT NULL,
                                           user_id integer NOT NULL,
                                           created_at timestamp without time zone DEFAULT now(),
                                           updated_at timestamp without time zone DEFAULT now()
);

CREATE TABLE public.goadmin_roles (
                                      id integer DEFAULT nextval('public.goadmin_roles_myid_seq'::regclass) NOT NULL,
                                      name character varying NOT NULL,
                                      slug character varying NOT NULL,
                                      created_at timestamp without time zone DEFAULT now(),
                                      updated_at timestamp without time zone DEFAULT now()
);

CREATE TABLE public.goadmin_session (
                                        id integer DEFAULT nextval('public.goadmin_session_myid_seq'::regclass) NOT NULL,
                                        sid character varying(50) NOT NULL,
                                        "values" character varying(3000) NOT NULL,
                                        created_at timestamp without time zone DEFAULT now(),
                                        updated_at timestamp without time zone DEFAULT now()
);

CREATE TABLE public.goadmin_user_permissions (
                                                 user_id integer NOT NULL,
                                                 permission_id integer NOT NULL,
                                                 created_at timestamp without time zone DEFAULT now(),
                                                 updated_at timestamp without time zone DEFAULT now()
);

CREATE TABLE public.goadmin_users (
                                      id integer DEFAULT nextval('public.goadmin_users_myid_seq'::regclass) NOT NULL,
                                      username character varying(100) NOT NULL,
                                      password character varying(100) NOT NULL,
                                      name character varying(100) NOT NULL,
                                      avatar character varying(255),
                                      remember_token character varying(100),
                                      created_at timestamp without time zone DEFAULT now(),
                                      updated_at timestamp without time zone DEFAULT now()
);

-- Insert data
INSERT INTO public.goadmin_menu (id, parent_id, type, "order", title, header, plugin_name, icon, uri, created_at, updated_at) VALUES
                                                                                                                                  (1, 0, 1, 2, 'Admin', NULL, '', 'fa-tasks', '', '2019-09-10 00:00:00', '2019-09-10 00:00:00'),
                                                                                                                                  (2, 1, 1, 2, 'Users', NULL, '', 'fa-users', '/info/manager', '2019-09-10 00:00:00', '2019-09-10 00:00:00'),
                                                                                                                                  (3, 1, 1, 3, 'Roles', NULL, '', 'fa-user', '/info/roles', '2019-09-10 00:00:00', '2019-09-10 00:00:00'),
                                                                                                                                  (4, 1, 1, 4, 'Permission', NULL, '', 'fa-ban', '/info/permission', '2019-09-10 00:00:00', '2019-09-10 00:00:00'),
                                                                                                                                  (5, 1, 1, 5, 'Menu', NULL, '', 'fa-bars', '/menu', '2019-09-10 00:00:00', '2019-09-10 00:00:00'),
                                                                                                                                  (6, 1, 1, 6, 'Operation log', NULL, '', 'fa-history', '/info/op', '2019-09-10 00:00:00', '2019-09-10 00:00:00'),
                                                                                                                                  (7, 0, 1, 1, 'Dashboard', NULL, '', 'fa-bar-chart', '/', '2019-09-10 00:00:00', '2019-09-10 00:00:00');

INSERT INTO public.goadmin_permissions (id, name, slug, http_method, http_path, created_at, updated_at) VALUES
                                                                                                            (1, 'All permission', '*', NULL, '*', '2019-09-10 00:00:00', '2019-09-10 00:00:00'),
                                                                                                            (2, 'Dashboard', 'dashboard', 'GET,PUT,POST,DELETE', '/', '2019-09-10 00:00:00', '2019-09-10 00:00:00');

INSERT INTO public.goadmin_roles (id, name, slug, created_at, updated_at) VALUES
                                                                              (1, 'Administrator', 'administrator', '2019-09-10 00:00:00', '2019-09-10 00:00:00'),
                                                                              (2, 'Operator', 'operator', '2019-09-10 00:00:00', '2019-09-10 00:00:00');

INSERT INTO public.goadmin_users (id, username, password, name, avatar, remember_token, created_at, updated_at) VALUES
                                                                                                                    (1, 'admin', '$2a$10$OxWYJJGTP2gi00l2x06QuOWqw5VR47MQCJ0vNKnbMYfrutij10Hwe', 'admin', NULL, 'tlNcBVK9AvfYH7WEnwB1RKvocJu8FfRy4um3DJtwdHuJy0dwFsLOgAc0xUfh', '2019-09-10 00:00:00', '2019-09-10 00:00:00'),
                                                                                                                    (2, 'operator', '$2a$10$rVqkOzHjN2MdlEprRflb1eGP0oZXuSrbJLOmJagFsCd81YZm0bsh.', 'Operator', NULL, NULL, '2019-09-10 00:00:00', '2019-09-10 00:00:00');

INSERT INTO public.goadmin_role_menu (role_id, menu_id, created_at, updated_at) VALUES
                                                                                    (1, 1, '2019-09-10 00:00:00', '2019-09-10 00:00:00'),
                                                                                    (1, 7, '2019-09-10 00:00:00', '2019-09-10 00:00:00'),
                                                                                    (2, 7, '2019-09-10 00:00:00', '2019-09-10 00:00:00');

INSERT INTO public.goadmin_role_permissions (role_id, permission_id, created_at, updated_at) VALUES
                                                                                                 (1, 1, '2019-09-10 00:00:00', '2019-09-10 00:00:00'),
                                                                                                 (1, 2, '2019-09-10 00:00:00', '2019-09-10 00:00:00'),
                                                                                                 (2, 2, '2019-09-10 00:00:00', '2019-09-10 00:00:00');

INSERT INTO public.goadmin_role_users (role_id, user_id, created_at, updated_at) VALUES
                                                                                     (1, 1, '2019-09-10 00:00:00', '2019-09-10 00:00:00'),
                                                                                     (2, 2, '2019-09-10 00:00:00', '2019-09-10 00:00:00');

INSERT INTO public.goadmin_user_permissions (user_id, permission_id, created_at, updated_at) VALUES
                                                                                                 (1, 1, '2019-09-10 00:00:00', '2019-09-10 00:00:00'),
                                                                                                 (2, 2, '2019-09-10 00:00:00', '2019-09-10 00:00:00');

-- Set sequence values
SELECT setval('public.goadmin_menu_myid_seq', 7, true);
SELECT setval('public.goadmin_operation_log_myid_seq', 1, true);
SELECT setval('public.goadmin_permissions_myid_seq', 2, true);
SELECT setval('public.goadmin_roles_myid_seq', 2, true);
SELECT setval('public.goadmin_site_myid_seq', 1, true);
SELECT setval('public.goadmin_session_myid_seq', 1, true);
SELECT setval('public.goadmin_users_myid_seq', 2, true);

-- Add primary keys
ALTER TABLE ONLY public.goadmin_menu
    ADD CONSTRAINT goadmin_menu_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.goadmin_operation_log
    ADD CONSTRAINT goadmin_operation_log_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.goadmin_permissions
    ADD CONSTRAINT goadmin_permissions_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.goadmin_roles
    ADD CONSTRAINT goadmin_roles_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.goadmin_site
    ADD CONSTRAINT goadmin_site_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.goadmin_session
    ADD CONSTRAINT goadmin_session_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.goadmin_users
    ADD CONSTRAINT goadmin_users_pkey PRIMARY KEY (id);

-- Set ownership
ALTER TABLE public.goadmin_menu OWNER TO postgres;
ALTER TABLE public.goadmin_operation_log OWNER TO postgres;
ALTER TABLE public.goadmin_site OWNER TO postgres;
ALTER TABLE public.goadmin_permissions OWNER TO postgres;
ALTER TABLE public.goadmin_role_menu OWNER TO postgres;
ALTER TABLE public.goadmin_role_permissions OWNER TO postgres;
ALTER TABLE public.goadmin_role_users OWNER TO postgres;
ALTER TABLE public.goadmin_roles OWNER TO postgres;
ALTER TABLE public.goadmin_session OWNER TO postgres;
ALTER TABLE public.goadmin_user_permissions OWNER TO postgres;
ALTER TABLE public.goadmin_users OWNER TO postgres;

-- Set permissions
REVOKE ALL ON SCHEMA public FROM PUBLIC;
REVOKE ALL ON SCHEMA public FROM postgres;
GRANT ALL ON SCHEMA public TO postgres;
GRANT ALL ON SCHEMA public TO PUBLIC;
