#!/usr/bin/env bash
# Maps Conorganizer's live-update surface from the source tree:
#   1. per live bucket, every Stream subscriber and every Broadcast call site
#   2. every POST/PUT/DELETE/PATCH route that broadcasts only through a
#      function it calls, at any depth, and every one that never broadcasts
#   3. handlers and functions that call http.Error or ReadSignals after datastar.NewSSE;
#      a NewSSE in a block that returns counts only inside that block
#   4. Datastar attributes the 1.0.4 client ignores or misreads
#
# It scans *.go and *.templ, skipping generated *_templ.go, *_test.go,
# .claude/, node_modules/, tmp/ and backup-* folders. Handler and function
# bodies are found by brace counting with string literals and // comments
# removed. Calls are matched by function name, preferring a definition in the
# route's own folder. Calls through a variable holding a function cannot be
# followed. Treat the last section as a list to check, not as a verdict.
#
# Usage, from anywhere:
#   bash .agents/skills/conorganizer-datastar-live/scripts/live-map.sh [repo-root]
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root="${1:-$(git -C "$script_dir" rev-parse --show-toplevel 2>/dev/null || true)}"
if [ -z "$root" ]; then
	echo "Pass the repo root as the first argument; the script is not inside a git checkout." >&2
	exit 1
fi
cd "$root"

if [ ! -f service/live/live.go ]; then
	echo "service/live/live.go not found under $root" >&2
	exit 1
fi

source_files() {
	find . \( -path ./.claude -o -path ./node_modules -o -path ./tmp -o -name 'backup-*' -o -path ./.git \) -prune -o \
		-type f \( -name '*.go' -o -name '*.templ' \) ! -name '*_templ.go' ! -name '*_test.go' -print \
		| sed 's|^\./||' | sort
}

# Shared awk helper: the code of one line with string, rune and raw-string
# literals emptied and // comments dropped. inraw carries a raw string across
# lines; callers reset it at each top-level declaration.
awk_code_of='
function code_of(line,    out, i, n, c, q) {
	out = ""; n = length(line); i = 1
	while (i <= n) {
		c = substr(line, i, 1)
		if (inraw) { if (c == "`") { inraw = 0; out = out "``" } i++; continue }
		if (c == "`") { inraw = 1; i++; continue }
		if (c == "\"" || c == "'\''") {
			q = c; i++
			while (i <= n) {
				c = substr(line, i, 1)
				if (c == "\\") { i += 2; continue }
				if (c == q) break
				i++
			}
			i++; out = out q q; continue
		}
		if (c == "/" && substr(line, i + 1, 1) == "/") break
		out = out c; i++
	}
	return out
}
function count(s, ch,    n) { n = gsub(ch, "", s); return n }
function after_sse(code,    rest, i, j, n, c, d, returns, late, found) {
	found = ""
	rest = code
	while ((i = index(rest, "NewSSE(")) > 0) {
		rest = substr(rest, i + 7)
		n = length(rest); d = 0; returns = 0
		for (j = 1; j <= n; j++) {
			c = substr(rest, j, 1)
			if (c == "{") d++
			else if (c == "}" && --d < 0) break
			else if (d == 0 && substr(rest, j, 6) == "return" && substr(rest, j - 1, 1) !~ /[A-Za-z0-9_]/ && substr(rest, j + 6, 1) !~ /[A-Za-z0-9_]/) returns = 1
		}
		late = returns ? substr(rest, 1, j - 1) : rest
		if (index(late, "http.Error(") > 0 && index(found, "http.Error") == 0) found = found (found == "" ? "" : ", ") "http.Error"
		if (index(late, "ReadSignals(") > 0 && index(found, "ReadSignals") == 0) found = found (found == "" ? "" : ", ") "ReadSignals"
		if (!returns) break
		rest = substr(rest, j)
	}
	return found
}
'

client_version="$(head -1 static/datastar.js 2>/dev/null | tr -d '\r' | sed 's|^// *||' || true)"
sdk_version="$(grep -E 'starfederation/datastar-go' go.mod 2>/dev/null | tr -d '\r' | awk '{print $2}' || true)"
echo "Live map for $root"
echo "Client: ${client_version:-unknown}   Go SDK: datastar-go ${sdk_version:-unknown}"
echo

