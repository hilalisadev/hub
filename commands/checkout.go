package commands

import (
	"fmt"
	"regexp"

	"github.com/hilalisadev/hub/git"
	"github.com/hilalisadev/hub/github"
	"github.com/hilalisadev/hub/utils"
)

var cmdCheckout = &Command{
	Run:          checkout,
	GitExtension: true,
	Usage:        "checkout <PULLREQ-URL> [<BRANCH>]",
	Long: `Check out the head of a pull request as a local branch.

## Examples:
		$ hub checkout https://github.com/jingweno/gh/pull/73
		> git fetch origin pull/73/head:jingweno-feature
		> git checkout jingweno-feature

## See also:

hub-merge(1), hub-am(1), hub(1), git-checkout(1)
`,
}

func init() {
	CmdRunner.Use(cmdCheckout)
}

func checkout(command *Command, args *Args) []byte {
	words := args.Words()

	if len(words) == 0 {
		return nil
	}

	checkoutURL := words[0]
	var newBranchName string
	if len(words) > 1 {
		newBranchName = words[1]
	}

	url, err := github.ParseURL(checkoutURL)
	if err != nil {
		// not a valid GitHub URL
		return nil
	}

	pullURLRegex := regexp.MustCompile("^pull/(\\d+)")
	projectPath := url.ProjectPath()
	if !pullURLRegex.MatchString(projectPath) {
		// not a valid PR URL
		return nil
	}

	err = sanitizeCheckoutFlags(args)
	utils.Check(err)

	id := pullURLRegex.FindStringSubmatch(projectPath)[1]
	gh := github.NewClient(url.Project.Host)
	pullRequest, err := gh.PullRequest(url.Project, id)
	utils.Check(err)

	newArgs, err := transformCheckoutArgs(args, pullRequest, newBranchName)
	utils.Check(err)

	if idx := args.IndexOfParam(newBranchName); idx >= 0 {
		args.RemoveParam(idx)
	}
	replaceCheckoutParam(args, checkoutURL, newArgs...)
	return nil
}

func transformCheckoutArgs(args *Args, pullRequest *github.PullRequest, newBranchName string) (newArgs []string, err error) {
	repo, err := github.LocalRepo()
	if err != nil {
		return
	}

	baseRemote, err := repo.RemoteForRepo(pullRequest.Base.Repo)
	if err != nil {
		return
	}

	var headRemote *github.Remote
	if pullRequest.IsSameRepo() {
		headRemote = baseRemote
	} else if pullRequest.Head.Repo != nil {
		headRemote, _ = repo.RemoteForRepo(pullRequest.Head.Repo)
	}

	if headRemote != nil {
		if newBranchName == "" {
			newBranchName = pullRequest.Head.Ref
		}
		remoteBranch := fmt.Sprintf("%s/%s", headRemote.Name, pullRequest.Head.Ref)
		refSpec := fmt.Sprintf("+refs/heads/%s:refs/remotes/%s", pullRequest.Head.Ref, remoteBranch)
		if git.HasFile("refs", "heads", newBranchName) {
			newArgs = append(newArgs, newBranchName)
			args.After("git", "merge", "--ff-only", fmt.Sprintf("refs/remotes/%s", remoteBranch))
		} else {
			newArgs = append(newArgs, "-b", newBranchName, "--no-track", remoteBranch)
			args.After("git", "config", fmt.Sprintf("branch.%s.remote", newBranchName), headRemote.Name)
			args.After("git", "config", fmt.Sprintf("branch.%s.merge", newBranchName), "refs/heads/"+pullRequest.Head.Ref)
		}
		args.Before("git", "fetch", headRemote.Name, refSpec)
	} else {
		if newBranchName == "" {
			newBranchName = pullRequest.Head.Ref
			if pullRequest.Head.Repo != nil && newBranchName == pullRequest.Head.Repo.DefaultBranch {
				newBranchName = fmt.Sprintf("%s-%s", pullRequest.Head.Repo.Owner.Login, newBranchName)
			}
		}
		newArgs = append(newArgs, newBranchName)

		b, errB := repo.CurrentBranch()
		isCurrentBranch := errB == nil && b.ShortName() == newBranchName

		ref := fmt.Sprintf("refs/pull/%d/head", pullRequest.Number)
		if isCurrentBranch {
			args.Before("git", "fetch", baseRemote.Name, ref)
			args.After("git", "merge", "--ff-only", "FETCH_HEAD")
		} else {
			args.Before("git", "fetch", baseRemote.Name, fmt.Sprintf("%s:%s", ref, newBranchName))
		}

		remote := baseRemote.Name
		mergeRef := ref
		if pullRequest.MaintainerCanModify && pullRequest.Head.Repo != nil {
			var project *github.Project
			project, err = github.NewProjectFromRepo(pullRequest.Head.Repo)
			if err != nil {
				return
			}

			remote = project.GitURL("", "", true)
			mergeRef = fmt.Sprintf("refs/heads/%s", pullRequest.Head.Ref)
		}

		if mc, err := git.Config(fmt.Sprintf("branch.%s.merge", newBranchName)); err != nil || mc == "" {
			args.After("git", "config", fmt.Sprintf("branch.%s.remote", newBranchName), remote)
			args.After("git", "config", fmt.Sprintf("branch.%s.merge", newBranchName), mergeRef)
		}
	}
	return
}

func sanitizeCheckoutFlags(args *Args) error {
	if i := args.IndexOfParam("-b"); i != -1 {
		return fmt.Errorf("Unsupported flag -b when checking out pull request")
	}

	if i := args.IndexOfParam("--orphan"); i != -1 {
		return fmt.Errorf("Unsupported flag --orphan when checking out pull request")
	}

	return nil
}

func replaceCheckoutParam(args *Args, checkoutURL string, replacement ...string) {
	idx := args.IndexOfParam(checkoutURL)
	args.RemoveParam(idx)
	args.InsertParam(idx, replacement...)
}
