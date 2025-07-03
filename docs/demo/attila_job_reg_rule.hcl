name            = "platform_namespace"
region_contexts = ["namespace"]

region_filter {
  expression {
    selector = "any(region_namespace, {.Name == \"platform\"})"
  }
}

region_picker {
  expression {
    selector = "filter(regions, .Group == \"europe\")"
  }
}
