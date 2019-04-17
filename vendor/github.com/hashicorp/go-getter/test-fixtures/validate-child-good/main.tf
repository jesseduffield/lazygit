module "child" {
    source = "./child"
    memory = "1G"
}

resource "aws_instance" "foo" {
    memory = "${module.child.result}"
}
