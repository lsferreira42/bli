FROM nginx:alpine

# Copiar arquivos estáticos para o diretório padrão do nginx
COPY index.html /usr/share/nginx/html/
COPY brainfuck.js /usr/share/nginx/html/
COPY brainfuck.js.map /usr/share/nginx/html/

# Configurar o nginx para escutar na porta 9093
RUN sed -i 's/listen\s*80;/listen 9093;/g' /etc/nginx/conf.d/default.conf

# Expor a porta 9093
EXPOSE 9093

# Iniciar o nginx
CMD ["nginx", "-g", "daemon off;"] 