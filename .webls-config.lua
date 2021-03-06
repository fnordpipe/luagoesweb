return {
  -- generic website options
  title       = "gorillua",
  description = "lua application server",
  website     = "gorillua.fnordpipe.org",
  cname       = "gorillua.fnordpipe.org",

  -- add "index.html" to the end of each link
  pagesuffix  = nil,

  -- folder that shall be scanned for content
  scanpath    = "content",

  -- the output directory that is later used as webroot
  www         = "www",

  -- modules that should be used for parsing
  modules     = { "html", "markdown", "gallery", "download", "footer" },

  -- define default colors
  colors      = {
    ["accent"]      = "#0269a4",
    ["border"]      = "#eeeeee",
    ["bg-page"]     = "#eeeeee",
    ["bg-content"]  = "#ffffff",
    ["bg-sidebar"]  = "#2e2d2b",
    ["fg-page"]     = "#000000",
    ["fg-sidebar"]  = "#ffffff",
    ["customcss"]   = "";
  },
}