buckets="$(grep -oE 'Bucket[A-Za-z]+ +Bucket = "[a-z]+"' service/live/live.go | tr -d '\r' | awk '{print $1 ":" $4}' | tr -d '"')"
if [ -z "$buckets" ]; then
	echo "No Bucket constants found in service/live/live.go" >&2
	exit 1
fi

# Subscribers: from a line that calls .Stream( until the Render: line, collect
# every live.Bucket* token. The route is the closest .Get("...", func before it,
# joined with the .Route("...", func prefixes open in the same function.
# Prefixes added by the calling function are not shown.
subscribers="$(source_files | tr '\n' '\0' | xargs -0 awk "$awk_code_of"'
FNR == 1 { depth = 0; sp = 0; inraw = 0; getpath = ""; collecting = 0 }
{
	sub(/\r$/, "")
	if ($0 ~ /^(func|templ|type|var|const|import|package)[ (]/) { depth = 0; sp = 0; inraw = 0 }
	if (match($0, /\.Route\("[^"]*", *func\(/)) {
		p = substr($0, RSTART + 8); sub(/".*/, "", p)
		sp++; spath[sp] = p; sdepth[sp] = depth
	}
	if (match($0, /\.Get\("[^"]*", *func\(/)) {
		p = substr($0, RSTART + 6); sub(/".*/, "", p)
		getpath = ""
		for (k = 1; k <= sp; k++) getpath = getpath spath[k]
		getpath = getpath p
		gsub(/\/\/+/, "/", getpath)
	}
	code = code_of($0)
	depth += count(code, "[{]") - count(code, "[}]")
	while (sp > 0 && depth <= sdepth[sp]) sp--
}
!collecting && /\.Stream\(/ { collecting = 1; start = FNR; buf = ""; route = getpath }
collecting {
	buf = buf " " $0
	if ($0 ~ /Render:/ || FNR - start > 12) {
		names = ""
		s = buf
		while (match(s, /live\.Bucket[A-Za-z]+/)) {
			names = names " " substr(s, RSTART + 5, RLENGTH - 5)
			s = substr(s, RSTART + RLENGTH)
		}
		if (names == "") names = " (buckets by variable)"
		printf "%s:%d\t%s\t%s\n", FILENAME, start, substr(names, 2), route
		collecting = 0
	}
}')"

# Broadcast call sites: any Broadcast( or broadcast( call with arguments,
# excluding declarations. A call is read until its own closing parenthesis.
broadcasts="$(source_files | tr '\n' '\0' | xargs -0 awk '
{ sub(/\r$/, "") }
!collecting && /[Bb]roadcast\([^)]/ && !/^func [^{]*[Bb]roadcast\(/ { collecting = 1; start = FNR; buf = "" }
collecting {
	buf = buf " " $0
	tail = buf; sub(/.*[Bb]roadcast\(/, "", tail)
	if (gsub(/\)/, ")", tail) > gsub(/\(/, "(", tail) || FNR - start > 6) {
		names = ""
		s = buf
		while (match(s, /live\.Bucket[A-Za-z]+/)) {
			names = names " " substr(s, RSTART + 5, RLENGTH - 5)
			s = substr(s, RSTART + RLENGTH)
		}
		if (names == "") names = " (buckets by variable)"
		printf "%s:%d\t%s\n", FILENAME, start, substr(names, 2)
		collecting = 0
	}
}')"

while IFS=: read -r const name; do
	echo "== Bucket $name ($const)"
	echo "   Stream subscribers:"
	printf '%s\n' "$subscribers" | awk -F '\t' -v c="$const" '$2 ~ ("(^| )" c "( |$)") { printf "     %-58s local route %-36s [%s]\n", $1, $3, $2 }'
	echo "   Broadcast call sites:"
	printf '%s\n' "$broadcasts" | awk -F '\t' -v c="$const" '$2 ~ ("(^| )" c "( |$)") { printf "     %-70s [%s]\n", $1, $2 }'
	echo
done <<EOF
$buckets
EOF

echo "== Calls whose buckets come from a variable (check by hand)"
printf '%s\n' "$subscribers" "$broadcasts" | awk -F '\t' '$2 ~ /variable/ { n++; printf "     %s\n", $1 } END { if (!n) print "     (none)" }'
echo

# Mutation routes and functions in one pass. A function body or a closure
# handler is followed by brace depth; a named handler ends where its call's
# parentheses close. A function broadcasts when its body calls Broadcast or
# calls a function that does, at any depth.
source_files | tr '\n' '\0' | xargs -0 awk "$awk_code_of"'
function dir_of(f,    d) { d = f; if (sub(/\/[^\/]*$/, "", d) == 0) d = "."; return d }
function ids_of(code,    s, id, out) {
	out = ""; s = code
	while (match(s, /[A-Za-z_][A-Za-z0-9_]*\(/)) {
		id = substr(s, RSTART, RLENGTH - 1)
		s = substr(s, RSTART + RLENGTH)
		if (id ~ /^(if|for|func|switch|return|append|len|make|new|panic|defer|go|range|map|case|select|recover|delete|close|cap|copy|print|println|string|int|int64|float64|byte|error|bool)$/) continue
		if (index(" " out " ", " " id " ") == 0) out = out " " id
	}
	return out
}
function resolve(name, dir,    keys, n, i) {
	if ((dir SUBSEP name) in fdirect) return fbc[dir SUBSEP name] ? dir SUBSEP name : ""
	if (!(name in fkeys)) return ""
	n = split(fkeys[name], keys, "\n")
	for (i = 2; i <= n; i++) if (fbc[keys[i]]) return keys[i]
	return ""
}
FNR == 1 { infunc = 0; inroute = 0; inraw = 0 }
{ sub(/\r$/, "") }
!infunc && /^func / {
	infunc = 1; fdepth = 0; fseen = 0; fcode = ""; fstart = FNR
	head = $0; sub(/^func (\([^)]*\) )?/, "", head); fname = head; sub(/[^A-Za-z0-9_].*/, "", fname)
	fkey = dir_of(FILENAME) SUBSEP fname
	if (!(fkey in fdirect)) { fdirect[fkey] = 0; fcalls[fkey] = ""; fkeys[fname] = fkeys[fname] "\n" fkey; fdir[fkey] = dir_of(FILENAME) }
	inraw = 0
}
!inroute && /\.(Post|Put|Delete|Patch)\("/ {
	inroute = 1; rstart = FNR; rdepth = 0; rparen = 0; rseen = 0; rcode = ""
	match($0, /\.(Post|Put|Delete|Patch)\("/)
	rmethod = substr($0, RSTART + 1, RLENGTH - 3)
	rest = substr($0, RSTART + RLENGTH)
	rpath = rest; sub(/".*/, "", rpath)
	rclosure = (rest ~ /func\(/)
	saved = inraw; rfirst = code_of(substr($0, RSTART + 1)); inraw = saved
	routefirst = 1
}
{
	code = code_of($0)
	if (infunc) {
		line = code
		if (fcode == "") sub(/^func (\([^)]*\) )?[A-Za-z_][A-Za-z0-9_]*/, "", line)
		fcode = fcode " " line
		opens = count(code, "[{]"); fdepth += opens - count(code, "[}]")
		if (opens > 0) fseen = 1
		if (fseen && fdepth <= 0) {
			if (fcode ~ /[Bb]roadcast\([^)]/) fdirect[fkey] = 1
			fcalls[fkey] = fcalls[fkey] ids_of(fcode)
			late = after_sse(fcode)
			if (late != "" && fcode !~ /\.(Get|Post|Put|Delete|Patch|Route)\(/) { nsse++; ssefile[nsse] = FILENAME; sseline[nsse] = fstart; ssewhat[nsse] = "func " fname "  [" late "]" }
			infunc = 0
		}
	}
	if (inroute) {
		rc = routefirst ? rfirst : code
		routefirst = 0
		rcode = rcode " " rc
		opens = count(rc, "[{]"); rdepth += opens - count(rc, "[}]")
		rparen += count(rc, "[(]") - count(rc, "[)]")
		if (opens > 0) rseen = 1
		if ((rclosure && rseen && rdepth <= 0) || (!rclosure && rparen <= 0)) {
			nr++
			rfile[nr] = FILENAME; rline[nr] = rstart; rmeth[nr] = rmethod; rpth[nr] = rpath
			rdirect[nr] = (rcode ~ /[Bb]roadcast\([^)]/)
			rcalls[nr] = ids_of(rcode)
			late = rclosure ? after_sse(rcode) : ""
			if (late != "") { nsse++; ssefile[nsse] = FILENAME; sseline[nsse] = rstart; ssewhat[nsse] = rmethod " \"" rpath "\"  [" late "]" }
			inroute = 0
		}
	}
}
END {
	for (k in fdirect) fbc[k] = fdirect[k]
	changed = 1
	while (changed) {
		changed = 0
		for (k in fdirect) {
			if (fbc[k]) continue
			n = split(fcalls[k], ids, " ")
			for (i = 1; i <= n; i++) {
				hit = resolve(ids[i], fdir[k])
				if (hit != "" && hit != k) { fbc[k] = 1; fvia[k] = hit; changed = 1; break }
			}
		}
	}
	print "== Mutation routes that broadcast through a function they call"
	for (r = 1; r <= nr; r++) {
		if (rdirect[r]) continue
		n = split(rcalls[r], ids, " ")
		hit = ""
		for (i = 1; i <= n && hit == ""; i++) hit = resolve(ids[i], dir_of(rfile[r]))
		if (hit == "") { never[r] = 1; continue }
		chain = ""
		for (guard = 0; hit != "" && guard < 10; guard++) {
			split(hit, part, SUBSEP)
			chain = chain (chain == "" ? "" : " > ") part[2]
			hit = (hit in fvia) ? fvia[hit] : ""
		}
		printf "     %-72s %-6s %-50s via %s\n", rfile[r] ":" rline[r], rmeth[r], "\"" rpth[r] "\"", chain
	}
	print ""
	print "== POST/PUT/DELETE/PATCH handlers that never call Broadcast"
	for (r = 1; r <= nr; r++) if (never[r]) printf "     %-72s %-6s %s\n", rfile[r] ":" rline[r], rmeth[r], "\"" rpth[r] "\""
	print ""
	print "     A handler in this list is fine when it changes nothing a live page renders: previews, data no page streams, auth."
	print "     The event image upload does not broadcast, although a first card or banner image changes the image URL on event cards."
	print "     Everything else is a missed broadcast, or a broadcast through a function value: open the handler and check."
	print ""
	print "== Handlers and functions with http.Error or ReadSignals after datastar.NewSSE (do not copy; check by hand)"
	if (!nsse) print "     (none)"
	for (i = 1; i <= nsse; i++) printf "     %-72s %s\n", ssefile[i] ":" sseline[i], ssewhat[i]
}'
echo

# Datastar 1.0.4 reads data-<plugin>:<key>__<modifier>. A hyphen after a plugin
# name makes the client look up a plugin that does not exist, so the attribute
# is ignored. data-on-intersect, data-on-interval and data-on-signal-patch are
# real plugins. A modifier before the key, as in data-signals__ifmissing:name,
# makes the client drop the key. HTML lowercases attribute names, so a
# camelCase key declares a lowercase signal. A keyless data-signals with a
# string value is spread into signals named 0, 1 and so on.
echo "== Datastar attributes the 1.0.4 client ignores or misreads"
hyphen_pattern='data-(on|bind|attr|class|signals|indicator|ref|style|computed|show|text|effect|init)-[a-z]'
order_pattern='data-[a-z-]+__[a-z.0-9]+:[A-Za-z_-]'
camel_pattern='data-(signals|bind|computed|ref|indicator|class|attr|style|on):[a-z0-9_-]*[A-Z]'
string_pattern='data-signals=\{ *fmt\.Sprintf\(`"'
datastar_hits="$(source_files | grep -E '\.templ$' | tr '\n' '\0' | xargs -0 grep -nE "$hyphen_pattern|$order_pattern|$camel_pattern|$string_pattern" 2>/dev/null \
	| tr -d '\r' | grep -vE 'data-on-(intersect|interval|signal-patch)' || true)"
if [ -n "$datastar_hits" ]; then
	printf '%s\n' "$datastar_hits" | while IFS= read -r hit; do
		location="$(printf '%s' "$hit" | cut -d: -f1-2)"
		attribute="$(printf '%s' "$hit" | grep -oE "$hyphen_pattern[A-Za-z0-9_.-]*|$order_pattern[A-Za-z0-9_-]*|$camel_pattern[A-Za-z0-9_-]*|$string_pattern" | head -1)"
		case "$attribute" in
			*__*:*) reason="modifier before the key" ;;
			*Sprintf*) attribute='data-signals={ fmt.Sprintf(`"...`'; reason="keyless string value" ;;
			*:*[A-Z]*) reason="camelCase key (HTML lowercases it)" ;;
			*) reason="hyphen form" ;;
		esac
		printf '     %-72s %-50s %s\n' "$location" "$attribute" "$reason"
	done
else
	echo "     (none)"
fi
