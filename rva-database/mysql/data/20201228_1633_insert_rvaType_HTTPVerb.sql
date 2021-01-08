insert into rvaType
(idRvaEntity,typeName,description,creatorAccount,updaterAccount)
values
((select idRvaEntity from rvaEntity where idRvaModule='RVA' and entityName='HTTPVerb'),'GET','HTTPVerb GET','system','system'),
((select idRvaEntity from rvaEntity where idRvaModule='RVA' and entityName='HTTPVerb'),'HEAD','HTTPVerb HEAD','system','system'),
((select idRvaEntity from rvaEntity where idRvaModule='RVA' and entityName='HTTPVerb'),'POST','HTTPVerb POST','system','system'),
((select idRvaEntity from rvaEntity where idRvaModule='RVA' and entityName='HTTPVerb'),'PUT','HTTPVerb PUT','system','system'),
((select idRvaEntity from rvaEntity where idRvaModule='RVA' and entityName='HTTPVerb'),'DELETE','HTTPVerb DELETE','system','system'),
((select idRvaEntity from rvaEntity where idRvaModule='RVA' and entityName='HTTPVerb'),'CONNECT','HTTPVerb CONNECT','system','system'),
((select idRvaEntity from rvaEntity where idRvaModule='RVA' and entityName='HTTPVerb'),'OPTIONS','HTTPVerb OPTIONS','system','system'),
((select idRvaEntity from rvaEntity where idRvaModule='RVA' and entityName='HTTPVerb'),'TRACE','HTTPVerb TRACE','system','system'),
((select idRvaEntity from rvaEntity where idRvaModule='RVA' and entityName='HTTPVerb'),'PATCH','HTTPVerb PATCH','system','system');