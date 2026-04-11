-- WARNING: This schema is for context only and is not meant to be run.
-- Table order and constraints may not be valid for execution.

CREATE TABLE gugu.aliexpress_seller_token (
                                              id text NOT NULL,
                                              seller_id text NOT NULL UNIQUE,
                                              havana_id text NOT NULL DEFAULT ''::text,
                                              app_user_id text NOT NULL DEFAULT ''::text,
                                              user_nick text NOT NULL DEFAULT ''::text,
                                              account text NOT NULL DEFAULT ''::text,
                                              account_platform text NOT NULL DEFAULT ''::text,
                                              locale text NOT NULL DEFAULT ''::text,
                                              sp text NOT NULL DEFAULT ''::text,
                                              access_token text NOT NULL,
                                              refresh_token text NOT NULL,
                                              access_token_expires_at timestamp with time zone NOT NULL,
                                              refresh_token_expires_at timestamp with time zone,
                                              last_refreshed_at timestamp with time zone NOT NULL DEFAULT now(),
                                              authorized_at timestamp with time zone NOT NULL DEFAULT now(),
                                              created_at timestamp with time zone NOT NULL DEFAULT now(),
                                              updated_at timestamp with time zone NOT NULL DEFAULT now(),
                                              app_type text NOT NULL DEFAULT 'AFFILIATE'::text UNIQUE,
                                              CONSTRAINT aliexpress_seller_token_pkey PRIMARY KEY (id)
);
CREATE TABLE gugu.app_user (
                               id text NOT NULL,
                               email text NOT NULL UNIQUE,
                               display_name text NOT NULL DEFAULT ''::text,
                               password_hash text NOT NULL DEFAULT ''::text,
                               auth_source text NOT NULL,
                               email_verified boolean NOT NULL DEFAULT false,
                               email_verified_at timestamp with time zone,
                               created_at timestamp with time zone NOT NULL DEFAULT now(),
                               CONSTRAINT app_user_pkey PRIMARY KEY (id)
);
CREATE TABLE gugu.email_verification (
                                         code text NOT NULL,
                                         user_id text NOT NULL,
                                         email text NOT NULL,
                                         expires_at timestamp with time zone NOT NULL,
                                         used_at timestamp with time zone,
                                         created_at timestamp with time zone NOT NULL DEFAULT now(),
                                         CONSTRAINT email_verification_pkey PRIMARY KEY (code)
);
CREATE TABLE gugu.oauth_identity (
                                     id text NOT NULL,
                                     user_id text NOT NULL,
                                     provider text NOT NULL,
                                     subject text NOT NULL,
                                     email text NOT NULL,
                                     created_at timestamp with time zone NOT NULL DEFAULT now(),
                                     last_login_at timestamp with time zone NOT NULL,
                                     CONSTRAINT oauth_identity_pkey PRIMARY KEY (id)
);
CREATE TABLE gugu.price_alert (
                                  id text NOT NULL,
                                  user_id text NOT NULL,
                                  sku_id text NOT NULL,
                                  channel text NOT NULL DEFAULT 'EMAIL'::text,
                                  enabled boolean NOT NULL DEFAULT true,
                                  created_at timestamp with time zone NOT NULL DEFAULT now(),
                                  CONSTRAINT price_alert_pkey PRIMARY KEY (id),
                                  CONSTRAINT price_alert_user_id_fkey FOREIGN KEY (user_id) REFERENCES gugu.app_user(id),
                                  CONSTRAINT price_alert_sku_id_fkey FOREIGN KEY (sku_id) REFERENCES gugu.sku(id)
);
CREATE TABLE gugu.product (
                              id text NOT NULL,
                              market text NOT NULL,
                              external_product_id text NOT NULL,
                              original_url text NOT NULL DEFAULT ''::text,
                              title text NOT NULL DEFAULT ''::text,
                              main_image_url text NOT NULL DEFAULT ''::text,
                              product_url text NOT NULL DEFAULT ''::text,
                              collection_source text NOT NULL DEFAULT ''::text,
                              last_collected_at timestamp with time zone NOT NULL DEFAULT now(),
                              created_at timestamp with time zone NOT NULL DEFAULT now(),
                              updated_at timestamp with time zone NOT NULL DEFAULT now(),
                              promotion_link text NOT NULL DEFAULT ''::text,
                              CONSTRAINT product_pkey PRIMARY KEY (id)
);
CREATE TABLE gugu.product_external_alias (
                                             id text NOT NULL,
                                             market text NOT NULL,
                                             alias_external_product_id text NOT NULL,
                                             product_id text NOT NULL,
                                             alias_type text NOT NULL DEFAULT 'VIEW'::text,
                                             created_at timestamp with time zone NOT NULL DEFAULT now(),
                                             updated_at timestamp with time zone NOT NULL DEFAULT now(),
                                             CONSTRAINT product_external_alias_pkey PRIMARY KEY (id),
                                             CONSTRAINT fk_product_external_alias_product FOREIGN KEY (product_id) REFERENCES gugu.product(id)
);
CREATE TABLE gugu.product_localization (
                                           product_id text NOT NULL,
                                           language text NOT NULL,
                                           title text NOT NULL DEFAULT ''::text,
                                           main_image_url text NOT NULL DEFAULT ''::text,
                                           product_url text NOT NULL DEFAULT ''::text,
                                           updated_at timestamp with time zone NOT NULL DEFAULT now(),
                                           CONSTRAINT product_localization_pkey PRIMARY KEY (product_id, language),
                                           CONSTRAINT fk_product_localization_product FOREIGN KEY (product_id) REFERENCES gugu.product(id)
);
CREATE TABLE gugu.sku (
                          id text NOT NULL,
                          product_id text NOT NULL,
                          external_sku_id text NOT NULL DEFAULT ''::text,
                          sku_name text NOT NULL DEFAULT ''::text,
                          color text NOT NULL DEFAULT ''::text,
                          size text NOT NULL DEFAULT ''::text,
                          image_url text NOT NULL DEFAULT ''::text,
                          sku_properties text NOT NULL DEFAULT ''::text,
                          created_at timestamp with time zone NOT NULL DEFAULT now(),
                          updated_at timestamp with time zone NOT NULL DEFAULT now(),
                          origin_sku_id text NOT NULL DEFAULT ''::text,
                          CONSTRAINT sku_pkey PRIMARY KEY (id),
                          CONSTRAINT fk_product_sku_product FOREIGN KEY (product_id) REFERENCES gugu.product(id)
);
CREATE TABLE gugu.sku_localization (
                                       sku_id text NOT NULL,
                                       language text NOT NULL,
                                       sku_name text NOT NULL DEFAULT ''::text,
                                       color_name text NOT NULL DEFAULT ''::text,
                                       size_name text NOT NULL DEFAULT ''::text,
                                       sku_properties text NOT NULL DEFAULT ''::text,
                                       image_url text NOT NULL DEFAULT ''::text,
                                       updated_at timestamp with time zone NOT NULL DEFAULT now(),
                                       CONSTRAINT sku_localization_pkey PRIMARY KEY (sku_id, language),
                                       CONSTRAINT fk_sku_localization_sku FOREIGN KEY (sku_id) REFERENCES gugu.sku(id)
);
CREATE TABLE gugu.sku_price_history (
                                        sku_id text NOT NULL,
                                        recorded_at timestamp with time zone NOT NULL,
                                        price text NOT NULL DEFAULT ''::text,
                                        currency text NOT NULL DEFAULT ''::text,
                                        change_value text NOT NULL DEFAULT ''::text,
                                        CONSTRAINT sku_price_history_pkey PRIMARY KEY (sku_id, currency, recorded_at),
                                        CONSTRAINT fk_sku_price_history_sku FOREIGN KEY (sku_id) REFERENCES gugu.sku(id)
);
CREATE TABLE gugu.sku_price_snapshot (
                                         sku_id text NOT NULL,
                                         snapshot_date date NOT NULL,
                                         price text NOT NULL DEFAULT ''::text,
                                         original_price text NOT NULL DEFAULT ''::text,
                                         currency text NOT NULL DEFAULT ''::text,
                                         CONSTRAINT sku_price_snapshot_pkey PRIMARY KEY (sku_id, currency, snapshot_date),
                                         CONSTRAINT fk_sku_price_snapshot_sku FOREIGN KEY (sku_id) REFERENCES gugu.sku(id)
);
CREATE TABLE gugu.sku_price_snapshot_staging (
                                                 run_id text NOT NULL,
                                                 sku_id text NOT NULL,
                                                 snapshot_date date NOT NULL,
                                                 price text NOT NULL DEFAULT ''::text,
                                                 original_price text NOT NULL DEFAULT ''::text,
                                                 currency text NOT NULL,
                                                 CONSTRAINT sku_price_snapshot_staging_pkey PRIMARY KEY (run_id, sku_id, currency),
                                                 CONSTRAINT fk_sku_price_snapshot_staging_run FOREIGN KEY (run_id) REFERENCES gugu.sku_snapshot_ingest_run(id),
                                                 CONSTRAINT fk_sku_price_snapshot_staging_sku FOREIGN KEY (sku_id) REFERENCES gugu.sku(id)
);
CREATE TABLE gugu.sku_snapshot_ingest_run (
                                              id text NOT NULL,
                                              product_id text NOT NULL,
                                              currency text NOT NULL,
                                              snapshot_date date NOT NULL,
                                              expected_sku_count integer NOT NULL DEFAULT 0,
                                              collected_sku_count integer NOT NULL DEFAULT 0,
                                              status text NOT NULL,
                                              started_at timestamp with time zone NOT NULL DEFAULT now(),
                                              finished_at timestamp with time zone,
                                              error_message text NOT NULL DEFAULT ''::text,
                                              CONSTRAINT sku_snapshot_ingest_run_pkey PRIMARY KEY (id),
                                              CONSTRAINT fk_sku_snapshot_ingest_run_product FOREIGN KEY (product_id) REFERENCES gugu.product(id)
);
CREATE TABLE gugu.user_login_session (
                                         id text NOT NULL,
                                         user_id text NOT NULL,
                                         refresh_token_hash text NOT NULL,
                                         token_family_id text NOT NULL,
                                         parent_session_id text,
                                         user_agent text NOT NULL DEFAULT ''::text,
                                         client_ip text NOT NULL DEFAULT ''::text,
                                         device_name text NOT NULL DEFAULT ''::text,
                                         expires_at timestamp with time zone NOT NULL,
                                         last_seen_at timestamp with time zone NOT NULL DEFAULT now(),
                                         rotated_at timestamp with time zone,
                                         revoked_at timestamp with time zone,
                                         reuse_detected_at timestamp with time zone,
                                         created_at timestamp with time zone NOT NULL DEFAULT now(),
                                         CONSTRAINT user_login_session_pkey PRIMARY KEY (id),
                                         CONSTRAINT fk_user_login_sessions_user FOREIGN KEY (user_id) REFERENCES gugu.app_user(id),
                                         CONSTRAINT fk_user_login_sessions_parent FOREIGN KEY (parent_session_id) REFERENCES gugu.user_login_session(id)
);
CREATE TABLE gugu.user_tracked_item (
                                        id text NOT NULL,
                                        user_id text NOT NULL,
                                        product_id text NOT NULL,
                                        original_url text NOT NULL DEFAULT ''::text,
                                        created_at timestamp with time zone NOT NULL DEFAULT now(),
                                        deleted_at timestamp with time zone,
                                        sku_id text,
                                        currency text NOT NULL DEFAULT 'KRW'::text,
                                        view_external_product_id text NOT NULL DEFAULT ''::text,
                                        preferred_language text NOT NULL DEFAULT 'KO'::text,
                                        tracking_scope text NOT NULL DEFAULT 'PRODUCT_ALL_SKU'::text,
                                        CONSTRAINT user_tracked_item_pkey PRIMARY KEY (id),
                                        CONSTRAINT fk_user_tracked_items_user FOREIGN KEY (user_id) REFERENCES gugu.app_user(id),
                                        CONSTRAINT fk_user_tracked_items_product FOREIGN KEY (product_id) REFERENCES gugu.product(id)
);
CREATE TABLE gugu.user_tracked_item_watch_sku (
                                                  tracked_item_id text NOT NULL,
                                                  sku_id text NOT NULL,
                                                  created_at timestamp with time zone NOT NULL DEFAULT now(),
                                                  CONSTRAINT user_tracked_item_watch_sku_pkey PRIMARY KEY (tracked_item_id, sku_id),
                                                  CONSTRAINT fk_user_tracked_item_watch_sku_tracked_item FOREIGN KEY (tracked_item_id) REFERENCES gugu.user_tracked_item(id),
                                                  CONSTRAINT fk_user_tracked_item_watch_sku_sku FOREIGN KEY (sku_id) REFERENCES gugu.sku(id)
);