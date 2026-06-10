#!/bin/bash

# ============================================
# Docker 镜像构建与推送脚本
# 用于构建前端和后端镜像并推送到阿里云镜像仓库
# ============================================

set -e

# 配置区域
REGISTRY="registry.cn-hangzhou.aliyuncs.com"
NAMESPACE="kk-mock"
WEB_IMAGE="${REGISTRY}/${NAMESPACE}/web"
SERVER_IMAGE="${REGISTRY}/${NAMESPACE}/server"

# 版本号（日期时间戳）
VERSION=$(date +"%Y%m%d%H%M")
IMAGE_TAG="${VERSION}"

# 项目根目录
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查 Docker 是否运行
check_docker() {
    if ! docker info > /dev/null 2>&1; then
        log_error "Docker 未运行，请先启动 Docker"
        exit 1
    fi
    log_info "Docker 检查通过"
}

# 登录阿里云镜像仓库
login_registry() {
    log_info "正在登录阿里云镜像仓库..."
    docker login --username=yangfans ${REGISTRY}
    log_info "登录成功"
}

# 构建前端镜像
build_web() {
    log_info "开始构建前端镜像..."

    WEB_DIR="${PROJECT_ROOT}/web"
    if [ ! -d "${WEB_DIR}" ]; then
        log_error "前端目录不存在: ${WEB_DIR}"
        exit 1
    fi

    docker build -t ${WEB_IMAGE}:${IMAGE_TAG} \
                 -f ${WEB_DIR}/Dockerfile \
                 ${WEB_DIR}

    log_info "前端镜像构建完成: ${WEB_IMAGE}:${IMAGE_TAG}"
}

# 构建后端镜像
build_server() {
    log_info "开始构建后端镜像..."

    SERVER_DIR="${PROJECT_ROOT}/server"
    if [ ! -d "${SERVER_DIR}" ]; then
        log_error "后端目录不存在: ${SERVER_DIR}"
        exit 1
    fi

    docker build -t ${SERVER_IMAGE}:${IMAGE_TAG} \
                 -f ${SERVER_DIR}/Dockerfile \
                 ${SERVER_DIR}

    log_info "后端镜像构建完成: ${SERVER_IMAGE}:${IMAGE_TAG}"
}

# 推送镜像到阿里云
push_images() {
    log_info "开始推送镜像到阿里云..."

    log_info "推送前端镜像: ${WEB_IMAGE}:${IMAGE_TAG}"
    docker push ${WEB_IMAGE}:${IMAGE_TAG}

    log_info "推送后端镜像: ${SERVER_IMAGE}:${IMAGE_TAG}"
    docker push ${SERVER_IMAGE}:${IMAGE_TAG}

    log_info "镜像推送完成"
}

# 标记 latest 标签
tag_latest() {
    log_info "标记 latest 标签..."

    docker tag ${WEB_IMAGE}:${IMAGE_TAG} ${WEB_IMAGE}:latest
    docker tag ${SERVER_IMAGE}:${IMAGE_TAG} ${SERVER_IMAGE}:latest

    log_info "推送 latest 标签..."
    docker push ${WEB_IMAGE}:latest
    docker push ${SERVER_IMAGE}:latest
}

# 显示帮助信息
show_help() {
    echo "Docker 镜像构建与推送脚本"
    echo ""
    echo "用法: $0 [选项]"
    echo ""
    echo "选项:"
    echo "  --no-login    跳过登录步骤（如果已经登录）"
    echo "  --no-latest    不标记并推送 latest 标签"
    echo "  --help         显示帮助信息"
    echo ""
    echo "示例:"
    echo "  $0              # 完整构建和推送流程"
    echo "  $0 --no-login   # 跳过登录（已登录时使用）"
    echo "  $0 --no-latest  # 不推送 latest 标签"
    echo ""
}

# 主函数
main() {
    echo "========================================"
    echo "  Docker 镜像构建与推送"
    echo "========================================"
    echo ""
    echo "配置信息:"
    echo "  镜像仓库: ${REGISTRY}"
    echo "  命名空间: ${NAMESPACE}"
    echo "  版本号: ${IMAGE_TAG}"
    echo "  前端镜像: ${WEB_IMAGE}:${IMAGE_TAG}"
    echo "  后端镜像: ${SERVER_IMAGE}:${IMAGE_TAG}"
    echo ""

    # 解析参数
    SKIP_LOGIN=false
    SKIP_LATEST=false

    while [[ $# -gt 0 ]]; do
        case $1 in
            --no-login)
                SKIP_LOGIN=true
                shift
                ;;
            --no-latest)
                SKIP_LATEST=true
                shift
                ;;
            --help)
                show_help
                exit 0
                ;;
            *)
                log_error "未知参数: $1"
                show_help
                exit 1
                ;;
        esac
    done

    # 执行构建流程
    check_docker

    if [ "$SKIP_LOGIN" = false ]; then
        login_registry
    else
        log_info "跳过登录步骤"
    fi

    build_web
    build_server
    push_images

    if [ "$SKIP_LATEST" = false ]; then
        tag_latest
    else
        log_info "跳过 latest 标签"
    fi

    echo ""
    echo "========================================"
    echo "  构建与推送完成!"
    echo "========================================"
    echo ""
    echo "镜像信息:"
    echo "  前端: ${WEB_IMAGE}:${IMAGE_TAG}"
    echo "  后端: ${SERVER_IMAGE}:${IMAGE_TAG}"
    echo ""
    echo "下一步:"
    echo "  1. 将 deploy/docker-compose-remote.yaml 复制到服务器"
    echo "  2. 在服务器上运行: docker-compose -f docker-compose-remote.yaml up -d"
    echo ""
}

main "$@"
