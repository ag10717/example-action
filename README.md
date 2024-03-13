# Example Action - Setup Build Number

## Current State

The current setup build number used in the github action for deploying a function app is a script written in `PowerShell` and it copied between all repositories it's requried in.

The general flow and state is:

- gets the current commit hash
- gets all tags
- gets latest tag
- splits latest tag into `major_version`, `minor_version`, `patch_version`
- gets the commit hash of the latest tagged version

- if there are no tags
  - output is the inputed `major_version`
END

- if the latest tag commit is the same as the current commit AND the branch is "main"
  - output is inputed `major_version` (if it is different from the current major)
  - output is inputed `major_version` from the latest tag
- if the latest commit message includes the word "FEATURE"
  - output is the inputed `major_version` + plus the current `minor_version` from the latest tag incremented by 1
- if the latest commit message does NOT include the word "FEATURE"
  - output is the inputed `major_version` + plus teh current `minor_version` + plus the current `patch_version` from the latest tag incremented by 1
END

- if the current source branch name includes "feature/"
  - output is the inputed `major_version` + plus the current `minor_version` from the latest tag incremented by 1
- if the current source branch does not include the word "feature/" and it is not handled from "main" 
  - output is the inputed `major_version` + plus the current `minor_version` + plus the current `patch_version` from the latest tag incremented by 1
END

version = `major_version`.`minor_version`.`patch_version`

- if the source branch is NOT "main" 
  - version = version + "PREVIEW.github_run_id" (github_run_id is the unique id of the current run)
  - output is added the runners env vars
END

## Future State

For future editions of custom scripts that produce a value used later (or provide a singular action that is re-used) we should consider creating a custom action.

This would be a docker container running whichever language we deem the best experience and would be able to be used in any repository.

### Concept

The concept I have choosen is writing the action in GOLANG. I have used this for some inspiration: https://thedevelopercafe.com/articles/custom-github-action-with-go-29d9ce66e5a8

The general flow of this action should be:

> NOTE: SUBJECT TO CHANGE

- get the latest tag
- if the source branch is "main"
  - increment the `minor_version`
- if the source branch is "main" AND it's flagged as "breaking" (in the commit message)
  - increment the `major_version`
- if the source branch is "feature/task/bugix" 
  - increment the `minor_version`
- if the source branch is any other 
  - increment the `patch_version`