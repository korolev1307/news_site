--user
    CREATE TABLE users (
        id integer primary key autoincrement,
        name varchar(100),
        surname varchar(100),
        patronumic varchar(100),
        login varchar(100),
        password varchar(1000),
        administrator boolean default 0,
        moderator boolean default 0
    );
    INSERT INTO "users" VALUES(1,'sanya','korolev','dmitrievich','korolev1307','password',1,0);

    --news

    CREATE TABLE news (
        id integer primary key autoincrement,
        title varchar(200),
        user_id references users(id), 
        content text,
        created_date timestamp,
        moderated_at timestamp default NULL,
        folder_name varchar(100),
        images boolean,
        files boolean,
        approved_by_administrator boolean default 0,    --default 0
        approved_by_moderator boolean default 0,        --default 0
        publishing_at_main_page boolean default 0,      --default 0
        publishing_at_lit_page boolean default 0,       --default 0
        publishing_at_EC boolean default 0,             --default 0
        moderated_by_id references users(id) default NULL   --default NULL
    );


    INSERT INTO "news" (id, title, user_id, content, created_date, folder_name, images, files) VALUES (1,'Publish on github',1,'Publish the source of tasks and picsort on github','2018-11-12 15:30:59','files/1',0,0);
