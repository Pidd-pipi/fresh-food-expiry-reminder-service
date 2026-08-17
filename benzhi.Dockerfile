# 评测用镜像：从仓库根目录构建 backend 下的 Go 项目。
FROM golang:1.22
WORKDIR /app
COPY backend/go.mod backend/go.sum ./backend/
RUN cd backend && go mod download
COPY . .
RUN cd backend && go build ./...
CMD ["bash"]

# 多架构交叉构建示例（如需交付双架构镜像）：
# docker buildx build --platform linux/arm64,linux/amd64 -f benzhi.Dockerfile -t <image> .
