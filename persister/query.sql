-- name: LoadCodebase :one
SELECT * FROM codebase WHERE repo_url = ? LIMIT 1;
-- name: ListCodebases :many
SELECT * FROM codebase ORDER BY repo_url;
-- name: InsertCodebase :one
INSERT INTO codebase ( repo_url ) VALUES ( ? ) RETURNING *;
-- name: DeleteCodebase :exec
DELETE FROM codebase WHERE repo_url = ?;
-- name: UpdateCodebase :exec
UPDATE codebase SET repo_url = ? WHERE repo_url = ? RETURNING *;
-- name: UpsertCodebase :one
INSERT INTO codebase (repo_url) VALUES (?)
ON CONFLICT (repo_url) DO UPDATE SET repo_url=excluded.repo_url RETURNING *;

-- name: LoadSurvey :one
SELECT * FROM survey WHERE id = ? LIMIT 1;
-- name: ListSurveys :many
SELECT * FROM survey ORDER BY timestamp;
-- name: ListCodebaseSurveys :many
SELECT * FROM survey WHERE codebase_id = ? ORDER BY timestamp;
-- name: InsertSurvey :one
INSERT INTO survey ( codebase_id,local_dir ) VALUES ( ?,? ) RETURNING *;
-- name: DeleteSurvey :exec
DELETE FROM survey WHERE id = ?;
-- name: DeleteCodebaseSurveys :exec
DELETE FROM survey WHERE codebase_id = ?;


-- name: LoadFile :one
SELECT * FROM file WHERE id = ? LIMIT 1;
-- name: ListFiles :many
SELECT * FROM file ORDER BY filepath;
-- name: InsertFile :one
INSERT INTO file ( survey_id,filepath ) VALUES ( ?,? ) RETURNING *;
-- name: DeleteFile :exec
DELETE FROM file WHERE id = ?;
-- name: UpdateFile :exec
UPDATE file SET filepath = ? WHERE id = ? RETURNING *;
-- name: UpsertFile :one
INSERT INTO file (survey_id,filepath) VALUES ( ?,? )
ON CONFLICT (survey_id,filepath) DO UPDATE SET filepath=excluded.filepath RETURNING *;


-- name: LoadSymbolType :one
SELECT * FROM symbol_type WHERE id = ? LIMIT 1;
-- name: ListSymbolTypes :many
SELECT * FROM symbol_type ORDER BY id;
-- name: ListSymbolTypesByName :many
SELECT * FROM symbol_type ORDER BY name;
-- name: InsertSymbolType :one
INSERT INTO symbol_type ( id,name ) VALUES ( ?,? ) RETURNING *;
-- name: DeleteSymbolType :exec
DELETE FROM symbol_type WHERE id = ?;
-- name: UpdateSymbolType :exec
UPDATE symbol_type SET name = ? WHERE id = ? RETURNING *;
-- name: UpsertSymbolType :one
INSERT INTO symbol_type ( id,name ) VALUES ( ?,? )
    ON CONFLICT (id) DO UPDATE SET name=excluded.name RETURNING *;

-- name: LoadType :one
SELECT * FROM type WHERE id = ? LIMIT 1;
-- name: ListTypes :many
SELECT * FROM type ORDER BY name;
-- name: InsertType :one
INSERT INTO type ( file_id, symbol_type_id, name,definition ) VALUES ( ?,?,?,? ) RETURNING *;
-- name: DeleteType :exec
DELETE FROM type WHERE id = ?;
-- name: UpdateType :exec
UPDATE type SET file_id= ?, symbol_type_id= ?, name= ?,definition= ? WHERE id = ? RETURNING *;

-- name: LoadImport :one
SELECT * FROM import WHERE id = ? LIMIT 1;
-- name: ListImports :many
SELECT * FROM import ORDER BY file_id;
-- name: InsertImport :one
INSERT INTO import ( file_id, package, alias ) VALUES ( ?,?,? ) RETURNING *;
-- name: DeleteImport :exec
DELETE FROM import WHERE id = ?;
-- name: UpdateImport :exec
UPDATE import SET file_id = ?, package = ?, alias = ? WHERE id = ? RETURNING *;
-- name: UpsertImport :one
INSERT INTO import ( file_id, package, alias ) VALUES ( ?,?,? ) 
ON CONFLICT (file_id,package,alias) DO UPDATE SET alias=excluded.alias RETURNING *;


-- name: LoadVariable :one
SELECT * FROM variable WHERE id = ? LIMIT 1;
-- name: ListVariables :many
SELECT * FROM variable ORDER BY name;
-- name: InsertVariable :one
INSERT INTO variable ( name,type_id,usage ) VALUES ( ?,?,? ) RETURNING *;
-- name: DeleteVariable :exec
DELETE FROM variable WHERE id = ?;
-- name: UpdateVariable :exec
UPDATE variable SET name = ? WHERE id = ? RETURNING *;


-- name: LoadMethod :one
SELECT * FROM method WHERE id = ? LIMIT 1;
-- name: ListMethods :many
SELECT * FROM method ORDER BY name;
-- name: InsertMethod :one
INSERT INTO method ( name,params,results,type_id ) VALUES ( ?,?,?,? ) RETURNING *;
-- name: DeleteMethod :exec
DELETE FROM method WHERE id = ?;
-- name: UpdateMethod :exec
UPDATE method SET name = ? WHERE id = ? RETURNING *;


-- name: LoadCategory :one
SELECT * FROM category WHERE id = ? LIMIT 1;
-- name: ListCategories :many
SELECT * FROM category ORDER BY name;
-- name: InsertCategory :one
INSERT INTO category ( name ) VALUES ( ? ) RETURNING *;
-- name: DeleteCategory :exec
DELETE FROM category WHERE id = ?;
-- name: UpdateCategory :exec
UPDATE category SET name = ? WHERE id = ? RETURNING *;
-- name: UpsertCategory :one
INSERT INTO category ( id,name ) VALUES ( ?,? )
ON CONFLICT (id) DO UPDATE SET name=excluded.name RETURNING *;








