FROM node:24-alpine

WORKDIR /app

# Copy package manifests first for caching
COPY web/website/package*.json ./

# Install dependencies (container keeps node_modules)
RUN npm ci

# Copy source (will be overridden by volume during dev)
COPY web/website ./

# Expose Vite dev port
EXPOSE 5173

# Start dev server and listen on all interfaces
CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0", "--port", "5173"]
