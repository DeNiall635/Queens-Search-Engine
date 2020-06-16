FROM ubuntu:18.04
COPY ./ /app
WORKDIR /app
RUN apt-get update
RUN apt-get install -y python3-pip
RUN pip3 install requests
RUN pip3 install flask
RUN pip3 install pymysql
EXPOSE 5000
CMD python3 ./search_engine.py
