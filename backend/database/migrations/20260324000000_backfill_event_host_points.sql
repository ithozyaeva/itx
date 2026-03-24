-- Backfill point_transactions for event hosts who missed points after the unique index was dropped
-- (bug window: ~2026-03-06 to 2026-03-24 when AwardPoints was silently failing)
INSERT INTO point_transactions (member_id, amount, reason, source_type, source_id, description)
SELECT
    eh.member_id,
    25,
    'event_host',
    'event',
    e.id,
    'Проведение события: ' || e.title
FROM event_hosts eh
JOIN events e ON e.id = eh.event_id
WHERE e.date < NOW()
  AND NOT EXISTS (
      SELECT 1 FROM point_transactions pt
      WHERE pt.member_id = eh.member_id
        AND pt.reason = 'event_host'
        AND pt.source_type = 'event'
        AND pt.source_id = e.id
  );
