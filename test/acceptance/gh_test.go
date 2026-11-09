package acceptance

// The stand-in for gh, the GitHub command-line tool. Scripts find it on
// PATH as gh. It keeps the state of the fake GitHub in the folder
// $GH_STAND_IN:
//
//   - state.json: the logged-in account and the repos. A script gives it;
//     the stand-in rewrites it after each successful writing call.
//   - calls.log: one line per call. A gh api call is logged as
//     "api <METHOD> <path> | <arguments>", where the path has no leading
//     slash and no query; any other call as its arguments.
//   - fail: optional. Each line "<METHOD> <path> <status> <message>" makes
//     the matching gh api call fail with that HTTP status.
//
// It answers like gh 2.x against the GitHub REST API, for these commands:
//
//   - gh --version, gh auth status
//   - gh repo view [<owner>/<name>] [--json <fields>] [--jq <path>]
//   - gh api <endpoint> with -X/--method, -f/--raw-field, -F/--field,
//     --input, -i/--include, --silent, --jq, and the ignored -H/--header,
//     --paginate, --slurp, --hostname, -p/--preview. Without -X, the
//     method is POST when fields or --input are given, else GET.
//   - gh secret set <name> --env <env> [--repo <owner>/<name>] [--body
//     <value>]; without --body, the value is standard input without its
//     trailing line breaks.
//   - gh secret list --env <env> [--repo <owner>/<name>] [--json <fields>]
//     [--jq <path>]
//
// --jq takes only paths such as .a.b, .[] and .[].name. Endpoints:
// user; user/repository_invitations/<id>; and below repos/<owner>/<name>:
// the repo itself, collaborators, collaborators/<login>,
// collaborators/<login>/permission, invitations, invitations/<id>,
// environments, environments/<env>, its deployment-branch-policies and
// secrets (list, public-key, <name>), rulesets and rulesets/<id>.
// {owner} and {repo} in an endpoint stand for the current repo.
//
// Not logged in, gh auth status exits 1 and every other command exits 4.
// An unknown repo, route or id is a 404. An HTTP error prints the error
// body on stdout and "gh: <message> (HTTP <status>)" as the first line on
// stderr, and exits 1.

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type ghState struct {
	LoggedIn bool               `json:"logged_in"`
	User     string             `json:"user"`
	Current  string             `json:"current,omitempty"`
	Repos    map[string]*ghRepo `json:"repos"`
}

type ghRepo struct {
	Admin    bool           `json:"admin"`
	Settings map[string]any `json:"settings,omitempty"`
	// Collaborators maps a login to its role: read, triage, write,
	// maintain or admin.
	Collaborators map[string]string         `json:"collaborators,omitempty"`
	Invitations   []*ghInvitation           `json:"invitations,omitempty"`
	Environments  map[string]*ghEnvironment `json:"environments,omitempty"`
	Rulesets      []map[string]any          `json:"rulesets,omitempty"`
}

type ghInvitation struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
	Role  string `json:"role"`
}

type ghEnvironment struct {
	DeploymentBranchPolicy map[string]any    `json:"deployment_branch_policy"`
	BranchPolicies         []*ghBranchPolicy `json:"branch_policies,omitempty"`
	Secrets                map[string]string `json:"secrets,omitempty"`
}

