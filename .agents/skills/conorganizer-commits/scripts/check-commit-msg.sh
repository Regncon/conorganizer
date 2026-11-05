#!/usr/bin/env bash
# Checks a commit message against the conorganizer rules:
#   subject: capitalized imperative, no trailing period, at most 72 characters,
#            no "NNN | " issue prefix
#   body:    one blank line after the subject, no Co-Authored-By or
#            Generated-with lines, lines wrapped at 72 (a note, not a failure)
# Reads the message from the file named as the first argument, or from stdin.
# Lines starting with # are ignored, as git ignores them in the editor.
# Merge, Revert, fixup! and squash! messages are skipped.
#
# Usage, from the repo root (quote the path, the checkout has a space in it):
#   bash ".agents/skills/conorganizer-commits/scripts/check-commit-msg.sh" message.txt
#   git log -1 --format=%B | bash ".agents/skills/conorganizer-commits/scripts/check-commit-msg.sh"
# Prints one line per problem and exits 1; prints nothing and exits 0 when the message passes.
set -eu

if [ "${1:-}" = "--help" ] || [ "${1:-}" = "-h" ]; then
	sed -n '2,14p' "$0" | sed 's/^# \{0,1\}//'
	exit 0
fi

if locale -a 2>/dev/null | grep -qi '^C\.utf-\?8$'; then
	export LC_ALL=C.UTF-8
fi

if [ $# -ge 1 ]; then
	input=$(cat -- "$1")
else
	input=$(cat)
fi

message=$(printf '%s\n' "$input" | tr -d '\r' | awk '
	/^#/ { next }
	{ lines[++n] = $0 }
	END {
		s = 1; while (s <= n && lines[s] == "") s++
		e = n; while (e >= s && lines[e] == "") e--
		for (i = s; i <= e; i++) print lines[i]
	}')

problems=0
fail() {
	printf 'commit-msg: %s\n' "$1" >&2
	problems=$((problems + 1))
}

if [ -z "$message" ]; then
	fail "the message is empty"
	exit 1
fi

subject=$(printf '%s\n' "$message" | sed -n '1p')
second=$(printf '%s\n' "$message" | sed -n '2p')
linecount=$(printf '%s\n' "$message" | awk 'END { print NR }')

case "$subject" in
	"Merge "*|"Revert \""*|"fixup! "*|"squash! "*) exit 0 ;;
esac

if printf '%s' "$subject" | grep -q '^[[:lower:]]'; then
	fail "the subject starts with a lowercase letter; capitalize the first word (\"Add ...\", \"Fix ...\")"
fi

if printf '%s' "$subject" | grep -Eq '^[0-9]+ *\| '; then
	fail "the subject starts with an issue number prefix; put the issue in the body as \"Closes #NNN\" instead"
fi

if printf '%s' "$subject" | sed -E 's/^[0-9]+ *\| *//' | grep -Eiq '^(added|adds|fixed|fixes|updated|updates|removed|removes|changed|changes|moved|moves|renamed|renames|improved|improves|refactored|refactors|implemented|implements)\b'; then
	fail "the subject is not imperative; write what the commit does (\"Add\", \"Fix\", \"Update\"), as if completing \"This commit will ...\""
fi

case "$subject" in
	*.) fail "the subject ends with a period; drop it" ;;
esac

length=$(printf '%s' "$subject" | wc -m | tr -d ' ')
if [ "$length" -gt 72 ]; then
	fail "the subject is $length characters; keep it at 72 or fewer, so git log and GitHub show it whole"
fi

if [ "$linecount" -gt 1 ] && [ -n "$second" ]; then
	fail "the second line is not blank; leave one blank line between the subject and the body"
fi

if printf '%s\n' "$message" | grep -Eiq '^[[:space:]]*(co-authored-by|claude-session):|generated[ -]with|🤖'; then
	fail "the message has an attribution line (Co-Authored-By, Claude-Session or Generated with); the repo keeps those out of its history"
fi

if [ "$linecount" -gt 2 ]; then
	printf '%s\n' "$message" | awk -v limit=72 'NR > 2 && length($0) > limit && $0 !~ /https?:\/\// { printf "commit-msg: note: body line %d is %d characters; wrap the body at 72\n", NR, length($0) }' >&2
fi

[ "$problems" -eq 0 ]
