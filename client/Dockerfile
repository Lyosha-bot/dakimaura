FROM nginx

COPY nginx.conf /etc/nginx/conf.d/default.conf

COPY . /usr/share/nginx/html

# ENTRYPOINT [ "nginx", "-g", "daemon off;" ]