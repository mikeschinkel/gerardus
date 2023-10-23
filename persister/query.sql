-- name: LoadProject :one
SELECT * FROM project WHERE id = ? LIMIT 1;
-- name: LoadProjectByName :one
SELECT * FROM project WHERE name = ? LIMIT 1;
-- name: LoadProjectByRepoURL :one
SELECT * FROM project WHERE repo_url = ? LIMIT 1;
-- name: LoadProjectRepoURL :one
SELECT repo_url FROM project WHERE id = ? LIMIT 1;
-- name: ListProjects :many
SELECT * FROM project ORDER BY name;
-- name: InsertProject :one
INSERT INTO project ( name,about,repo_url, website ) VALUES ( ?,?,?,? ) RETURNING *;
-- name: DeleteProject :exec
DELETE FROM project WHERE id = ?;
-- name: DeleteProjectByName :exec
DELETE FROM project WHERE name = ?;
-- name: UpdateProject :exec
UPDATE project SET name = ?, about = ?, repo_url = ?, website = ? WHERE id = ? RETURNING *;
-- name: UpdateProjectByName :exec
UPDATE project SET repo_url = ?, about = ?, website = ? WHERE name = ? RETURNING *;
-- name: UpsertProject :one
INSERT INTO project ( name,about,repo_url,website ) VALUES ( ?,?,?,? )
ON CONFLICT (name) DO UPDATE SET about=excluded.about, repo_url=excluded.repo_url, website=excluded.website  RETURNING *;


-- name: LoadCodebase :one
SELECT * FROM codebase WHERE id = ? LIMIT 1;
-- name: LoadCodebaseByProjectNameAndVersionTag :one
SELECT c.id FROM codebase c JOIN project p ON p.id=c.project_id WHERE p.name = ? AND c.version_tag = ? LIMIT 1;
-- name: LoadCodebaseIdByRepoURL :one
SELECT c.id FROM codebase c JOIN project p ON p.id=c.project_id WHERE p.repo_url = ? LIMIT 1;
-- name: ListCodebases :many
SELECT * FROM codebase ORDER BY project_id,version_tag;
-- name: InsertCodebase :one
INSERT INTO codebase ( project_id,version_tag,source_url ) VALUES ( ?,?,? ) RETURNING *;
-- name: DeleteCodebase :exec
DELETE FROM codebase WHERE id = ?;
-- name: DeleteCodebaseByProjectIdAndVersionTag :exec
DELETE FROM codebase WHERE project_id = ? AND version_tag = ?;
-- name: UpdateCodebase :exec
UPDATE codebase SET project_id = ?, version_tag = ?, source_url = ? WHERE id = ? RETURNING *;
-- name: UpdateCodebaseByProjectIdAndVersionTag :exec
UPDATE codebase SET source_url = ? WHERE project_id = ? AND version_tag = ? RETURNING *;
-- name: UpsertCodebase :one
INSERT INTO codebase ( project_id,version_tag,source_url ) VALUES ( ?,?,? )
ON CONFLICT (project_id,version_tag) DO UPDATE SET source_url=excluded.source_url,version_tag=excluded.version_tag RETURNING *;

-- name: LoadSurvey :one
SELECT * FROM survey WHERE id = ? LIMIT 1;
-- name: LoadSurveyByRepoURL :one
SELECT * FROM survey s JOIN codebase c ON c.id=s.codebase_id  JOIN project p ON p.id=c.project_id WHERE p.repo_url = ? LIMIT 1;
-- name: ListSurveys :many
SELECT
    sv.id,
    cb.source_url,
    sv.local_dir,
    sv.timestamp
FROM survey AS sv
    JOIN codebase cb ON cb.id=sv.codebase_id
    JOIN project cb ON cb.id=sv.codebase_id
