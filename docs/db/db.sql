CREATE sequence next_id;

CREATE
OR REPLACE FUNCTION public.next_id()
 RETURNS bigint
 LANGUAGE plpgsql
AS $function$
DECLARE
seq_id bigint;
BEGIN
SELECT nextval('public.next_id')
INTO seq_id;
return seq_id;
END;
$function$
;

CREATE TABLE public.user
(
    user_id   int8         NOT NULL DEFAULT next_id(),
    fullname  varchar(80)  NOT NULL,
    no_hp     varchar(20)  NOT NULL,
    email     varchar(200) NOT NULL,
    username  varchar(20)  NOT NULL,
    passwd    varchar(200) NOT NULL,
    is_active bool         NOT NULL,
    create_by int8         NOT NULL,
    create_dt timestamptz(0) NOT NULL,
    update_by int8         NOT NULL,
    update_dt timestamptz(0) NOT NULL,
    delete_by int8 NULL,
    delete_dt timestamptz(0) NULL,
    CONSTRAINT user_pk PRIMARY KEY (user_id)
);

CREATE TABLE public.userpokemon
(
    userpokemon_id int8        NOT NULL DEFAULT next_id(),
    user_id        int8        NOT NULL,
    pokemon        varchar(80) NOT NULL,
    nickname       varchar(80) NOT NULL,
    create_by      int8        NOT NULL,
    create_dt      timestamptz(0) NOT NULL,
    update_by      int8        NOT NULL,
    update_dt      timestamptz(0) NOT NULL,
    CONSTRAINT userpokemon_pk PRIMARY KEY (userpokemon_id)
);
