export function set_page_title(title) {
  document.title = title;
}

export function set_timeout(delay, fn) {
  setTimeout(fn, delay);
}

export function set_timeout_async(delay, fn) {
  setTimeout(async () => fn, delay);
}

export function is_same_node(element_one, element_two) {
  return element_one.isSameNode(element_two);
}

export function push_to_url(path) {
  document.location.replace(path);
}
