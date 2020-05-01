load("system", "GO")
load("std/golang", "go_module")

neon = go_module(
    go = GO,
    name = "neon",
    go_mod = glob("go.mod"),
    go_sum = glob("go.sum"),
    sources = glob("**/*.go"),
)