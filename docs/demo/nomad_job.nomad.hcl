job "example" {
  namespace = "platform"
  group "cache" {
    task "redis" {
      driver = "docker"

      config {
        image = "redis:7"
      }
    }
  }
}
