resource "misskey_note" "test" {
  text             = <<-EOT
        テスト1
        テスト2
        テスト3
    EOT
  visibility       = "specified"
  visible_user_ids = []
}
