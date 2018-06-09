CREATE OR REPLACE FUNCTION pgspy_notificator()
  RETURNS TRIGGER AS $trigger$
DECLARE
  rec     RECORD;
  channel TEXT;
BEGIN
  IF tg_nargs < 1
  THEN
    channel := 'pgspy';
  ELSE
    channel := tg_argv [0];
  END IF;

  CASE TG_OP
    WHEN 'INSERT'
    THEN
      PERFORM pg_notify(
          tg_argv [0],
          json_build_object(
              'table', TG_TABLE_NAME,
              'operation', tg_op,
              'time', extract(epoch from date_trunc('seconds', now())),
              'new', row_to_json(NEW)
          ) :: TEXT
      );
      rec := NEW;
    WHEN 'UPDATE'
    THEN
      PERFORM pg_notify(
          tg_argv [0],
          json_build_object(
              'table', TG_TABLE_NAME,
              'operation', tg_op,
              'time', extract(epoch from date_trunc('seconds', now())),
              'new', row_to_json(NEW),
              'old', row_to_json(OLD)
          ) :: TEXT
      );
      rec := NEW;
    WHEN 'DELETE'
    THEN
      PERFORM pg_notify(
          tg_argv [0],
          json_build_object(
              'table', TG_TABLE_NAME,
              'operation', tg_op,
              'time', extract(epoch from date_trunc('seconds', now())),
              'old', row_to_json(OLD)
          ) :: TEXT
      );
      rec := OLD;
  ELSE
    RAISE EXCEPTION 'Unknown TG_OP: "%"', TG_OP;
  END CASE;

  RETURN rec;
END;
$trigger$
LANGUAGE plpgsql;
