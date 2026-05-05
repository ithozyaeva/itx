-- Источники билетов: вместо неконтролируемой покупки/чек-ина —
-- идемпотентная выдача билета за конкретное достижение.
-- source_type описывает тип активности (check_in, daily_task, all_dailies_bonus,
-- challenge, attend_event, purchase, legacy), source_id — конкретный артефакт
-- (id задачи/события/раффла/билета). UNIQUE (raffle_id, member_id, source_type,
-- source_id) гарантирует, что повторный вызов хука не плодит билет.

ALTER TABLE raffle_tickets ADD COLUMN IF NOT EXISTS source_type TEXT;
ALTER TABLE raffle_tickets ADD COLUMN IF NOT EXISTS source_id BIGINT;

-- Backfill: существующие билеты получают source_type='legacy' и source_id=id
-- (id уникален для всех существующих строк, поэтому UNIQUE не схлопнется).
UPDATE raffle_tickets
   SET source_type = 'legacy', source_id = id
 WHERE source_type IS NULL OR source_id IS NULL;

ALTER TABLE raffle_tickets ALTER COLUMN source_type SET NOT NULL;
ALTER TABLE raffle_tickets ALTER COLUMN source_id   SET NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS uniq_raffle_ticket_source
    ON raffle_tickets (raffle_id, member_id, source_type, source_id);

-- Для purchase у одной покупки count билетов — чтобы UNIQUE не схлопнулся,
-- каждому билету в покупке нужен свой source_id. Берём из выделенного
-- sequence: один SQL-запрос вставляет любое количество строк.
CREATE SEQUENCE IF NOT EXISTS raffle_ticket_purchase_seq;
