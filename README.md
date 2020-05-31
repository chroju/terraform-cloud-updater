terraform-cloud-updater
======

This action automates the terraform version update of your Terraform cloud workspace, based on the remote config file.

## Usage

This action parses the remote config file and checks which Terraform version is configured in your Terraform cloud workspace. If the latest Terraform version is not set, it will notice. It is also possible to automatically update the latest version.

For example, if you use `. /terraform` directory, including the following remote config file.

```hcl
terraform {
  backend "remote" {
    hostname = "app.terraform.io"
    organization = "sample"

    workspaces {
      name = "sample"
    }
  }
  required_version = "> 0.12.0, <= 0.12.23"
}
```

GitHub Actions are configuread as follows.

```yaml
on:
  push:
  pull_request:
    branches:
      - master

jobs:
  tf_cloud_update:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: check update
        id: check
        uses: chroju/terraform-cloud-updater@v1
        with:
          working_dir: ./terraform
          comment_pr: true
        env:
          TFE_TOKEN: ${{ secrets.TFE_TOKEN }}
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
      - name: result
        run: echo "${{ steps.check.outputs.output }}"
        if: "${{ steps.check.outputs.is_available_update == 'true' }}"
```

When action finds a new version, it will comment on the pull request. Action will also check for consistency with the required version and will not automatically update to a non-conforming version (In this case, it will not be updated to `0.12.24` ). Sample pull request is [here](https://github.com/chroju/terraform-cloud-updater/pull/17).

## Inputs

* `working_dir` - (Optional) Terraform working directory. Defaults to `./` (root of the GitHub repository) .
* `auto_update` - (Optional) Not only notice, automatically update Terraform Cloud workspace to the latest version compatible with required version. Defaults to `false` .
* `comment_pr` - (Optional) Whether or not to post a comment on GitHub pull requests. If you set it to true, you need to set the `GITHUB_TOKEN` environment variable. Defaults to `false` .

## Outputs

* `is_available_update` - Wherther or not available update ('true' or 'false') .
* `output` - result message (empty if new version is not released) .

## Environment Variables

* `TFE_TOKEN` - (Required) Terraform Cloud API token.
* `GITHUB_TOKEN` -  (Optional) The GitHub API token used to post comments to pull requests. Not required if the `comment_pr` input is set to `false` .

## Notes

### Support for Terraform Enterprise

Since Terraform Cloud and Terraform Enterprise have a common API, this action is likely to be available for Terraform Enterprise as well. However, we have not verified the operation using Terraform Enterprise.
