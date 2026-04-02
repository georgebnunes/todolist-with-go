# ── Variáveis ────────────────────────────────────────────────
variable "aws_region" {
  default = "us-east-1"
}

variable "lambda_zip_path" {
  description = "Caminho para o zip do Lambda gerado pelo make zip"
  default     = "../lambda.zip"
}