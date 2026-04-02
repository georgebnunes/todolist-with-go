# ── Outputs ──────────────────────────────────────────────────
output "api_url" {
  description = "URL base da API"
  value       = aws_apigatewayv2_stage.default.invoke_url
}

output "dynamodb_table_name" {
  value = aws_dynamodb_table.todos.name
}
