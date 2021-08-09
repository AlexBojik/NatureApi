package sql

var BaseLayerList = "select id, name, url, description, minZoom, maxZoom from base_layers order by id"
var GroupLayerList = "select id, name, icon from group_layers order by id"
var DictionaryList = "select id, name from dictionaries order by name"
var LayerList = "select group_id, id, name, type, url, color, commonName, commonDescription, symbol, cluster, line_width, line_color, warning, limitation, 'order' from layers where group_id = ? and (not limitation or ?) order by `order`"
var NewsList = "select id, descr, created, start_date, end_date from news where ?=? order by start_date desc"
var NewsFilteredList = "select id, descr, created, start_date, end_date FROM news where start_date < ? and IFNULL(end_date, '2100-01-01 01:00:00') > ? order by id"
var Layer = "select o.id, o.layerId, o.name, ST_AsGeoJSON(c.g) as geojson from objects o inner join coordinates c on o.layerId = ? and o.id = c.id"
var BaseLayersCreate = "insert into base_layers (name, url, description) values (?, ?, ?)"
var NewsCreate = "insert into news (created, start_date, end_date, descr) values (?, ?, ?, ?)"
var BaseLayerUpdate = "update base_layers SET name=?, url=?, description=? where id=?"
var NewsUpdate = "update news SET created=?, start_date=?, end_date=?, descr=? where id=?"
var BaseLayerDelete = "delete from base_layers where id=?"
var NewsDelete = "delete from news where id=?"
var FeaturesByFilter = "select o.id, o.layerId, o.name, ST_AsGeoJSON(c.g) as geojson from objects o inner join coordinates c on o.id = c.id inner join (select id from objects o where name like CONCAT('%', ?, '%') union select objectId from fields_values fv where value like CONCAT('%', ?, '%') union select objectId id from fields_values fv inner join dictionary_values dv on fv.value_num = dv.id where name like CONCAT('%',? ,'%') LIMIT 100) f on o.id = f.id"
var AdditionalFields = "select f.name, IFNULL(dv.name, fv.value) as value from fields_values fv INNER JOIN fields f ON f.id = fv.fieldId AND (not f.limitation OR ?) LEFT JOIN dictionary_values dv  ON f.`type` = dv.dictId AND fv.value_num = dv.id where objectId = ? order by fv.fieldId"
var FeaturesByIds = "select o.id, o.layerId, o.name, ST_AsGeoJSON(c.g) as geojson from objects o inner join coordinates c on o.id = c.id where o.id in (%s)"
var FeaturesByRegion = "select o.id, o.layerId, o.name, ST_AsGeoJSON(c.g) as geojson from objects o inner join coordinates c on o.id = c.id WHERE ST_Contains(ST_GeomFromText(?), g) LIMIT 100"
var FieldsList = "select id, name from dictionary_values"
var FeaturesByFields = "select o.id, o.layerId, o.name, ST_AsGeoJSON(c.g) as geojson from objects o inner join coordinates c on o.id = c.id inner join (select objectId id from fields_values where value_num in (%s)) fv on o.id = fv.id"
var CoordinatesUpdate = "update coordinates SET g=ST_GeomFromText(?, 4326) where id=?"
var ClusterCoordinate = "select o.id, o.layerId, o.name, ST_AsGeoJSON(ST_Centroid(g)) as geojson from coordinates c  inner join objects o on c.id  = o.id  WHERE layerId = ? and ST_GeometryType(g) IN ('POLYGON', 'MULTIPOLYGON')"
var User = "select name, token, phone, email, snils, regAddr, proAddr, doc, admin, layers, dicts, messages, info, groupId, block from users where token = ?"
var UserList = "select id, name, token, phone, email, snils, regAddr, proAddr, doc, admin, layers, dicts, messages, info, groupId, block from users where groupId = ? or ? = 0"
var UserGroupList = "select id, name, admin, layers, dicts, messages, info from user_groups"
var UserCreate = "insert into users (name, token, phone, email, snils, regAddr, proAddr, doc) values (?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE name=name, token=token, phone=phone, email=email, token=token, snils=snils, regAddr=regAddr, proAddr=proAddr, doc=doc"
var UserGroupCreate = "insert into user_groups (name, admin, layers, dicts, messages, info) values (?, ?, ?, ?, ?, ?)"
var UserGroupUpdate = "update user_groups SET name=?, admin=?, layers=?, dicts=?, messages=?, info=? where id=?"
var UserUpdate = "update users SET admin=?, layers=?, dicts=?, messages=?, info=?, groupId=?, block=? where token=?"
var Message = "select data FROM images where message_id = ?"
var MessageCreate = "insert into user_messages (token, text, point, time) values(?, ?, ST_GeomFromText(?, 4326), ?)"
var MessageCount = "select count(id) from user_messages WHERE status = 0"
var MessageList = "select um.id, um.text, u.name, um.status, um.time, um.end, um.comment, um.employerId, e.name as employerName from user_messages um inner join users u ON um.token = u.token left join users e ON um.employerId = e.token"
var MessageNotSendList = "select um.id, um.text, u.name, um.status, ST_AsGeoJSON(um.point) as point, u.token from user_messages um inner join users u ON um.token = u.token where um.status = 0"
var ImageCreate = "insert into images (message_id, data) values(?, ?)"
var MessageUpdateStatus = "update user_messages set status = ? where id = ?"
var LayerCreate = "insert into layers (name, type, group_id, url, color, commonName, commonDescription, warning, symbol, cluster, `order`, line_width, line_color, limitation) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
var LayerUpdate = "update layers SET name=?, type=?, group_id=?, url=?, color=?, commonName=?, commonDescription=?, warning=?, symbol=?, cluster=?, `order`=?, line_width=?, line_color=?, limitation=? where id=?"
var ObjectCreate = "insert into objects (layerId, name) values (?, ?)"
var FieldsValuesCreate = "insert into fields_values (objectId, fieldId, value, value_num) values (?, ?, ?, ?)"
