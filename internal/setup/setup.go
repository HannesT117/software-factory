// Package setup applies the admin part of the setup of a GitHub repo
// through the engineer's own gh: merge settings, the invitation of the
// agent account, the environment cross-review with its secret, and the
// rulesets main and branches. Each step reads the current state, compares
// it with the target, and writes only when they differ.
package setup

import (
	"bufio"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strings"

	"github.com/HannesT117/software-factory/internal/report"
)

//go:embed rulesets/main.json rulesets/branches.json
var rulesetFiles embed.FS

// The names that the setup fixes.
const (
	Environment = "cross-review"
	SecretName  = "OPPER_API_KEY"
	steps       = 6
)

// Options are the arguments of factory setup.
type Options struct {
	Agent  string
	Repo   string
	DryRun bool
	Stdin  io.Reader
}

var (
	loginPattern = regexp.MustCompile(`^[A-Za-z0-9](?:[A-Za-z0-9-]*[A-Za-z0-9])?$`)
	repoPattern  = regexp.MustCompile(`^[A-Za-z0-9-]+/[A-Za-z0-9._-]+$`)
)

// ValidLogin reports whether login has the form of a GitHub login.
func ValidLogin(login string) bool { return loginPattern.MatchString(login) }

// run holds one run of factory setup.
type run struct {
	r    *report.Report
	gh   gh
	opts Options
	repo string
}

// stepError is an error of a step: its code and detail.
type stepError struct {
	code   string
	detail string
}

func (e *stepError) Error() string { return e.code + " " + e.detail }

// fail turns an error of a call into the error of a step.
func fail(err error) *stepError {
	var se *stepError
	if errors.As(err, &se) {
		return se
	}
	return &stepError{code: "api-failed", detail: err.Error()}
}

// Run applies the setup and returns the exit code.
func Run(r *report.Report, opts Options) int {
	path, err := exec.LookPath("gh")
	if err != nil {
		return r.Abort("gh", "no-gh", "gh is not on PATH")
	}
	u := &run{r: r, gh: gh{path: path}, opts: opts}
	if _, err := u.gh.run("gh auth status", nil, "auth", "status"); err != nil {
		return r.Abort("gh", "no-gh", "gh is not logged in; run gh auth login")
	}
	if code, ok := u.checkRepo(); !ok {
		return code
	}
	if code, ok := u.checkAgent(); !ok {
		return code
	}
	for _, s := range []struct {
		name string
		fn   func() (bool, string, error)
	}{
		{"merge-settings", u.mergeSettings},
		{"agent-account", u.agentAccount},
		{"environment", u.environment},
		{"secret", u.secret},
		{"ruleset main", func() (bool, string, error) { return u.ruleset("main") }},
		{"ruleset branches", func() (bool, string, error) { return u.ruleset("branches") }},
	} {
		differs, fields, err := s.fn()
		if err != nil {
			se := fail(err)
			r.Error(s.name, 0, se.code, se.detail)
			return r.Finish(steps)
		}
		line := "unchanged " + s.name
		if differs {
			line = "changed " + s.name
			if opts.DryRun {
				line = "would change " + s.name
			}
			if fields != "" {
				line += ": " + fields
			}
		}
		r.Line(line)
	}
	return r.Finish(steps)
}

// checkRepo finds the repo and checks the admin right of the logged-in
// account. ok is false when the run stops, with code as its exit code.
func (u *run) checkRepo() (code int, ok bool) {
	repo := u.opts.Repo
	if repo == "" {
		out, err := u.gh.run("gh repo view", nil, "repo", "view", "--json", "nameWithOwner", "--jq", ".nameWithOwner")
		if err != nil {
			return u.r.Abort(".", "no-repo", "the current folder has no GitHub repo; name one with --repo <owner>/<name>"), false
		}
		repo = strings.TrimSpace(string(out))
	}
	if !repoPattern.MatchString(repo) {
		return u.r.Abort(repo, "no-repo", "a repo has the form <owner>/<name>"), false
	}
	u.repo = repo
	v, err := u.gh.api("GET", "repos/"+repo, nil)
	if isNotFound(err) {
		return u.r.Abort(repo, "no-repo", "GitHub has no such repo, or the logged-in account cannot see it"), false
	}
	if err != nil {
		u.r.Error(repo, 0, "api-failed", err.Error())
		return u.r.Finish(0), false
	}
	if admin, _ := lookup(v, "permissions", "admin").(bool); !admin {
		return u.r.Abort(repo, "not-admin", "the account logged in to gh has no admin right on the repo"), false
	}
	return 0, true
}

