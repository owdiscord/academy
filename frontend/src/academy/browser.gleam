import gleam/dynamic/decode

@external(javascript, "./browser_ffi.mjs", "set_page_title")
pub fn set_page_title(_title: String) -> Nil {
  panic as "target not supported"
}

@external(javascript, "./browser_ffi.mjs", "set_timeout")
pub fn set_timeout(_delay: Int, _callback: fn() -> anything) -> anything {
  panic as "target not supported"
}

@external(javascript, "./browser_ffi.mjs", "is_same_node")
pub fn is_same_node(_left: decode.Dynamic, _right: decode.Dynamic) -> Bool {
  panic as "target not supported"
}

@external(javascript, "./browser_ffi.mjs", "push_to_url")
pub fn push_to_url(_url: String) -> Nil {
  panic as "target not supported"
}
