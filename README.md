<div align="center">
  <h1>Terraform Provider for Typesense</h1>
  <strong>This is a Terraform provider for Typesense</strong>
</div>

<hr>

## Support

- Supports v28.0+ version of Typesense.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= v0.12.0 (v0.11.x may work but not supported actively)

## Building The Provider

Clone repository to: `$GOPATH/src/github.com/ronati/terraform-provider-typesense`

```console
$ mkdir -p $GOPATH/src/github.com/ronati; cd $GOPATH/src/github.com/ronati
$ git clone git@github.com:ronati/terraform-provider-typesense
Enter the provider directory and build the provider

$ cd $GOPATH/src/github.com/ronati/terraform-provider-typesense
$ make build
```

# Contribute

**All contributions are welcome!**

## Commit format

This project is setup with automatic semver versioning based on your commit semantic. It uses [`commitizen`](https://commitizen.github.io/cz-cli/) to enforce the format and help contributors format their commit message. We follow the [conventional commit format](https://www.conventionalcommits.org/en/v1.0.0/). Once you want to commit your work, you need to:

## Notes for project's maintainers

When you merge a PR from `beta` into `main` and it successfully published a new version on the `latest` channel, **don't forget to create a PR from `main` to `beta`**. **This is mandatory** for `semantic-release` to take it into account for next `beta` version.
