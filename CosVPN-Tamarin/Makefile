# If you set these variables in your environment, the makefile will use those instead

# On linux, you can use xopen (or just any  command to launch a browser)
TAMARIN_BROWSER 	?= $(shell which xdg-open || which open)
# The local tamarin path
TAMARIN_BIN 		?= tamarin-prover
# All additional arguments to tamarin
TAMARIN_ARGS 		?= "--heuristic=i"
# The remote path to tamarin and directory that the .spthy gets copied to
TAMARIN_REMOTEBIN 	?= tamarin-prover
# The target protocol to be proved
TAMARIN_PROTO		?= wireguard
# The remote host for running
TAMARIN_HOST		?= proof.wireguard.com

logfile = log

default: auto

.PHONY: check-ports check-file remote local

%.spthy: %.m4
	m4 $< > $@

# Find the smallest unused port > 3000 on $(host); thanks http://unix.stackexchange.com/a/55940
# There's a race condition (if someone binds that port in the middle of the make) but it's hardly important
# Choose a local port for Tamarin (FYI you need -n in lsof or it takes 20s for DNS resolution... yuck)
check-ports:
ifndef TAMARIN_HOST
	$(error Running remotely requires a TAMARIN_HOST to be defined in your environment)
endif
	$(eval remoteport := $(shell ssh -t $(TAMARIN_HOST) \
	"ss -tln | awk 'NR > 1{gsub(/.*:/,"\"\"",\$$4); print \$$4}' | sort -un |\
	awk -v n=3001 '\$$0 < n {next}; \$$0 == n {n++; next}; {exit}; END {print n}'"))
	$(eval localport := $(shell lsof -n -i -P | grep -i "listen" |\
	grep -o "[0-9]* (LISTEN)" | cut -d' ' -f1 | sort -un | awk -v n=3001\
	'$$0 < n {next}; $$0 == n {n++; next}; {exit}; END {print n}'))

remote: check-ports $(TAMARIN_PROTO).spthy
	$(eval remotedir := $(shell ssh $(TAMARIN_HOST) "mktemp -d"))
	scp $(TAMARIN_PROTO).spthy $(TAMARIN_HOST):$(remotedir) > /dev/null
	@echo ========= local port $(localport), remote port $(remoteport) =========
	(tail -f $(logfile) | grep -m1 "Application launched" | xargs echo "" >> $(logfile);\
	$(TAMARIN_BROWSER) http://localhost:$(localport)/thy/trace/1/overview/help) &
	ssh -L $(localport):localhost:$(remoteport) -S'none' -t $(TAMARIN_HOST) \
		"$(TAMARIN_REMOTEBIN) interactive $(remotedir)/ --interface='*4' --port=$(remoteport) $(TAMARIN_ARGS)" 2>&1 | tee $(logfile)

local: $(TAMARIN_PROTO).spthy
	(sleep 1 && $(TAMARIN_BROWSER) http://localhost:3001/thy/trace/1/overview/help) &
	$(TAMARIN_BIN) interactive . --interface='*4' --port=3001 $(TAMARIN_ARGS)

auto: $(TAMARIN_PROTO).spthy
	$(TAMARIN_BIN) --prove $(TAMARIN_ARGS) $<
