# terraform-provider-misskey

During local development, place the following content in `~/.terraformrc`.
It enables you to load the Terraform provider from the local directory specified.

```
provider_installation {
  dev_overrides {
    "maeda6uiui/misskey" = "/home/maeda6uiui/WorkingDir/terraform-provider-misskey/example/bin/"
  }

  direct {}
```