type ghBranchPolicy struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// ghCall is one run of the stand-in.
type ghCall struct {
	dir    string
	state  *ghState
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

// httpError is an HTTP error of the fake GitHub.
type httpError struct {
	status  int
	message string
	extra   string
}

func ghStandIn() int {
	dir := os.Getenv("GH_STAND_IN")
	if dir == "" {
		fmt.Fprintln(os.Stderr, "gh stand-in: GH_STAND_IN is not set")
		return 1
	}
	c := &ghCall{dir: dir, stdin: os.Stdin, stdout: os.Stdout, stderr: os.Stderr}
	return c.run(os.Args[1:])
}

func (c *ghCall) run(args []string) int {
	// A gh api call is logged by api, which knows method and path.
	if len(args) == 0 || args[0] != "api" {
		c.log(strings.Join(args, " "))
	}
	if len(args) > 0 && (args[0] == "--version" || args[0] == "version") {
		fmt.Fprintln(c.stdout, "gh version 2.63.0 (stand-in)")
		return 0
	}
	state, err := c.load()
	if err != nil {
		if len(args) > 0 && args[0] == "api" {
			c.log("api ? ? | " + strings.Join(args[1:], " "))
		}
		fmt.Fprintln(c.stderr, "gh stand-in:", err)
		return 1
	}
	c.state = state
	if len(args) >= 2 && args[0] == "auth" && args[1] == "status" {
		if !state.LoggedIn {
			fmt.Fprintln(c.stderr, "You are not logged into any GitHub hosts. To log in, run: gh auth login")
			return 1
		}
		fmt.Fprintf(c.stdout, "github.com\n  ✓ Logged in to github.com account %s (keyring)\n", state.User)
		return 0
	}
	if !state.LoggedIn {
		if len(args) > 0 && args[0] == "api" {
			c.log("api ? ? | " + strings.Join(args[1:], " "))
		}
		fmt.Fprintln(c.stderr, "To get started with GitHub CLI, please run:  gh auth login")
		fmt.Fprintln(c.stderr, "Alternatively, populate the GH_TOKEN environment variable with a GitHub API authentication token.")
		return 4
	}
	switch {
	case len(args) >= 1 && args[0] == "api":
		return c.api(args[1:])
	case len(args) >= 2 && args[0] == "repo" && args[1] == "view":
		return c.repoView(args[2:])
	case len(args) >= 2 && args[0] == "secret" && args[1] == "set":
		return c.secretSet(args[2:])
	case len(args) >= 2 && args[0] == "secret" && args[1] == "list":
		return c.secretList(args[2:])
	}
	return c.unsupported("command: gh " + strings.Join(args, " "))
}

func (c *ghCall) unsupported(what string) int {
	fmt.Fprintln(c.stderr, "gh stand-in: unsupported "+what)
	return 1
}

func (c *ghCall) log(line string) {
	f, err := os.OpenFile(filepath.Join(c.dir, "calls.log"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
	if err != nil {
		fmt.Fprintln(c.stderr, "gh stand-in:", err)
		return
	}
	defer f.Close()
	fmt.Fprintln(f, strings.ReplaceAll(line, "\n", `\n`))
}

func (c *ghCall) load() (*ghState, error) {
	data, err := os.ReadFile(filepath.Join(c.dir, "state.json"))
	if err != nil {
		return nil, err
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	var s ghState
	if err := dec.Decode(&s); err != nil {
		return nil, fmt.Errorf("state.json: %v", err)
	}
	if s.Repos == nil {
		s.Repos = map[string]*ghRepo{}
	}
	for _, r := range s.Repos {
		if r.Settings == nil {
			r.Settings = map[string]any{}
		}
		if r.Collaborators == nil {
			r.Collaborators = map[string]string{}
		}
		if r.Environments == nil {
			r.Environments = map[string]*ghEnvironment{}
		}
		for _, e := range r.Environments {
			if e.Secrets == nil {
				e.Secrets = map[string]string{}
			}
		}
	}
	return &s, nil
}

func (c *ghCall) save() error {
	data, err := json.MarshalIndent(c.state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(c.dir, "state.json"), append(data, '\n'), 0o644)
}

// ---- gh repo view

func (c *ghCall) repoView(args []string) int {
	name, jsonFields, jq := "", "", ""
	for i := 0; i < len(args); i++ {
		a := args[i]
		flag, value, hasValue := splitFlag(a)
		switch flag {
		case "--json", "--jq", "-q":
			if !hasValue {
				if i+1 >= len(args) {
					return c.unsupported("use of " + flag + " without a value")
				}
				i++
				value = args[i]
			}
			if flag == "--json" {
				jsonFields = value
			} else {
				jq = value
			}
		default:
			if strings.HasPrefix(a, "-") || name != "" {
				return c.unsupported("argument of gh repo view: " + a)
			}
			name = a
		}
	}
	if name == "" {
		name = c.state.Current
		if name == "" {
			fmt.Fprintln(c.stderr, "could not determine base repo: no git remotes found")
			return 1
		}
	}
	repo := c.state.Repos[name]
	if repo == nil {
		fmt.Fprintf(c.stderr, "GraphQL: Could not resolve to a Repository with the name '%s'. (repository)\n", name)
		return 1
	}
	if jsonFields == "" {
		fmt.Fprintf(c.stdout, "name:\t%s\n", name)
		return 0
	}
	owner, short, _ := strings.Cut(name, "/")
	permission := "WRITE"
	if repo.Admin {
		permission = "ADMIN"
	}
	all := map[string]any{
		"nameWithOwner":    name,
		"name":             short,
		"owner":            map[string]any{"login": owner},
		"url":              "https://github.com/" + name,
		"viewerPermission": permission,
	}
	out := map[string]any{}
	for _, f := range strings.Split(jsonFields, ",") {
		v, ok := all[f]
		if !ok {
			return c.unsupported("field of gh repo view --json: " + f)
		}
		out[f] = v
	}
	return c.print(out, jq)
}

// ---- gh secret

func (c *ghCall) secretFlags(args []string, allowed map[string]bool) (positional []string, flags map[string]string, bad string) {
	flags = map[string]string{}
	aliases := map[string]string{"-e": "--env", "-R": "--repo", "-b": "--body", "-a": "--app", "-q": "--jq"}
	for i := 0; i < len(args); i++ {
		a := args[i]
		if !strings.HasPrefix(a, "-") {
			positional = append(positional, a)
			continue
		}
		flag, value, hasValue := splitFlag(a)
		if long, ok := aliases[flag]; ok {
			flag = long
		}
		if !allowed[flag] {
			return nil, nil, a
		}
		if !hasValue {
			if i+1 >= len(args) {
				return nil, nil, a
			}
			i++
			value = args[i]
		}
		flags[flag] = value
	}
	return positional, flags, ""
}

func (c *ghCall) secretRepo(flags map[string]string) (string, *ghRepo, int) {
	name := flags["--repo"]
	if name == "" {
		name = c.state.Current
	}
	repo := c.state.Repos[name]
	if repo == nil {
		fmt.Fprintln(c.stderr, "failed to fetch public key: HTTP 404: Not Found (https://api.github.com/repos/"+name+")")
		return "", nil, 1
	}
	return name, repo, 0
}

func (c *ghCall) secretSet(args []string) int {
	positional, flags, bad := c.secretFlags(args, map[string]bool{"--env": true, "--repo": true, "--body": true, "--app": true})
	if bad != "" || len(positional) != 1 {
		return c.unsupported("arguments of gh secret set: " + strings.Join(args, " "))
	}
	if app, ok := flags["--app"]; ok && app != "actions" {
		return c.unsupported("--app " + app)
	}
	envName, ok := flags["--env"]
	if !ok {
		return c.unsupported("gh secret set without --env")
	}
	_, repo, code := c.secretRepo(flags)
	if code != 0 {
		return code
	}
	env := repo.Environments[envName]
	if env == nil {
		fmt.Fprintln(c.stderr, "failed to fetch public key: HTTP 404: Not Found")
		return 1
	}
	value, ok := flags["--body"]
	if !ok {
		data, err := io.ReadAll(c.stdin)
		if err != nil {
			fmt.Fprintln(c.stderr, "failed to read from standard input:", err)
			return 1
		}
		value = strings.TrimRight(string(data), "\r\n")
	}
	env.Secrets[positional[0]] = value
	if err := c.save(); err != nil {
		fmt.Fprintln(c.stderr, "gh stand-in:", err)
		return 1
	}
	return 0
}

func (c *ghCall) secretList(args []string) int {
	positional, flags, bad := c.secretFlags(args, map[string]bool{"--env": true, "--repo": true, "--json": true, "--jq": true, "--app": true})
	if bad != "" || len(positional) != 0 {
		return c.unsupported("arguments of gh secret list: " + strings.Join(args, " "))
	}
	envName, ok := flags["--env"]
	if !ok {
		return c.unsupported("gh secret list without --env")
	}
	_, repo, code := c.secretRepo(flags)
	if code != 0 {
		return code
	}
	env := repo.Environments[envName]
	if env == nil {
		fmt.Fprintln(c.stderr, "failed to get secrets: HTTP 404: Not Found")
		return 1
	}
	names := sortedKeys(env.Secrets)
	if fields, ok := flags["--json"]; ok {
		list := []any{}
		for _, n := range names {
			item := map[string]any{}
			for _, f := range strings.Split(fields, ",") {
				switch f {
				case "name":
					item["name"] = n
				case "updatedAt":
					item["updatedAt"] = "2026-01-01T00:00:00Z"
				default:
					return c.unsupported("field of gh secret list --json: " + f)
				}
			}
			list = append(list, item)
		}
		return c.print(list, flags["--jq"])
	}
	for _, n := range names {
		fmt.Fprintf(c.stdout, "%s\t2026-01-01T00:00:00Z\n", n)
	}
	return 0
}

// ---- gh api

func (c *ghCall) api(args []string) int {
	var (
		method, endpoint, input, jq string
		include, silent             bool
		body                        map[string]any
		hasFields                   bool
	)
	logBad := func(what string) int {
		c.log("api ? ? | " + strings.Join(args, " "))
		return c.unsupported(what)
	}
	for i := 0; i < len(args); i++ {
		a := args[i]
		if !strings.HasPrefix(a, "-") || a == "-" {
			if endpoint != "" {
				return logBad("second endpoint of gh api: " + a)
			}
			endpoint = a
			continue
		}
		flag, value, hasValue := splitFlag(a)
		if strings.HasPrefix(flag, "-X") && len(flag) > 2 {
			flag, value, hasValue = "-X", flag[2:], true
		}
		next := func() bool {
			if hasValue {
				return true
			}
			if i+1 >= len(args) {
				return false
			}
			i++
			value = args[i]
			return true
		}
		switch flag {
		case "-X", "--method":
			if !next() {
				return logBad("use of " + flag)
			}
			method = strings.ToUpper(value)
		case "-f", "--raw-field", "-F", "--field":
			if !next() {
				return logBad("use of " + flag)
			}
			key, raw, ok := strings.Cut(value, "=")
			if !ok {
				return logBad("field without =: " + value)
			}
			var v any = raw
			if flag == "-F" || flag == "--field" {
				typed, err := c.typedField(raw)
				if err != nil {
					return logBad(err.Error())
				}
				v = typed
			}
			if body == nil {
				body = map[string]any{}
			}
			if err := setField(body, key, v); err != nil {
				return logBad(err.Error())
			}
			hasFields = true
		case "--input":
			if !next() {
				return logBad("use of --input")
			}
			input = value
		case "-q", "--jq":
			if !next() {
				return logBad("use of --jq")
			}
			jq = value
		case "-H", "--header", "--hostname", "-p", "--preview":
			if !next() {
				return logBad("use of " + flag)
			}
		case "-i", "--include":
			include = true
		case "--silent":
			silent = true
		case "--paginate", "--slurp":
		default:
			return logBad("flag of gh api: " + a)
		}
	}
	if endpoint == "" {
		return logBad("gh api without an endpoint")
	}
	if method == "" {
		method = "GET"
		if hasFields || input != "" {
			method = "POST"
		}
	}
	if input != "" {
		var data []byte
		var err error
		if input == "-" {
			data, err = io.ReadAll(c.stdin)
		} else {
			data, err = os.ReadFile(input)
		}
		if err != nil {
			return logBad("--input: " + err.Error())
		}
		dec := json.NewDecoder(bytes.NewReader(data))
		dec.UseNumber()
		var v map[string]any
		if err := dec.Decode(&v); err != nil {
			return logBad("--input is no JSON object: " + err.Error())
		}
		body = v
	} else if method == "GET" {
		body = nil // fields of a GET are query parameters
	}
	path := strings.TrimPrefix(endpoint, "/")
	path, _, _ = strings.Cut(path, "?")
	if c.state.Current != "" {
		owner, name, _ := strings.Cut(c.state.Current, "/")
		path = strings.ReplaceAll(path, "{owner}", owner)
		path = strings.ReplaceAll(path, "{repo}", name)
	}
	c.log(fmt.Sprintf("api %s %s | %s", method, path, strings.Join(args, " ")))

	status, result, herr := 0, any(nil), (*httpError)(nil)
	if herr = c.injected(method, path); herr == nil {
		status, result, herr = c.route(method, path, body)
	}
	if herr != nil {
		errBody := map[string]any{
			"message":           herr.message,
			"documentation_url": "https://docs.github.com/rest",
			"status":            strconv.Itoa(herr.status),
		}
		if include {
			fmt.Fprintf(c.stdout, "HTTP/2.0 %d %s\r\nContent-Type: application/json; charset=utf-8\r\n\r\n", herr.status, herr.message)
		}
		fmt.Fprintln(c.stdout, encodeJSON(errBody))
		fmt.Fprintf(c.stderr, "gh: %s (HTTP %d)\n", herr.message, herr.status)
		if herr.extra != "" {
			fmt.Fprintln(c.stderr, herr.extra)
		}
		return 1
	}
	if method != "GET" {
		if err := c.save(); err != nil {
			fmt.Fprintln(c.stderr, "gh stand-in:", err)
			return 1
		}
	}
	if include {
		fmt.Fprintf(c.stdout, "HTTP/2.0 %d %s\r\nContent-Type: application/json; charset=utf-8\r\n\r\n", status, statusText(status))
	}
	if silent || status == 204 {
		return 0
	}
	return c.print(result, jq)
}

func statusText(status int) string {
	switch status {
	case 200:
		return "OK"
	case 201:
		return "Created"
	case 204:
		return "No Content"
	}
	return ""
}

func (c *ghCall) typedField(raw string) (any, error) {
	switch {
	case raw == "true":
		return true, nil
	case raw == "false":
		return false, nil
	case raw == "null":
		return nil, nil
	case regexp.MustCompile(`^-?[0-9]+$`).MatchString(raw):
		return json.Number(raw), nil
	case strings.HasPrefix(raw, "@"):
		var data []byte
		var err error
		if raw == "@-" {
			data, err = io.ReadAll(c.stdin)
		} else {
			data, err = os.ReadFile(raw[1:])
		}
		return string(data), err
	}
	return raw, nil
}

// setField sets key, such as a, a[b] or a[], in body, as gh api does.
func setField(body map[string]any, key string, v any) error {
	name, rest, nested := strings.Cut(key, "[")
	if !nested {
		body[key] = v
		return nil
	}
	parts := []string{name}
	for _, p := range strings.Split(strings.TrimSuffix(rest, "]"), "][") {
		parts = append(parts, p)
	}
	var cur any = body
	for i, p := range parts {
		last := i == len(parts)-1
		m, ok := cur.(map[string]any)
		if !ok {
			return errors.New("field key " + key)
		}
		if last {
			m[p] = v
			return nil
		}
		if parts[i+1] == "" {
			list, _ := m[p].([]any)
			if i+1 == len(parts)-1 {
				m[p] = append(list, v)
				return nil
			}
			return errors.New("field key " + key)
		}
		next, ok := m[p].(map[string]any)
		if !ok {
			next = map[string]any{}
			m[p] = next
		}
		cur = next
	}
	return nil
}

func (c *ghCall) injected(method, path string) *httpError {
	data, err := os.ReadFile(filepath.Join(c.dir, "fail"))
	if err != nil {
		return nil
	}
	sc := bufio.NewScanner(bytes.NewReader(data))
	for sc.Scan() {
		f := strings.SplitN(strings.TrimSpace(sc.Text()), " ", 4)
		if len(f) < 4 || f[0] != method || strings.TrimPrefix(f[1], "/") != path {
			continue
		}
		status, err := strconv.Atoi(f[2])
		if err != nil {
			continue
		}
		return &httpError{status: status, message: f[3], extra: "stand-in: refused by the fail file"}
	}
	return nil
}

func notFound() *httpError { return &httpError{status: 404, message: "Not Found"} }

func invalid(detail string) *httpError {
	return &httpError{status: 422, message: "Validation Failed", extra: detail}
}

func (c *ghCall) route(method, path string, body map[string]any) (int, any, *httpError) {
	seg := strings.Split(path, "/")
	switch {
	case path == "user" && method == "GET":
		return 200, map[string]any{"login": c.state.User}, nil
	case len(seg) == 3 && seg[0] == "user" && seg[1] == "repository_invitations":
		return c.answerInvitation(method, seg[2])
	case len(seg) >= 3 && seg[0] == "repos":
		name := seg[1] + "/" + seg[2]
		repo := c.state.Repos[name]
		if repo == nil {
			return 0, nil, notFound()
		}
		return c.repoRoute(name, repo, method, seg[3:], body)
	}
	return 0, nil, notFound()
}

func (c *ghCall) answerInvitation(method, idText string) (int, any, *httpError) {
	id, err := strconv.Atoi(idText)
	if err != nil {
		return 0, nil, notFound()
	}
	for _, repo := range c.state.Repos {
		for i, inv := range repo.Invitations {
			if inv.ID != id {
				continue
			}
			switch method {
			case "PATCH":
				repo.Collaborators[inv.Login] = inv.Role
			case "DELETE":
			default:
				return 0, nil, notFound()
			}
			repo.Invitations = append(repo.Invitations[:i], repo.Invitations[i+1:]...)
			return 204, nil, nil
		}
	}
	return 0, nil, notFound()
}

var roleOfPermission = map[string]string{"pull": "read", "triage": "triage", "push": "write", "maintain": "maintain", "admin": "admin"}

func (c *ghCall) repoJSON(name string, repo *ghRepo) map[string]any {
	owner, short, _ := strings.Cut(name, "/")
	out := map[string]any{}
	for k, v := range repo.Settings {
		out[k] = v
	}
	out["name"] = short
	out["full_name"] = name
	out["owner"] = map[string]any{"login": owner, "type": "User"}
	out["private"] = false
	out["default_branch"] = "main"
	out["permissions"] = map[string]any{"admin": repo.Admin, "maintain": repo.Admin, "push": true, "triage": true, "pull": true}
	return out
}

func (c *ghCall) repoRoute(name string, repo *ghRepo, method string, seg []string, body map[string]any) (int, any, *httpError) {
	if len(seg) == 0 {
		switch method {
		case "GET":
			return 200, c.repoJSON(name, repo), nil
		case "PATCH":
			for k, v := range body {
				switch k {
				case "name", "full_name", "owner", "permissions", "id":
					return 0, nil, invalid("stand-in: " + k + " is not changed here")
				}
				repo.Settings[k] = v
			}
			return 200, c.repoJSON(name, repo), nil
		}
		return 0, nil, notFound()
	}
	switch seg[0] {
	case "collaborators":
		return c.collaborators(name, repo, method, seg[1:], body)
	case "invitations":
		return c.invitations(name, repo, method, seg[1:])
	case "environments":
		return c.environments(repo, method, seg[1:], body)
	case "rulesets":
		return c.rulesets(name, repo, method, seg[1:], body)
	}
	return 0, nil, notFound()
}

func permissionOfRole(role string) string {
	switch role {
	case "maintain":
		return "write"
	case "triage":
		return "read"
	}
	return role
}

func (c *ghCall) collaborators(name string, repo *ghRepo, method string, seg []string, body map[string]any) (int, any, *httpError) {
	switch {
	case len(seg) == 0 && method == "GET":
		list := []any{}
		for _, login := range sortedKeys(repo.Collaborators) {
			role := repo.Collaborators[login]
			list = append(list, map[string]any{
				"login":       login,
				"role_name":   role,
				"permissions": rolePermissions(role),
			})
		}
		return 200, list, nil
	case len(seg) == 1:
		login := seg[0]
		_, isCollaborator := repo.Collaborators[login]
		switch method {
		case "GET":
			if isCollaborator {
				return 204, nil, nil
			}
			return 0, nil, notFound()
		case "PUT":
			permission := "push"
			if p, ok := body["permission"].(string); ok {
				permission = p
			}
			role, ok := roleOfPermission[permission]
			if !ok {
				return 0, nil, invalid("permission " + permission)
			}
			if isCollaborator {
				return 204, nil, nil
			}
			for _, inv := range repo.Invitations {
				if inv.Login == login {
					inv.Role = role
					return 201, c.invitationJSON(name, inv), nil
				}
			}
			inv := &ghInvitation{ID: c.nextInvitationID(), Login: login, Role: role}
			repo.Invitations = append(repo.Invitations, inv)
			return 201, c.invitationJSON(name, inv), nil
		case "DELETE":
			delete(repo.Collaborators, login)
			return 204, nil, nil
		}
	case len(seg) == 2 && seg[1] == "permission" && method == "GET":
		role, ok := repo.Collaborators[seg[0]]
		if !ok {
			return 0, nil, notFound()
		}
		return 200, map[string]any{
			"permission": permissionOfRole(role),
			"role_name":  role,
			"user":       map[string]any{"login": seg[0]},
		}, nil
	}
	return 0, nil, notFound()
}

func rolePermissions(role string) map[string]any {
	rank := map[string]int{"read": 1, "triage": 2, "write": 3, "maintain": 4, "admin": 5}[role]
	return map[string]any{"pull": rank >= 1, "triage": rank >= 2, "push": rank >= 3, "maintain": rank >= 4, "admin": rank >= 5}
}

func (c *ghCall) nextInvitationID() int {
	id := 0
	for _, r := range c.state.Repos {
		for _, inv := range r.Invitations {
			id = max(id, inv.ID)
		}
	}
	return id + 1
}

func (c *ghCall) invitationJSON(name string, inv *ghInvitation) map[string]any {
	return map[string]any{
		"id":          inv.ID,
		"invitee":     map[string]any{"login": inv.Login},
		"inviter":     map[string]any{"login": c.state.User},
		"permissions": inv.Role,
		"repository":  map[string]any{"full_name": name},
	}
}

func (c *ghCall) invitations(name string, repo *ghRepo, method string, seg []string) (int, any, *httpError) {
	if len(seg) == 0 && method == "GET" {
		list := []any{}
		for _, inv := range repo.Invitations {
			list = append(list, c.invitationJSON(name, inv))
		}
		return 200, list, nil
	}
	if len(seg) == 1 && method == "DELETE" {
		for i, inv := range repo.Invitations {
			if strconv.Itoa(inv.ID) == seg[0] {
				repo.Invitations = append(repo.Invitations[:i], repo.Invitations[i+1:]...)
				return 204, nil, nil
			}
		}
	}
	return 0, nil, notFound()
}

func environmentJSON(name string, env *ghEnvironment) map[string]any {
	var policy any
	if env.DeploymentBranchPolicy != nil {
		policy = env.DeploymentBranchPolicy
	}
	return map[string]any{"name": name, "deployment_branch_policy": policy, "protection_rules": []any{}}
}

func (c *ghCall) environments(repo *ghRepo, method string, seg []string, body map[string]any) (int, any, *httpError) {
	if len(seg) == 0 {
		if method != "GET" {
			return 0, nil, notFound()
		}
		list := []any{}
		for _, n := range sortedKeys(repo.Environments) {
			list = append(list, environmentJSON(n, repo.Environments[n]))
		}
		return 200, map[string]any{"total_count": len(list), "environments": list}, nil
	}
	envName := seg[0]
	env := repo.Environments[envName]
	if len(seg) == 1 {
		switch method {
		case "GET":
			if env == nil {
				return 0, nil, notFound()
			}
			return 200, environmentJSON(envName, env), nil
		case "PUT":
			if env == nil {
				env = &ghEnvironment{Secrets: map[string]string{}}
			}
			if raw, ok := body["deployment_branch_policy"]; ok {
				if raw == nil {
					env.DeploymentBranchPolicy = nil
					env.BranchPolicies = nil
				} else {
					policy, ok := raw.(map[string]any)
					if !ok {
						return 0, nil, invalid("deployment_branch_policy is no object")
					}
					protected, ok1 := policy["protected_branches"].(bool)
					custom, ok2 := policy["custom_branch_policies"].(bool)
					if !ok1 || !ok2 || protected == custom {
						return 0, nil, invalid("deployment_branch_policy needs exactly one of protected_branches and custom_branch_policies true")
					}
					env.DeploymentBranchPolicy = map[string]any{"protected_branches": protected, "custom_branch_policies": custom}
					if !custom {
						env.BranchPolicies = nil
					}
				}
			}
			repo.Environments[envName] = env
			return 200, environmentJSON(envName, env), nil
		case "DELETE":
			if env == nil {
				return 0, nil, notFound()
			}
			delete(repo.Environments, envName)
			return 204, nil, nil
		}
		return 0, nil, notFound()
	}
	if env == nil {
		return 0, nil, notFound()
	}
	switch seg[1] {
	case "deployment-branch-policies":
		return c.branchPolicies(env, method, seg[2:], body)
	case "secrets":
		return c.secrets(env, method, seg[2:])
	}
	return 0, nil, notFound()
}

func policyJSON(p *ghBranchPolicy) map[string]any {
	return map[string]any{"id": p.ID, "node_id": fmt.Sprintf("BP_%d", p.ID), "name": p.Name, "type": p.Type}
}

func (c *ghCall) branchPolicies(env *ghEnvironment, method string, seg []string, body map[string]any) (int, any, *httpError) {
	custom, _ := env.DeploymentBranchPolicy["custom_branch_policies"].(bool)
	if len(seg) == 0 {
		switch method {
		case "GET":
			list := []any{}
			for _, p := range env.BranchPolicies {
				list = append(list, policyJSON(p))
			}
			return 200, map[string]any{"total_count": len(list), "branch_policies": list}, nil
		case "POST":
			if !custom {
				return 0, nil, notFound()
			}
			name, _ := body["name"].(string)
			if name == "" {
				return 0, nil, invalid("name is missing")
			}
			typ := "branch"
			if t, ok := body["type"].(string); ok {
				typ = t
			}
			id := 0
			for _, p := range env.BranchPolicies {
				if p.Name == name && p.Type == typ {
					return 0, nil, invalid("name has already been taken")
				}
				id = max(id, p.ID)
			}
			p := &ghBranchPolicy{ID: id + 1, Name: name, Type: typ}
			env.BranchPolicies = append(env.BranchPolicies, p)
			return 200, policyJSON(p), nil
		}
		return 0, nil, notFound()
	}
	if len(seg) != 1 {
		return 0, nil, notFound()
	}
	for i, p := range env.BranchPolicies {
		if strconv.Itoa(p.ID) != seg[0] {
			continue
		}
		switch method {
		case "GET":
			return 200, policyJSON(p), nil
		case "PUT":
			if name, ok := body["name"].(string); ok && name != "" {
				p.Name = name
			}
			return 200, policyJSON(p), nil
		case "DELETE":
			env.BranchPolicies = append(env.BranchPolicies[:i], env.BranchPolicies[i+1:]...)
			return 204, nil, nil
		}
	}
	return 0, nil, notFound()
}

func (c *ghCall) secrets(env *ghEnvironment, method string, seg []string) (int, any, *httpError) {
	secretJSON := func(n string) map[string]any {
		return map[string]any{"name": n, "created_at": "2026-01-01T00:00:00Z", "updated_at": "2026-01-01T00:00:00Z"}
	}
	switch {
	case len(seg) == 0 && method == "GET":
		list := []any{}
		for _, n := range sortedKeys(env.Secrets) {
			list = append(list, secretJSON(n))
		}
		return 200, map[string]any{"total_count": len(list), "secrets": list}, nil
	case len(seg) == 1 && seg[0] == "public-key" && method == "GET":
		return 200, map[string]any{"key_id": "1", "key": "c3RhbmQtaW4ga2V5IG9mIDMyIGJ5dGVzIGZvciB0ZXN0cw=="}, nil
	case len(seg) == 1 && method == "GET":
		if _, ok := env.Secrets[seg[0]]; ok {
			return 200, secretJSON(seg[0]), nil
		}
	case len(seg) == 1 && method == "DELETE":
		if _, ok := env.Secrets[seg[0]]; ok {
			delete(env.Secrets, seg[0])
			return 204, nil, nil
		}
	}
	return 0, nil, notFound()
}

func rulesetID(r map[string]any) string {
	return fmt.Sprint(r["id"])
}

func rulesetSummary(r map[string]any) map[string]any {
	out := map[string]any{}
	for _, k := range []string{"id", "name", "target", "source_type", "source", "enforcement", "node_id"} {
		if v, ok := r[k]; ok {
			out[k] = v
		}
	}
	return out
}

func (c *ghCall) rulesets(name string, repo *ghRepo, method string, seg []string, body map[string]any) (int, any, *httpError) {
	nameTaken := func(n string, except string) bool {
		for _, r := range repo.Rulesets {
			if r["name"] == n && rulesetID(r) != except {
				return true
			}
		}
		return false
	}
	if len(seg) == 0 {
		switch method {
		case "GET":
			list := []any{}
			for _, r := range repo.Rulesets {
				list = append(list, rulesetSummary(r))
			}
			return 200, list, nil
		case "POST":
			n, _ := body["name"].(string)
			if n == "" {
				return 0, nil, invalid("name is missing")
			}
			if _, ok := body["enforcement"].(string); !ok {
				return 0, nil, invalid("enforcement is missing")
			}
			if nameTaken(n, "") {
				return 0, nil, invalid("name must be unique")
			}
			id := 0
			for _, r := range repo.Rulesets {
				existing, _ := strconv.Atoi(rulesetID(r))
				id = max(id, existing)
			}
			r := map[string]any{}
			for k, v := range body {
				r[k] = v
			}
			r["id"] = json.Number(strconv.Itoa(id + 1))
			r["node_id"] = fmt.Sprintf("RRS_%d", id+1)
			r["source_type"] = "Repository"
			r["source"] = name
			repo.Rulesets = append(repo.Rulesets, r)
			return 201, r, nil
		}
		return 0, nil, notFound()
	}
	if len(seg) != 1 {
		return 0, nil, notFound()
	}
	for i, r := range repo.Rulesets {
		if rulesetID(r) != seg[0] {
			continue
		}
		switch method {
		case "GET":
			return 200, r, nil
		case "PUT":
			if n, ok := body["name"].(string); ok && nameTaken(n, seg[0]) {
				return 0, nil, invalid("name must be unique")
			}
			for k, v := range body {
				switch k {
				case "id", "node_id", "source_type", "source":
					continue
				}
				r[k] = v
			}
			return 200, r, nil
		case "DELETE":
			repo.Rulesets = append(repo.Rulesets[:i], repo.Rulesets[i+1:]...)
			return 204, nil, nil
		}
	}
	return 0, nil, notFound()
}

// ---- output

// print writes v as compact JSON, or the results of the --jq path, one per
// line, with strings unquoted.
func (c *ghCall) print(v any, jq string) int {
	if jq == "" {
		fmt.Fprintln(c.stdout, encodeJSON(v))
		return 0
	}
	results, err := jqPath(v, jq)
	if err != nil {
		return c.unsupported(err.Error())
	}
	for _, r := range results {
		if s, ok := r.(string); ok {
			fmt.Fprintln(c.stdout, s)
		} else {
			fmt.Fprintln(c.stdout, encodeJSON(r))
		}
	}
	return 0
}

var jqSegment = regexp.MustCompile(`^\.(\[\]|[A-Za-z_][A-Za-z0-9_]*)`)

func jqPath(v any, expr string) ([]any, error) {
	if expr == "." {
		return []any{v}, nil
	}
	values := []any{v}
	for rest := expr; rest != ""; {
		m := jqSegment.FindString(rest)
		if m == "" {
			return nil, errors.New("--jq expression: " + expr)
		}
		rest = rest[len(m):]
		var next []any
		for _, cur := range values {
			if m == ".[]" {
				switch x := cur.(type) {
				case []any:
					next = append(next, x...)
				case map[string]any:
					for _, k := range sortedKeys(x) {
						next = append(next, x[k])
					}
				default:
					return nil, errors.New("--jq cannot iterate over " + encodeJSON(cur))
				}
				continue
			}
			obj, _ := cur.(map[string]any)
			next = append(next, obj[m[1:]])
		}
		values = next
	}
	return values, nil
}

func splitFlag(a string) (flag, value string, hasValue bool) {
	if strings.HasPrefix(a, "--") {
		if f, v, ok := strings.Cut(a, "="); ok {
			return f, v, true
		}
	}
	return a, "", false
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
