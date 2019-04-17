module "child" {
    source = "./child"
}

module "child2" {
    source = "./child"
    memory = "${module.child.memory_max}"
}
