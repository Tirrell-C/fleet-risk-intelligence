#!/bin/bash
set -e

echo "🚀 Setting up Fleet Risk Intelligence MVP development environment..."

# Install Go tools
echo "📦 Installing Go tools..."
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/swaggo/swag/cmd/swag@latest
go install github.com/99designs/gqlgen@latest

# Install Pulumi
echo "☁️  Installing Pulumi..."
curl -fsSL https://get.pulumi.com | sh
echo 'export PATH=$PATH:$HOME/.pulumi/bin' >> ~/.bashrc

# Install Node.js dependencies globally
echo "🌐 Installing Node.js tools..."
npm install -g yarn pnpm typescript @types/node

# Install Docker Compose
echo "🐳 Installing Docker Compose..."
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Install additional development tools
echo "🔧 Installing additional tools..."
sudo apt-get update
sudo apt-get install -y jq make postgresql-client redis-tools

# Set up Git hooks directory
mkdir -p .git/hooks

echo "✅ Development environment setup complete!"
echo "🎯 Ready to build the Fleet Risk Intelligence MVP!"