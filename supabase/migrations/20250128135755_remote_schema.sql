alter table "public"."ai_models" drop constraint "valid_model_type";

alter table "public"."api_keys" drop constraint "rate_limit_range";

alter table "public"."model_requests" drop constraint "valid_status";

alter table "public"."api_keys" drop constraint "api_keys_user_id_fkey";

alter table "public"."model_requests" drop constraint "model_requests_model_id_fkey";

alter table "public"."model_requests" drop constraint "model_requests_user_id_fkey";

drop index if exists "public"."users_email_idx";

alter table "public"."ai_models" drop column "function_url";

alter table "public"."ai_models" add column "functionURL" text;

alter table "public"."ai_models" alter column "created_at" set default now();

alter table "public"."ai_models" alter column "created_at" set not null;

alter table "public"."ai_models" alter column "huggingface_id" set default ''::text;

alter table "public"."ai_models" alter column "huggingface_id" set data type text using "huggingface_id"::text;

alter table "public"."ai_models" alter column "id" set default gen_random_uuid();

alter table "public"."ai_models" alter column "model_type" set default ''::text;

alter table "public"."ai_models" alter column "model_type" set data type text using "model_type"::text;

alter table "public"."ai_models" alter column "name" set default ''::text;

alter table "public"."ai_models" alter column "name" set data type text using "name"::text;

alter table "public"."ai_models" alter column "version" set default ''::text;

alter table "public"."ai_models" alter column "version" set data type text using "version"::text;

alter table "public"."api_keys" drop column "last_used";

alter table "public"."api_keys" add column "last_update" timestamp with time zone;

alter table "public"."api_keys" alter column "created_at" set default now();

alter table "public"."api_keys" alter column "created_at" set not null;

alter table "public"."api_keys" alter column "id" set default gen_random_uuid();

alter table "public"."api_keys" alter column "key_hash" set default ''::text;

alter table "public"."api_keys" alter column "name" drop not null;

alter table "public"."api_keys" alter column "name" set data type text using "name"::text;

alter table "public"."api_keys" alter column "rate_limit" drop default;

alter table "public"."api_keys" alter column "rate_limit" drop not null;

alter table "public"."api_keys" alter column "rate_limit" set data type bigint using "rate_limit"::bigint;

alter table "public"."model_requests" drop column "token_count";

alter table "public"."model_requests" drop column "token_used";

alter table "public"."model_requests" add column "tokens_used" bigint default '0'::bigint;

alter table "public"."model_requests" add column "tokent_count" bigint;

alter table "public"."model_requests" alter column "api_key_id" drop not null;

alter table "public"."model_requests" alter column "created_at" set default now();

alter table "public"."model_requests" alter column "created_at" set not null;

alter table "public"."model_requests" alter column "id" set default gen_random_uuid();

alter table "public"."model_requests" alter column "input_data" drop not null;

alter table "public"."model_requests" alter column "model_id" drop not null;

alter table "public"."model_requests" alter column "processing_time" drop default;

alter table "public"."model_requests" alter column "processing_time" set data type bigint using "processing_time"::bigint;

alter table "public"."model_requests" alter column "status" set default 'pending'::text;

alter table "public"."model_requests" alter column "status" drop not null;

alter table "public"."model_requests" alter column "status" set data type text using "status"::text;

alter table "public"."model_requests" alter column "user_id" drop not null;

alter table "public"."users" alter column "created_at" drop default;

alter table "public"."users" alter column "created_at" set not null;

alter table "public"."users" alter column "failed_login_attempts" set default '0'::bigint;

alter table "public"."users" alter column "failed_login_attempts" set data type bigint using "failed_login_attempts"::bigint;

alter table "public"."users" disable row level security;

CREATE UNIQUE INDEX ai_models_model_type_key ON public.ai_models USING btree (model_type);

CREATE UNIQUE INDEX ai_models_name_key ON public.ai_models USING btree (name);

CREATE UNIQUE INDEX ai_models_version_key ON public.ai_models USING btree (version);

CREATE UNIQUE INDEX api_keys_user_id_key ON public.api_keys USING btree (user_id);

CREATE UNIQUE INDEX users_email_key ON public.users USING btree (email);

alter table "public"."ai_models" add constraint "ai_models_model_type_key" UNIQUE using index "ai_models_model_type_key";

alter table "public"."ai_models" add constraint "ai_models_name_key" UNIQUE using index "ai_models_name_key";

alter table "public"."ai_models" add constraint "ai_models_version_key" UNIQUE using index "ai_models_version_key";

alter table "public"."api_keys" add constraint "api_keys_user_id_key" UNIQUE using index "api_keys_user_id_key";

alter table "public"."users" add constraint "users_email_key" UNIQUE using index "users_email_key";

alter table "public"."api_keys" add constraint "api_keys_user_id_fkey" FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE not valid;

alter table "public"."api_keys" validate constraint "api_keys_user_id_fkey";

alter table "public"."model_requests" add constraint "model_requests_model_id_fkey" FOREIGN KEY (model_id) REFERENCES ai_models(id) ON DELETE CASCADE not valid;

alter table "public"."model_requests" validate constraint "model_requests_model_id_fkey";

alter table "public"."model_requests" add constraint "model_requests_user_id_fkey" FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE not valid;

alter table "public"."model_requests" validate constraint "model_requests_user_id_fkey";

create policy "Enable insert for authenticated users only"
on "public"."users"
as permissive
for insert
to authenticated
with check (true);


create policy "Enable read access for all users"
on "public"."users"
as permissive
for select
to public
using (true);


create policy "Enable update for users based on email"
on "public"."users"
as permissive
for update
to public
using (((( SELECT auth.jwt() AS jwt) ->> 'email'::text) = email))
with check (((( SELECT auth.jwt() AS jwt) ->> 'email'::text) = email));


create policy "Enable users to view their own data only"
on "public"."users"
as permissive
for select
to authenticated
using ((( SELECT auth.uid() AS uid) = id));



