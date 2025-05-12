CREATE TABLE schedules (
	id bigserial NOT NULL,
	user_id int8 NULL,
	name_medication text NULL,
	medication_per_day int8 NULL,
	duration_medication int8 NULL,
	CONSTRAINT schedules_pkey PRIMARY KEY (id)
);