ORDER BY timestamp;
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
-- name: ListFilesBySurvey :many
SELECT * FROM file WHERE survey_id= ? ORDER BY filepath;
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
SELECT * FROM type_view ORDER BY name;
-- name: ListTypesBySurvey :many
SELECT * FROM type_view WHERE survey_id = ?;
-- name: ListTypesByFile :many
SELECT * FROM type_view WHERE file_id = ?;
-- name: InsertType :one
INSERT INTO type ( file_id, survey_id,symbol_type_id, name,definition ) VALUES ( ?,?,?,?,? ) RETURNING *;
-- name: DeleteType :exec
DELETE FROM type WHERE id = ?;
-- name: UpdateType :exec
UPDATE type SET file_id= ?, symbol_type_id= ?, name= ?,definition= ? WHERE id = ? RETURNING *;

-- name: LoadPackage :one
SELECT * FROM package WHERE id = ? LIMIT 1;
-- name: ListPackages :many
SELECT * FROM package ORDER BY path;
-- name: InsertPackage :one
INSERT INTO package ( path,source ) VALUES ( ?,? ) RETURNING *;
-- name: DeletePackage :exec
DELETE FROM package WHERE id = ?;
-- name: UpdatePackage :exec
UPDATE package SET path = ?, source = ? WHERE id = ? RETURNING *;
-- name: UpsertPackage :one
INSERT INTO package ( path, source ) VALUES ( ?,? )
ON CONFLICT (path, source) DO UPDATE SET path=excluded.path, source=excluded.source RETURNING *;


-- name: LoadImport :one
SELECT * FROM import WHERE id = ? LIMIT 1;
-- name: ListImports :many
SELECT * FROM import ORDER BY file_id;
-- name: InsertImport :one
INSERT INTO import ( file_id, survey_id, package_id, alias ) VALUES ( ?,?,?,? ) RETURNING *;
-- name: DeleteImport :exec
DELETE FROM import WHERE id = ?;
-- name: UpdateImport :exec
UPDATE import SET file_id = ?, survey_id = ?, package_id = ? WHERE id = ? RETURNING *;
-- name: UpsertImport :one
INSERT INTO import ( file_id, survey_id, package_id, alias ) VALUES ( ?,?,?,? )
ON CONFLICT (file_id, survey_id, package_id, alias) DO UPDATE SET alias=excluded.alias RETURNING *;


-- name: LoadVariable :one
SELECT * FROM variable WHERE id = ? LIMIT 1;
-- name: ListVariables :many
SELECT * FROM variable ORDER BY survey_id,type_id,is_result,name;
-- name: InsertVariable :one
INSERT INTO variable ( name,type_id,survey_id,usage ) VALUES ( ?,?,?,? ) RETURNING *;
-- name: DeleteVariable :exec
DELETE FROM variable WHERE id = ?;
-- name: UpdateVariable :exec
UPDATE variable SET name = ?,survey_id = ?,type_id = ?,usage=? WHERE id = ? RETURNING *;


-- name: LoadMethod :one
SELECT * FROM method WHERE id = ? LIMIT 1;
-- name: ListMethods :many
SELECT * FROM method ORDER BY survey_id,file_id,type_id,name;
-- name: InsertMethod :one
INSERT INTO method ( name,params,results,type_id,file_id,survey_id ) VALUES ( ?,?,?,?,?,? ) RETURNING *;
-- name: DeleteMethod :exec
DELETE FROM method WHERE id = ?;
-- name: UpdateMethod :exec
UPDATE method SET name = ?,params = ?,results = ?,type_id = ?,file_id = ?,survey_id = ? WHERE id = ? RETURNING *;


-- name: LoadCategory :one
SELECT * FROM category WHERE id = ? LIMIT 1;
-- name: ListCategories :many
SELECT * FROM category ORDER BY name;
-- name: InsertCategory :one
INSERT INTO category ( survey_id,name ) VALUES ( ?,? ) RETURNING *;
-- name: DeleteCategory :exec
DELETE FROM category WHERE id = ?;
-- name: UpdateCategory :exec
UPDATE category SET name = ? WHERE id = ? RETURNING *;
-- name: UpsertCategory :one
INSERT INTO category ( survey_id,name ) VALUES ( ?,? )
ON CONFLICT (survey_id,name) DO UPDATE SET name=excluded.name RETURNING *;





