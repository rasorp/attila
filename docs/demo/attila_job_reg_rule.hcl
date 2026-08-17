name = "europe-platform"

region_context { kind = "namespace" }

region_picker "europe-region" {
  provider = "expr"

  config {
    expression = "regions.filter(r, r.group == \"europe\")"
  }
}

region_picker "platform-namespace" {
  provider = "filter"

  config {
    expression = "any(region_namespace, {.Name == \"platform\"})"
  }
}
