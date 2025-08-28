#!/bin/bash

echo "🚀 开始部署前端到后端服务器..."

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 检查是否在正确的目录
if [ ! -f "backend/main.go" ]; then
    echo -e "${RED}❌ 错误：请在项目根目录运行此脚本${NC}"
    exit 1
fi

# 步骤1：构建前端
echo -e "${YELLOW}📦 步骤1：构建前端生产版本...${NC}"
cd frontend/res_pro

if [ ! -f "package.json" ]; then
    echo -e "${RED}❌ 错误：前端目录不存在${NC}"
    exit 1
fi

# 安装依赖（如果需要）
if [ ! -d "node_modules" ]; then
    echo "📥 安装依赖..."
    npm install
fi

# 构建项目
echo "🔨 构建项目..."
npm run build:h5

if [ $? -ne 0 ]; then
    echo -e "${RED}❌ 前端构建失败${NC}"
    exit 1
fi

echo -e "${GREEN}✅ 前端构建成功！${NC}"

# 步骤2：创建后端静态文件目录
echo -e "${YELLOW}📁 步骤2：准备后端静态文件目录...${NC}"
cd ../../backend

# 创建static目录
mkdir -p static

# 复制前端构建文件到后端static目录
echo "📋 复制前端文件到后端..."
cp -r ../frontend/res_pro/dist/build/h5/* static/

# 重命名index.html（如果需要）
if [ -f "static/index.html" ]; then
    echo "✅ 找到index.html文件"
else
    echo -e "${RED}❌ 未找到index.html文件，请检查构建输出${NC}"
    exit 1
fi

# 步骤3：重启后端服务
echo -e "${YELLOW}🔄 步骤3：重启后端服务...${NC}"

# 查找并停止现有的后端进程
echo "🛑 停止现有后端服务..."
pkill -f "reschedule-program" || true

# 等待进程完全停止
sleep 2

# 重新构建并启动后端
echo "🔨 重新构建后端..."
go build -o server .

if [ $? -ne 0 ]; then
    echo -e "${RED}❌ 后端构建失败${NC}"
    exit 1
fi

# 启动后端服务
echo "🚀 启动后端服务..."
nohup ./server > server.log 2>&1 &

# 等待服务启动
sleep 3

# 检查服务是否启动成功
if pgrep -f "server" > /dev/null; then
    echo -e "${GREEN}✅ 后端服务启动成功！${NC}"
    echo -e "${GREEN}🌐 前端现在可以通过 http://47.108.201.156 访问${NC}"
    echo -e "${GREEN}📊 后端日志：tail -f backend/server.log${NC}"
else
    echo -e "${RED}❌ 后端服务启动失败，请检查日志${NC}"
    exit 1
fi

echo -e "${GREEN}🎉 部署完成！${NC}"
echo ""
echo "📋 部署总结："
echo "   ✅ 前端已构建并部署到后端"
echo "   ✅ 后端已重启并配置静态文件服务"
echo "   ✅ 前端可通过 http://47.108.201.156 访问"
echo ""
echo "🔧 常用命令："
echo "   查看后端日志：tail -f backend/server.log"
echo "   停止后端：pkill -f server"
echo "   重启后端：./deploy.sh"
