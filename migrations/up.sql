create table if not EXISTS users(
    id serial primary key,
    name varchar (50) not null
);

create table if not exists rides(
    ride_id varchar (50) not null,
    scooter_id varchar (50) not null,
    user_id varchar (50) not null,
    status varchar (50),
    start_time int,
    stop_time int
);
