package pikoci

// ReleaseTemplate es el HCL que ggit publica para crear releases al empujar un
// tag. Es independiente de PipelineTemplate: se activa en un pipeline aparte.
//
// Variables esperadas al publicar el pipeline:
//
//	ggit_url    URL base de ggit
//	ggit_token  token del usuario bot
//	repo        owner/name del repositorio
//	git_url     URL de clone del repositorio
//	target      rama sobre la que se crea el release
const ReleaseTemplate = `# Pipeline de releases gestionado por ggit.
variable "ggit_url" {
  type    = string
  default = ""
}
variable "ggit_token" {
  type    = string
  default = ""
}
variable "repo" {
  type    = string
  default = ""
}
variable "git_url" {
  type    = string
  default = ""
}
variable "serial_group" {
  type    = string
  default = "ggit-ci"
}
variable "target" {
  type    = string
  default = "main"
}

resource "git" "repo" {
  params {
    url   = var.git_url
    name  = "repo"
    token = var.ggit_token
    tag   = "true"
  }
}

resource_type "ggit-release" {
  params = ["ggit_url", "ggit_token", "repo", "target"]

  push "exec" {
    path = "/bin/sh"
    args = [
      "-ec",
      <<-EOT
      TAG=$(git -C repo describe --tags --exact-match)
      BODY=$(cat notes.md 2>/dev/null || echo "Release $TAG")

      jq -n --arg tag "$TAG" --arg target "$param_target" --arg body "$BODY" \
        '{tag_name:$tag,target:$target,name:$tag,body:$body}' \
      | curl -sf -X POST "$param_ggit_url/api/v2/$param_repo/releases" \
          -H "Authorization: token $param_ggit_token" \
          -H "Content-Type: application/json" -d @-
      EOT
    ]
  }
}

resource "ggit-release" "release" {
  params {
    ggit_url   = var.ggit_url
    ggit_token = var.ggit_token
    repo       = var.repo
    target     = var.target
  }
}

job "release" {
  concurrency   = 1
  serial_groups = [var.serial_group]
  interruptible = true

  get "git" "repo" { trigger = true }

  task "notes" {
    timeout = "5m"
    run "shell" { cmd = "cd repo && git log -1 --pretty=%B > ../notes.md" }
  }

  put "ggit-release" "release" {}
}
`
