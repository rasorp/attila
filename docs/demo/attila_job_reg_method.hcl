name      = "europe-platform"
namespace = "default"

selector "namespace_platform" {
  provider = "filter"

  config {
   expression = "job.Namespace == \"platform\""
  }
}

rule {
  name = "europe-platform"
}
