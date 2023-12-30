INSERT  INTO app (id,name,secret)
VALUES (1,'admin','test')
ON CONFLICT DO NOTHING;