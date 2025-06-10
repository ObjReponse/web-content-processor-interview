#!/bin/bash

# 🚀 Setup script for interview repository
# Запустите этот скрипт после создания репозитория из template

set -e

echo "🔧 Setting up interview repository..."

# Проверяем что не инициализировано уже
if [ -f .template-initialized ]; then
    echo "⚠️  Repository already initialized"
    exit 0
fi

# Получаем URL репозитория из git remote
REPO_URL=$(git remote get-url origin)
echo "🔍 Repository URL: $REPO_URL"

if [[ $REPO_URL == *"github.com"* ]]; then
    # Извлекаем путь репозитория из URL
    # Поддерживаем и HTTPS и SSH URLs
    if [[ $REPO_URL == git@github.com:* ]]; then
        # SSH format: git@github.com:user/repo.git
        REPO_PATH=$(echo $REPO_URL | sed 's/git@github\.com://; s/\.git$//')
    else
        # HTTPS format: https://github.com/user/repo.git
        REPO_PATH=$(echo $REPO_URL | sed 's/.*github\.com\///; s/\.git$//')
    fi

    echo "📍 Repository path: $REPO_PATH"
else
    echo "❌ Cannot determine repository path from: $REPO_URL"
    echo "💡 Make sure origin remote points to GitHub"
    exit 1
fi

# Проверяем что это не template репозиторий
if [[ $REPO_PATH == *"web-content-processor-interview" ]] && [[ $REPO_PATH == *"your-company"* ]]; then
    echo "⚠️  This appears to be the template repository"
    echo "💡 Please create a new repository from this template first"
    exit 1
fi

# Обновляем go.mod
echo "📝 Updating go.mod..."
sed -i.bak "s|github.com/your-company/web-content-processor-interview|github.com/$REPO_PATH|g" go.mod

# Обновляем импорты в main.go
echo "📝 Updating imports in main.go..."
sed -i.bak "s|github.com/your-company/web-content-processor-interview/processor|github.com/$REPO_PATH/processor|g" main.go

# Удаляем backup файлы
rm -f go.mod.bak main.go.bak

# Создаем файл-маркер
echo "Repository initialized on $(date)" > .template-initialized

# Обновляем go.sum
echo "📦 Updating dependencies..."
go mod tidy

echo "✅ Setup completed successfully!"
echo "🏃 You can now run: make test"

# Удаляем сам setup script
echo "🧹 Cleaning up setup script..."
rm -f setup.sh