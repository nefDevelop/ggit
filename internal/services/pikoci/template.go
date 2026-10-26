package pikoci

// PipelineTemplate es el HCL que ggit publica en PikoCI. Reporta el estado del
// build a la API v2 de ggit con el token del usuario bot.
//
// Variables esperadas al publicar el pipeline:
//
//	ggit_url    URL base de ggit
//	ggit_token  token del usuario bot
//	repo        owner/name del repositorio
//	git_url     URL de clone del repositorio
//	branch      rama a construir
//	pikoci_url  URL pública de PikoCI, para enlazar el build
const PipelineTemplate = `# Pipeline de CI gestionado por ggit.
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
variable "branch" {
  type    = string
  default = "main"
}
variable "pikoci_url" {
  type    = string
  default = ""
}

resource "git" "repo" {
  params {
    url    = var.git_url
    branch = var.branch
    name   = "repo"
    token  = var.ggit_token
  }
}

notification_type "ggit-status" {
  params = ["ggit_url", "ggit_token", "repo", "pikoci_url"]

  notify "exec" {
    path = "/bin/sh"
    args = [
      "-ec",
      <<-EOT
      SHA="$${notify_sha:-$GET_REPO_REF}"

      if [ -n "$BUILD_NUMBER" ]; then
        TARGET="$param_pikoci_url/teams/$BUILD_TEAM_NAME/pipelines/$BUILD_PIPELINE_NAME/jobs/$BUILD_JOB_NAME/builds/$BUILD_NUMBER"
      else
        TARGET="$param_pikoci_url"
      fi

      curl -sf -X POST "$param_ggit_url/api/v2/$param_repo/statuses/$SHA" \
        -H "Authorization: token $param_ggit_token" \
        -H "Content-Type: application/json" \
        -d "{\"state\":\"$notify_state\",\"context\":\"ci/pikoci\",\"description\":\"$notify_description\",\"target_url\":\"$TARGET\"}"
      EOT
    ]
  }
}

notification "ggit-status" "ci" {
  params {
    ggit_url   = var.ggit_url
    ggit_token = var.ggit_token
    repo       = var.repo
    pikoci_url = var.pikoci_url
  }
}

job "test" {
  get "git" "repo" { trigger = true }

  notify "ggit-status" "ci" {
    state       = "pending"
    description = "Pipeline en marcha"
  }

  task "test" {
    run "shell" { cmd = "cd repo && make test" }
  }

  on_success {
    notify "ggit-status" "ci" {
      state       = "success"
      description = "Pipeline correcto"
    }
  }

  on_failure {
    notify "ggit-status" "ci" {
      state       = "failure"
      description = "Pipeline falló"
    }
  }
}
`
