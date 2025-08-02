schema_version = 1

project {
  license        = "BUSL-1.1"
  copyright_holder = "Kopexa GmbH"
  copyright_year = 2025

  # (OPTIONAL) A list of globs that should not have copyright/license headers.
  # Supports doublestar glob patterns for more flexibility in defining which
  # files or folders should be ignored
  header_ignore = [
    # "vendor/**",
    # "**autogen**",
    "node_modules",
    "**/*.tf",
    "**/testdata/**",
    "**/*.pb.go",
    "**/*_string.go",
    "**/*.html",
    ".git/**",
    "**/*.sql",
    "build/reports/test-unit.xml",
  ]
}