output {
  one = "${replace(var.sub_domain, ".", "\\.")}"
  two = "${replace(var.sub_domain, ".", "\\\\.")}"
  many = "${replace(var.sub_domain, ".", "\\\\\\\\.")}"
}
