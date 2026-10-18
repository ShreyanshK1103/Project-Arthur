#!/usr/bin/env bash
# End-to-end smoke test for the Arthur backend.
# Needs: API + worker running, jq, curl.
#
#   REPO=https://github.com/<you>/<a-vite-or-cra-repo> ./smoke_test.sh
#
# Optional: BASE (default http://localhost:8080/v1), BRANCH (default main),
#           BUILD_CMD (default "npm run build"), OUT_DIR (default dist)

set -u
BASE="${BASE:-http://localhost:8080/v1}"
ROOT="${BASE%/v1}"
REPO="${REPO:?set REPO to a public https GitHub repo that builds with npm}"
BRANCH="${BRANCH:-main}"
BUILD_CMD="${BUILD_CMD:-npm run build}"
OUT_DIR="${OUT_DIR:-dist}"

TMP="$(mktemp -d)"
JAR1="$TMP/user1.jar"
JAR2="$TMP/user2.jar"
RUN="$(date +%s)"
PASS=0
FAIL=0

ok()   { echo "  PASS  $1"; PASS=$((PASS+1)); }
bad()  { echo "  FAIL  $1"; FAIL=$((FAIL+1)); }
expect() { # expect <description> <actual> <expected>
  if [ "$2" = "$3" ]; then ok "$1"; else bad "$1 (got '$2', want '$3')"; fi
}
code() { curl -s -o /dev/null -w '%{http_code}' "$@"; }
api()  { curl -s "$@"; }

echo "== health"
expect "GET /healthz" "$(code "$BASE/healthz")" "200"

echo "== auth"
for n in 1 2; do
  jar="$TMP/user$n.jar"
  api -X POST "$BASE/auth/register" -H 'Content-Type: application/json' \
    -d "{\"name\":\"smoke$n\",\"email\":\"smoke$n-$RUN@example.com\",\"password\":\"Passw0rd!$RUN\"}" >/dev/null
  c=$(code -c "$jar" -X POST "$BASE/auth/login" -H 'Content-Type: application/json' \
    -d "{\"email\":\"smoke$n-$RUN@example.com\",\"password\":\"Passw0rd!$RUN\"}")
  expect "login user$n" "$c" "200"
done
expect "/auth/me with cookie"    "$(code -b "$JAR1" "$BASE/auth/me")" "200"
expect "/auth/me without cookie" "$(code "$BASE/auth/me")" "401"

echo "== repo URL validation"
for bad_url in "file:///etc" "http://github.com/a/b" "https://user:pw@github.com/a/b" "--upload-pack=x"; do
  c=$(code -b "$JAR1" -X POST "$BASE/projects" -H 'Content-Type: application/json' \
    -d "{\"name\":\"bad\",\"repo_url\":\"$bad_url\",\"branch\":\"$BRANCH\"}")
  expect "reject repo_url '$bad_url'" "$c" "400"
done

echo "== create project"
P=$(api -b "$JAR1" -X POST "$BASE/projects" -H 'Content-Type: application/json' \
  -d "{\"name\":\"Smoke Test $RUN\",\"repo_url\":\"$REPO\",\"branch\":\"$BRANCH\",\"install_command\":\"npm install\",\"build_command\":\"$BUILD_CMD\",\"output_dir\":\"$OUT_DIR\"}")
PID=$(echo "$P" | jq -r '.project_id')
SLUG=$(echo "$P" | jq -r '.slug')
[ "$PID" != "null" ] && [ -n "$PID" ] && ok "project created ($SLUG)" || { bad "project create: $P"; exit 1; }
expect "no active deployment yet" "$(echo "$P" | jq -r '.active_deployment_id')" "null"
expect "project URL has slug" "$(echo "$P" | jq -r '.url' | grep -c "$SLUG")" "1"

deploy_and_wait() { # prints deployment id; returns final status in $STATUS
  local d did after=0 resp
  d=$(api -b "$JAR1" -X POST "$BASE/deployments" -H 'Content-Type: application/json' -d "{\"project_id\":\"$PID\"}")
  did=$(echo "$d" | jq -r '.id')
  for _ in $(seq 1 150); do
    resp=$(api -b "$JAR1" "$BASE/deployments/$did/logs?after=$after")
    after=$(echo "$resp" | jq -r '.next_after')
    if [ "$(echo "$resp" | jq -r '.done')" = "true" ]; then
      STATUS=$(echo "$resp" | jq -r '.status')
      LASTLOG=$(echo "$resp" | jq -r '.logs[-1].message')
      echo "$did"; return
    fi
    sleep 2
  done
  STATUS="timeout"; echo "$did"
}

