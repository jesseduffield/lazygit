module "child" {
    source = "./child"
}

resource "aws_instance" "foo" {
    memory = "${module.child.memory}"
}
