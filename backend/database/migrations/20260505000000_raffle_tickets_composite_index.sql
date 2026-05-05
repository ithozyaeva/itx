-- Композитный индекс ускоряет два горячих запроса:
--   GetMemberTicketCount(raffle_id, member_id) — вызывается при каждом
--   check-in и в каждом просмотре карточки daily-раффла.
--   GetAll() с подзапросом my_tickets по текущему юзеру.
-- Существующих idx_raffle_tickets_raffle / idx_raffle_tickets_member
-- недостаточно при больших объёмах (раньше один пользователь мог иметь
-- миллионы записей в одной раффле).
CREATE INDEX IF NOT EXISTS idx_raffle_tickets_raffle_member
    ON raffle_tickets (raffle_id, member_id);
