drop table zenos.trans_test
create table zenos.trans_test(id bigint not null,num bigint not null,primary key(id))
insert into zenos.trans_test(id,num) values (1,1),(2,2),(3,3),(4,4)