// checkAgent checks that the agent account is no admin of the repo.
func (u *run) checkAgent() (code int, ok bool) {
	v, err := u.gh.api("GET", "repos/"+u.repo+"/collaborators/"+u.opts.Agent+"/permission", nil)
	if isNotFound(err) {
		return 0, true
	}
	if err != nil {
		u.r.Error(u.opts.Agent, 0, "api-failed", err.Error())
		return u.r.Finish(0), false
	}
	if p, _ := lookup(v, "permission").(string); p == "admin" {
		u.r.Error(u.opts.Agent, 0, "agent-is-admin", "the agent account has the permission admin on "+u.repo+"; setup does not change it, remove it under Settings, Collaborators")
		return u.r.Finish(0), false
	}
	return 0, true
}

// mergeTargets are the merge settings, in the order of the output.
var mergeTargets = []struct {
	field string
	value any
}{
	{"allow_squash_merge", true},
	{"allow_merge_commit", false},
	{"allow_rebase_merge", false},
	{"allow_auto_merge", true},
	{"delete_branch_on_merge", true},
	{"squash_merge_commit_title", "PR_TITLE"},
	{"squash_merge_commit_message", "PR_BODY"},
}

func (u *run) mergeSettings() (bool, string, error) {
	v, err := u.gh.api("GET", "repos/"+u.repo, nil)
	if err != nil {
		return false, "", err
	}
	current, _ := v.(map[string]any)
	patch := map[string]any{}
	var diffs []string
	for _, t := range mergeTargets {
		cur, present := current[t.field]
		if present && equal(cur, t.value) {
			continue
		}
		patch[t.field] = t.value
		diffs = append(diffs, diff(t.field, cur, present, t.value))
	}
	if len(patch) == 0 {
		return false, "", nil
	}
	if !u.opts.DryRun {
		if _, err := u.gh.api("PATCH", "repos/"+u.repo, patch); err != nil {
			return false, "", err
		}
	}
	return true, strings.Join(diffs, ", "), nil
}

func (u *run) agentAccount() (bool, string, error) {
	login := u.opts.Agent
	_, err := u.gh.api("GET", "repos/"+u.repo+"/collaborators/"+login, nil)
	if err == nil {
		return false, "", nil
	}
	if !isNotFound(err) {
		return false, "", err
	}
	list, err := u.gh.list("repos/"+u.repo+"/invitations?per_page=100", "")
	if err != nil {
		return false, "", err
	}
	for _, inv := range list {
		if l, _ := lookup(inv, "invitee", "login").(string); strings.EqualFold(l, login) {
			return false, "", nil
		}
	}
	if !u.opts.DryRun {
		if _, err := u.gh.api("PUT", "repos/"+u.repo+"/collaborators/"+login, map[string]any{"permission": "push"}); err != nil {
			return false, "", err
		}
	}
	return true, "", nil
}

var branchPolicy = map[string]any{"protected_branches": false, "custom_branch_policies": true}

