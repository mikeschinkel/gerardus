.PHONY: .force
.force:

export GOPROXY = https://proxy.golang.org/
export GOSUMDB = sum.golang.org


define for_each_module
	@root="$$(pwd)" && for file in $$(find . | grep go.mod | grep -v .archive); \
		do true \
  		&& dir="$$(dirname "$${file}")" \
  		&& cd "$${root}/$${dir}" \
  		&& echo "$(1) $${dir}" \
  		&& $(2); \
  		done
endef

test: .force
	$(call \
		for_each_module,\
		Testing gerardus,\
		GERARDUS_SOURCE_DIR=/Users/mikeschinkel/Projects/gerardus \
		go test -tags test -timeout=0 ./...)


vulncheck: .force
	$(call for_each_module,Vulnerability Checking gerardus,govulncheck ./... > /dev/null)

tidy: .force
	$(call for_each_module,Tidying gerardus,go mod tidy)

why: .force
	$(call for_each_module,Why for $(package) in gerardus,go mod why -m $(package))

graph: .force
	$(call for_each_module,Graph for $(package) in gerardus,go mod graph | grep $(package))

test-release: .force
	@echo "Running GoReleaser..."
	@GOLANG_VERSION="$$(go version | cut -w -f 3)" \
		COMMITTER_NAME="$$(git log  --max-count 1 --pretty=format:'%cn')" \
		goreleaser release --snapshot --skip-publish --clean


test-run: .force
	@rm -f ./test/testdata/bendavis.out.sql \
	  	&& time ./bin/gerardus run \
		--input ./test/testdata/bendavis.sql \
		--output ./test/testdata/bendavis.out.sql \
		--find bendavis.com \
		--replace dww.local \
		--loglevel none

release: .force  # TODO Need to add validation and ideally logic to determine the version
	@echo "Releasing $(version)..."
	git tag "$(version)"
	git push
	git push --tags

revoke: .force
	@echo "Revoking $(version)..."
	git tag -d "$(version)"
	git push origin --delete "$(version)"


# Phony target to force make to always evaluate targets
.PHONY: build gen

# Actual files that sqlc depends on
SQLC_DEPS=./persister/query.sql ./persister/schema.sql ./persister/sqlc.yaml
SQLC_OUT= ./persister/db.go ./persister/models.go ./persister/query.sql.go

# Initialize RUN_ONCE variable
RUN_ONCE=1

# Rule to run go generate when one of the dependencies changes
$(SQLC_OUT): $(SQLC_DEPS)
	@RUN_ONCE=$(intcmp $(RUN_ONCE),1,$(shell go generate -x ./... && echo 0))

gen: $(SQLC_OUT)

# Rule to build the application, depending on gen
build: .force
	@cd cmd && \
	echo "Building gerardus" && \
	go build -o ../bin/gerardus . && \
	echo "Done"

