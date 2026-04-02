# ============================================================
# Terraform — Todo Lambda + DynamoDB + API Gateway
# ============================================================

terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = var.aws_region
}

# ── DynamoDB Table ───────────────────────────────────────────
resource "aws_dynamodb_table" "todos" {
  name         = "todos"
  billing_mode = "PAY_PER_REQUEST" # On-demand — sem provisionar capacidade
  hash_key     = "id"

  attribute {
    name = "id"
    type = "S" # String
  }

  tags = {
    Project = "todo-lambda"
  }
}

# ── IAM Role para o Lambda ───────────────────────────────────
resource "aws_iam_role" "lambda_exec" {
  name = "todo-lambda-exec-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action    = "sts:AssumeRole"
      Effect    = "Allow"
      Principal = { Service = "lambda.amazonaws.com" }
    }]
  })
}


# Permissão para o Lambda escrever logs no CloudWatch
resource "aws_iam_role_policy_attachment" "lambda_basic" {
  role       = aws_iam_role.lambda_exec.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

# Permissão para o Lambda acessar o DynamoDB
resource "aws_iam_role_policy" "lambda_dynamodb" {
  name = "todo-lambda-dynamodb-policy"
  role = aws_iam_role.lambda_exec.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Action = [
        "dynamodb:PutItem",
        "dynamodb:GetItem",
        "dynamodb:UpdateItem",
        "dynamodb:DeleteItem",
        "dynamodb:Scan",
        "dynamodb:Query"
      ]
      Resource = aws_dynamodb_table.todos.arn
    }]
  })
}

# ── Lambda Function ──────────────────────────────────────────
resource "aws_lambda_function" "todo" {
  function_name = "todo-lambda"
  role          = aws_iam_role.lambda_exec.arn
  handler       = "bootstrap"         # nome do binário Go
  runtime       = "provided.al2023"   # runtime customizado para Go
  filename      = var.lambda_zip_path

  source_code_hash = filebase64sha256(var.lambda_zip_path)

  environment {
    variables = {
      DYNAMODB_TABLE_NAME = aws_dynamodb_table.todos.name
    }
  }

  tags = {
    Project = "todo-lambda"
  }
}

# ── API Gateway (HTTP API — mais barato que REST API) ────────
resource "aws_apigatewayv2_api" "todo_api" {
  name          = "todo-api"
  protocol_type = "HTTP"
}

resource "aws_apigatewayv2_integration" "lambda" {
  api_id                 = aws_apigatewayv2_api.todo_api.id
  integration_type       = "AWS_PROXY"
  integration_uri        = aws_lambda_function.todo.invoke_arn
  payload_format_version = "2.0"
}

# Rotas — $default captura tudo e repassa ao Lambda
resource "aws_apigatewayv2_route" "todos_collection" {
  api_id    = aws_apigatewayv2_api.todo_api.id
  route_key = "POST /todos"
  target    = "integrations/${aws_apigatewayv2_integration.lambda.id}"
}

resource "aws_apigatewayv2_route" "todos_item" {
  api_id    = aws_apigatewayv2_api.todo_api.id
  route_key = "ANY /todos/{id}"
  target    = "integrations/${aws_apigatewayv2_integration.lambda.id}"
}

resource "aws_apigatewayv2_stage" "default" {
  api_id      = aws_apigatewayv2_api.todo_api.id
  name        = "$default"
  auto_deploy = true
}

# Permissão para o API Gateway invocar o Lambda
resource "aws_lambda_permission" "apigw" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.todo.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.todo_api.execution_arn}/*/*"
}