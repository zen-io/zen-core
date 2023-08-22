import { Target } from "../targets"

export type RunRequest = {
  target: Target
  options: RunOptions
}

export type RunOptions = {
  dryDrun: boolean
  deployEnvironment: string
  useEnvironments: boolean
}

export type WorkerConfig = {
  modulePath: string
  knownTargetTypes: { [key: string]: string }
}
