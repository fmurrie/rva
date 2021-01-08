create procedure stepOne
(
in stepName varchar(100)
)
begin

select concat(stepName,'-stepOne') as stepName;
select concat(stepName,'-stepOneOne') as stepName;
select concat(stepName,'-stepOneOneOne') as stepName;

end