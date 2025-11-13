create table public.tutors(
    id int generated always as identity primary key,
    full_name varchar(50) not null,
    email varchar(50) not null,
    constraint email_unique unique (email)
);

create table public.questions(
    id int generated always as identity primary key,
    question_text text not null,
    tutor_id int references tutors (id) on delete set null,
    created_at timestamp default now(),
    is_edit bool default false,
    constraint question_text_unique unique(question_text)
);

create table public.answers(
    id int generated always as identity primary key,
    answer_text text not null,
    tutor_id int references tutors (id) on delete set null,
    question_id int references questions (id) on delete cascade,
    created_at timestamp default now(),
    is_edit bool default false,
    constraint answer_text_unique unique(answer_text),
    constraint question_id_unique unique (question_id)
);

create table public.tags(
    id int generated always as identity primary key,
    tutor_id int references tutors(id) on delete set null,
    tag varchar(25) not null,
    constraint tag_unique unique (tag)
);

create table public.questions_tags(
    question_id int references questions(id) on delete cascade,
    tag_id int references tags(id) on delete cascade,
    constraint questions_tags_pk primary key (question_id, tag_id)
);

create table public.question_versions(
    id int generated always as identity primary key,
    question_id int references questions(id) on delete cascade,
    question_text text not null,
    tutor_id int references tutors(id) on delete set null,
    created_at timestamp default now(),
    version_number int not null,
    constraint unique_question_version unique (question_id, version_number)
);

create table public.answer_versions(
    id int generated always as identity primary key,
    answer_id int references answers(id) on delete cascade,
    answer_text text not null,
    tutor_id int references tutors(id) on delete set null,
    created_at timestamp default now(),
    version_number int not null,
    constraint unique_answer_version unique (answer_id, version_number)
);