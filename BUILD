
load("std/golang", "go_module")
load("system", "go")
load("system/stdenv", "runInStdenv")

neon = go_module(
    go = go,
    runInStdenv = runInStdenv,
    name = "neon",
    go_mod = glob("go.mod"),
    go_sum = glob("go.sum"),
    sources = glob("**/*.go"),
)
