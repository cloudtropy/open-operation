package main

import (
  "cloudtropy.com/alert/g"
  "cloudtropy.com/alert/router"
)

func main() {

  g.ParseConfig("./configure.json")

  router.RunServer()
}
