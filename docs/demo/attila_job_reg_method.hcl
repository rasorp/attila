name = "europe-platform"

selector "namespace_platform" {
  provider = "filter"

  config {
   expression = "job.Namespace == \"platform\""
  }
}

rule {
  name = "europe-platform"
}