echo "== deployment 1 (polls logs until done)"
D1=$(deploy_and_wait)
expect "deployment 1 status" "$STATUS" "success"
echo "        last log line: $LASTLOG"
expect "site served on project URL" "$(code -H "Host: $SLUG.localhost:8080" "$ROOT/")" "200"
PREFIX="${D1:0:12}"
expect "site served on deployment permalink" "$(code -H "Host: $PREFIX.localhost:8080" "$ROOT/")" "200"
expect "unknown slug returns 404" "$(code -H "Host: nope-0000.localhost:8080" "$ROOT/")" "404"
expect "wildcard prefix '_' returns 404" "$(code -H "Host: _.localhost:8080" "$ROOT/")" "404"

echo "== deployment 2, then rollback"
D2=$(deploy_and_wait)
expect "deployment 2 status" "$STATUS" "success"
ACTIVE=$(api -b "$JAR1" "$BASE/projects/$PID" | jq -r '.active_deployment_id')
expect "newest deployment is active" "$ACTIVE" "$D2"

c=$(code -b "$JAR1" -X POST "$BASE/projects/$PID/deployments/$D1/promote")
expect "promote deployment 1 (rollback)" "$c" "200"
ACTIVE=$(api -b "$JAR1" "$BASE/projects/$PID" | jq -r '.active_deployment_id')
expect "active_deployment_id is deployment 1" "$ACTIVE" "$D1"
IS_ACTIVE=$(api -b "$JAR1" "$BASE/projects/$PID/deployments" | jq -r ".[] | select(.id==\"$D1\") | .is_active")
expect "is_active true on deployment 1" "$IS_ACTIVE" "true"
expect "promote random uuid -> 404" "$(code -b "$JAR1" -X POST "$BASE/projects/$PID/deployments/00000000-0000-0000-0000-000000000000/promote")" "404"

echo "== authorization (user2 must not touch user1's data)"
expect "user2 reads user1 logs"       "$(code -b "$JAR2" "$BASE/deployments/$D1/logs")" "404"
expect "user2 promotes user1 deploy"  "$(code -b "$JAR2" -X POST "$BASE/projects/$PID/deployments/$D1/promote")" "404"
expect "user2 reads user1 project"    "$(code -b "$JAR2" "$BASE/projects/$PID")" "404"
expect "user2 deletes user1 project"  "$(code -b "$JAR2" -X DELETE "$BASE/projects/$PID")" "404"
expect "logs: unauthenticated"        "$(code "$BASE/deployments/$D1/logs")" "401"
expect "logs: after=abc"              "$(code -b "$JAR1" "$BASE/deployments/$D1/logs?after=abc")" "400"
expect "logs: after=-1"               "$(code -b "$JAR1" "$BASE/deployments/$D1/logs?after=-1")" "400"

echo "== webhook"
expect "unsigned webhook rejected" "$(code -X POST "$BASE/github/webhook" -H 'X-GitHub-Event: push' -H 'Content-Type: application/json' -d '{}')" "401"
expect "bad signature rejected"    "$(code -X POST "$BASE/github/webhook" -H 'X-GitHub-Event: push' -H 'X-Hub-Signature-256: sha256=deadbeef' -H 'Content-Type: application/json' -d '{}')" "401"

echo "== delete project"
expect "owner deletes project" "$(code -b "$JAR1" -X DELETE "$BASE/projects/$PID")" "200"
expect "project gone"          "$(code -b "$JAR1" "$BASE/projects/$PID")" "404"
expect "project URL now 404"   "$(code -H "Host: $SLUG.localhost:8080" "$ROOT/")" "404"
expect "deployment logs gone"  "$(code -b "$JAR1" "$BASE/deployments/$D1/logs")" "404"

echo "== logout"
expect "logout returns 200"        "$(code -b "$JAR1" -c "$JAR1" -X POST "$BASE/auth/logout")" "200"
expect "/auth/me after logout"     "$(code -b "$JAR1" "$BASE/auth/me")" "401"
expect "logout again without cookies" "$(code -X POST "$BASE/auth/logout")" "200"

echo
echo "passed: $PASS   failed: $FAIL"
rm -rf "$TMP"
[ "$FAIL" -eq 0 ]
