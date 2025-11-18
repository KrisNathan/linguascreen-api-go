# Use the official Bun image as base
FROM oven/bun:latest

# Set working directory
WORKDIR /app

# Copy package files
COPY package.json bun.lockb* ./

# Install dependencies
RUN bun install --frozen-lockfile

# Copy source code
COPY . .

# Pre-download Tesseract.js assets for English and Japanese
RUN bun -e "\
import { createWorker } from 'tesseract.js'; \
(async () => { \
  console.log('Pre-downloading language data...'); \
  const worker = await createWorker(['eng', 'vie', 'jpn', 'kor', 'chi_sim', 'chi_tra', 'jav']); \
  await worker.terminate(); \
  console.log('Pre-downloading Japanese language data...'); \
  console.log('Language data pre-download complete'); \
})();"

# Expose the port (Hono default is 3000)
EXPOSE 3000

# Run the application
CMD ["bun", "run", "src/index.ts"]