func (u *run) environment() (bool, string, error) {
	envPath := "repos/" + u.repo + "/environments/" + Environment
	v, err := u.gh.api("GET", envPath, nil)
	exists := err == nil
	if err != nil && !isNotFound(err) {
		return false, "", err
	}
	changed := false
	policyOK := exists && contains(lookup(v, "deployment_branch_policy"), branchPolicy)
	if !policyOK {
		changed = true
		if !u.opts.DryRun {
			if _, err := u.gh.api("PUT", envPath, map[string]any{"deployment_branch_policy": branchPolicy}); err != nil {
				return false, "", err
			}
		}
	}

	// The branch policies exist only once the environment has custom
	// branch policies; a dry run on a new environment has none yet.
	var policies []any
	if policyOK || !u.opts.DryRun {
		policies, err = u.gh.list(envPath+"/deployment-branch-policies?per_page=100", "branch_policies")
		if err != nil {
			return false, "", err
		}
	}
	hasMain := false
	for _, p := range policies {
		name, _ := lookup(p, "name").(string)
		typ, _ := lookup(p, "type").(string)
		if name == "main" && (typ == "" || typ == "branch") && !hasMain {
			hasMain = true
			continue
		}
		changed = true
		if !u.opts.DryRun {
			if _, err := u.gh.api("DELETE", envPath+"/deployment-branch-policies/"+fmt.Sprint(lookup(p, "id")), nil); err != nil {
				return false, "", err
			}
		}
	}
	if !hasMain {
		changed = true
		if !u.opts.DryRun {
			if _, err := u.gh.api("POST", envPath+"/deployment-branch-policies", map[string]any{"name": "main", "type": "branch"}); err != nil {
				return false, "", err
			}
		}
	}
	return changed, "", nil
}

func (u *run) secret() (bool, string, error) {
	_, err := u.gh.api("GET", "repos/"+u.repo+"/environments/"+Environment+"/secrets/"+SecretName, nil)
	if err == nil {
		return false, "", nil
	}
	if !isNotFound(err) {
		return false, "", err
	}
	if u.opts.DryRun {
		return true, "", nil
	}
	u.r.Line("reading " + SecretName + " from standard input")
	value := ""
	if u.opts.Stdin != nil {
		line, err := bufio.NewReader(u.opts.Stdin).ReadString('\n')
		if err != nil && err != io.EOF {
			return false, "", &stepError{code: "no-secret", detail: "standard input: " + err.Error()}
		}
		value = strings.TrimSpace(line)
	}
	if value == "" {
		return false, "", &stepError{code: "no-secret", detail: "standard input holds no value for " + SecretName}
	}
	if _, err := u.gh.run("gh secret set", []byte(value), "secret", "set", SecretName, "--env", Environment, "--repo", u.repo); err != nil {
		return false, "", err
	}
	return true, "", nil
}

// rulesetFields are the fields of a ruleset that the setup compares, in
// the order of the output.
var rulesetFields = []string{"name", "target", "enforcement", "bypass_actors", "conditions", "rules"}

// Ruleset returns the content of the embedded ruleset name.
func Ruleset(name string) (map[string]any, error) {
	data, err := rulesetFiles.ReadFile("rulesets/" + name + ".json")
	if err != nil {
		return nil, err
	}
	var v map[string]any
	dec := json.NewDecoder(strings.NewReader(string(data)))
	dec.UseNumber()
	if err := dec.Decode(&v); err != nil {
		return nil, fmt.Errorf("ruleset %s: %v", name, err)
	}
	return v, nil
}

func (u *run) ruleset(name string) (bool, string, error) {
	target, err := Ruleset(name)
	if err != nil {
		return false, "", &stepError{code: "api-failed", detail: err.Error()}
	}
	list, err := u.gh.list("repos/"+u.repo+"/rulesets?per_page=100", "")
	if err != nil {
		return false, "", err
	}
	id := ""
	for _, rs := range list {
		source, _ := lookup(rs, "source_type").(string)
		if n, _ := lookup(rs, "name").(string); n == name && (source == "" || source == "Repository") {
			id = fmt.Sprint(lookup(rs, "id"))
			break
		}
	}
	var current map[string]any
	if id != "" {
		v, err := u.gh.api("GET", "repos/"+u.repo+"/rulesets/"+id, nil)
		if err != nil {
			return false, "", err
		}
		current, _ = v.(map[string]any)
	}
	var diffs []string
	for _, f := range rulesetFields {
		cur, present := current[f]
		if present && contains(cur, target[f]) {
			continue
		}
		diffs = append(diffs, diff(f, cur, present, target[f]))
	}
	if len(diffs) == 0 {
		return false, "", nil
	}
	if !u.opts.DryRun {
		if id == "" {
			_, err = u.gh.api("POST", "repos/"+u.repo+"/rulesets", target)
		} else {
			_, err = u.gh.api("PUT", "repos/"+u.repo+"/rulesets/"+id, target)
		}
		if err != nil {
			return false, "", err
		}
	}
	return true, strings.Join(diffs, ", "), nil
}
